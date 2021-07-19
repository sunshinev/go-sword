package gosword

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

	"github.com/sunshinev/go-sword/assets/view"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/sunshinev/go-sword/assets/resource"

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

func Init(configFile string) *Sword {
	// 初始化配置
	err := config.Config{}.LoadConfig(configFile)

	if err != nil {
		log.Fatalf("sword init err %v", err)
	}
	return &Sword{}
}

func (s *Sword) Run() {
	// 数据表列表
	http.HandleFunc("/api/model/table_list", s.tableList)
	// 预览
	http.HandleFunc("/api/model/preview", s.Preview)
	// 创建生成文件
	http.HandleFunc("/api/model/generate", s.Generate)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 静态文件路由
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

	//Start server
	go func() {
		err := http.ListenAndServe(":"+config.GlobalConfig.ServerPort, nil)
		if err != nil {
			log.Fatalf("Go-sword start err: %v", err)
		}
	}()
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
	// body,err := ioutil.ReadFile("view/"+path+".html")
	// 这里使用go-bindata释放到views目录，通过go文件加载资源，少了文件读写；使用bindata每次修改完html文件之后，需要重新生成views资源
	body, err := view.Asset("view" + path + ".html")
	if err != nil {
		panic(err)
	}
	_, err = writer.Write(body)
	if err != nil {
		panic(err)
	}
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
