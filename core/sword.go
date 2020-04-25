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

	// Check MySQL connection
	err := s.connDb()
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Default Route
	http.HandleFunc("/api/model/table_list", s.tableList)
	http.HandleFunc("/api/model/preview", s.Preview)
	http.HandleFunc("/api/model/generate", s.Generate)

	// Static file route
	fs := assetfs.AssetFS{
		Asset:     resource.Asset,
		AssetDir:  resource.AssetDir,
		AssetInfo: resource.AssetInfo,
		Prefix:    "resource/dist",
	}
	http.Handle("/", http.FileServer(&fs))

	// Render vue component
	http.HandleFunc("/render", render.Render)

	// Start server
	err = http.ListenAndServe(":"+s.Config.ServerPort, nil)
	if err != nil {
		log.Fatalf("Go-sword start err: %v", err)
	}
}

// Check MySQL database connection
func (s *Sword) connDb() (err error) {
	c := s.Config.Database
	// user:password@(localhost)/dbname?charset=utf8&parseTime=True&loc=Local
	db, err := sql.Open("mysql", c.User+":"+c.Password+"@tcp("+c.Host+":"+strconv.Itoa(c.Port)+")/"+c.Database+"?&parseTime=True")
	if err != nil {
		return
	}

	s.Db = db

	return nil
}

// Get database table list
func (s *Sword) tableList(w http.ResponseWriter, r *http.Request) {
	rows, err := s.Db.Query("SHOW TABLES")
	if err != nil {
		panic(err.Error())
	}

	tables := []string{}

	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
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

	g := Generator{}
	g.Init(s.Config)
	g.Preview(s.Config.Database, data["table_name"])

	ret, err := json.Marshal(&g.FileList)
	if err != nil {
		panic(err.Error())
	}

	_, err = w.Write(ret)
	if err != nil {
		panic(err.Error())
	}

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

	g := Generator{}
	g.Init(s.Config)
	g.Generate(s.Config.Database, data["table_name"])

	ret, err := json.Marshal(&g.FileList)

	_, err = w.Write(ret)
	if err != nil {
		panic(err.Error())
	}

}
