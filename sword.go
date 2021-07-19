package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/sunshinev/go-sword/config"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/sunshinev/go-sword/assets/resource"

	"github.com/sunshinev/go-sword/assets/view"

	_ "github.com/go-sql-driver/mysql"
)

type gZipWriter struct {
	gz *gzip.Writer
	http.ResponseWriter
}

func (u *gZipWriter) Write(p []byte) (int, error) {
	return u.gz.Write(p)
}

type Sword struct {
}

func Init() *Sword {
	// 初始化配置
	config.Config{}.InitConfig()

	return &Sword{}
}

func (s *Sword) Run() {
	// Default Route
	http.HandleFunc("/api/model/table_list", s.tableList)
	http.HandleFunc("/api/model/preview", s.Preview)
	http.HandleFunc("/api/model/generate", s.Generate)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		//Static file route
		fs := assetfs.AssetFS{
			Asset:     resource.Asset,
			AssetDir:  resource.AssetDir,
			AssetInfo: resource.AssetInfo,
			Prefix:    "resource/dist",
		}

		handle := http.FileServer(&fs)
		w.Header().Set("Content-Encoding", "gzip")

		gz := gzip.NewWriter(w)
		newWriter := &gZipWriter{
			gz:             gz,
			ResponseWriter: w,
		}

		defer gz.Close()

		handle.ServeHTTP(newWriter, r)
	})

	// Render vue component
	http.HandleFunc("/render", s.Render)

	s.Welcome()

	// Start server
	err := http.ListenAndServe(":"+config.GlobalConfig.ServerPort, nil)
	if err != nil {
		log.Fatalf("Go-sword start err: %v", err)
	}
}

type GenerateParams struct {
	TableName string   `json:"table_name"`
	Files     []string `json:"files"`
}

// Get database table list
func (s *Sword) tableList(w http.ResponseWriter, r *http.Request) {
	rows, err := config.GlobalConfig.DbConn.Query("SHOW TABLES")
	if err != nil {
		panic(err.Error())
	}

	tables := []string{}

	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		tables = append(tables, tableName)
	}

	jsonData, err := json.Marshal(Ret{
		Code: http.StatusOK,
		Data: List{
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

	g := Generator{}.Init()
	g.Preview(data["table_name"])

	ret, err := json.Marshal(Ret{
		Code: http.StatusOK,
		Data: List{
			List: &g.FileList,
		},
	})
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

	var data = &GenerateParams{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err.Error())
	}

	if data.TableName == "" {
		panic("tableName is empty")
	}

	if len(data.Files) == 0 {
		panic("Files is empty")
	}

	g := Generator{}.Init()
	g.Generate(s, data.TableName, data.Files)

	ret, err := json.Marshal(Ret{
		Code: http.StatusOK,
		Data: List{
			List: &g.GFileList,
		},
	})

	_, err = w.Write(ret)
	if err != nil {
		panic(err.Error())
	}

}

func (s *Sword) Render(writer http.ResponseWriter, request *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic %v", err)
		}
	}()
	// 解析参数，映射到文件
	err := request.ParseForm()
	if err != nil {
		panic(err.Error())
	}

	path := request.FormValue("path")

	if path == "" {
		panic("lose path param")
	}

	// 从view目录中寻找文件
	//body := s.readFile("view" + path + ".html")
	body, err := view.Asset("view" + path + ".html")

	_, err = writer.Write(body)

	if err != nil {
		panic(err.Error())
	}
}

func (s *Sword) readFile(path string) []byte {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err.Error())
	}

	return body
}

func (s *Sword) Welcome() {
	c := config.GlobalConfig
	str := "Go-Sword will create new project named " + c.RootPath + " in current directory" +
		"\n\n[Server info]" +
		"\nServer port : " + c.ServerPort +
		"\nProject module : " + c.RootPath +
		"\n\n[db info]" +
		"\nMySQL host : " + c.DatabaseSet.Host +
		"\nMySQL port : " + strconv.Itoa(c.DatabaseSet.Port) +
		"\nMySQL user : " + c.DatabaseSet.User +
		"\nMySQL password : " + c.DatabaseSet.Password +
		"\n\nStart successful, server is running ...\n" +
		"Please request: " +
		strings.Join([]string{"http://localhost:", c.ServerPort}, "") +
		"\n"

	fmt.Println(str)
}
