package core

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/sunshinev/go-sword/assets/resource"

	assetfs "github.com/elazarl/go-bindata-assetfs"

	"github.com/sunshinev/go-sword/config"
	"github.com/sunshinev/go-sword/controller/render"
	"github.com/sunshinev/go-sword/model"
	"github.com/sunshinev/go-sword/response"

	_ "github.com/go-sql-driver/mysql"
)

// Engine
type Sword struct {
	Config *config.Config
	Db     *sql.DB
}

// Create Sword
func Default() *Sword {
	return &Sword{
		Config: &config.Config{},
	}
}

// Set Config like Database
func (s *Sword) SetConfig(cfg *config.Config) {
	s.Config = cfg
}

func (s *Sword) Run() {

	// Cache Panic
	defer func() {
		if err := recover(); err != nil {
			log.Printf("%v", err)
		}
	}()

	err := s.connDb()
	if err != nil {
		log.Fatalf("%v", err)
	}

	//http.HandleFunc("/sword/api/model/create", e.modelCreate)
	http.HandleFunc("/api/model/table_list", s.tableList)
	http.HandleFunc("/api/model/preview", s.Preview)
	http.HandleFunc("/api/model/generate", s.Generate)

	// home page
	fs := assetfs.AssetFS{
		Asset:     resource.Asset,
		AssetDir:  resource.AssetDir,
		AssetInfo: resource.AssetInfo,
		Prefix:    "resource/dist",
	}
	http.Handle("/", http.FileServer(&fs))

	// render vue component
	http.HandleFunc("/render", render.Render)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("服务启动失败 %v", err)
	}
}

func (s *Sword) connDb() (err error) {
	c := s.Config.Database
	db, err := sql.Open("mysql", c.User+":"+c.Password+"@tcp("+c.Host+":"+strconv.Itoa(c.Port)+")/"+c.Database+"?&parseTime=True")
	if err != nil {
		return
	}

	s.Db = db

	return nil
}

func (s *Sword) tableList(w http.ResponseWriter, r *http.Request) {
	// user:password@(localhost)/dbname?charset=utf8&parseTime=True&loc=Local

	rows, err := s.Db.Query("show tables")
	if err != nil {
		panic(err.Error())
	}

	tables := []string{}

	for rows.Next() {
		var tableName string
		rows.Scan(&tableName)
		tables = append(tables, tableName)
	}

	jsonData, err := json.Marshal(response.Ret{
		Code: http.StatusOK,
		Data: response.List{
			List: tables,
		},
	})

	_, err = w.Write(jsonData)
	if err != nil {
		panic(err.Error())
	}
}

func (s *Sword) Preview(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	var data map[string]string
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err.Error())
	}

	if data["table_name"] == "" {
		panic("tableName is empty")
	}

	m := model.ModelCreate{}
	m.Init(s.Config)
	m.Preview(s.Config.Database, data["table_name"])

	ret, err := json.Marshal(&m.FileList)

	w.Write(ret)
}

func (s *Sword) Generate(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	var data map[string]string
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err.Error())
	}

	if data["table_name"] == "" {
		panic("tableName is empty")
	}

	m := model.ModelCreate{}
	m.Init(s.Config)
	m.Generate(s.Config.Database, data["table_name"])

	ret, err := json.Marshal(&m.FileList)

	w.Write(ret)
}
