package render

import (
	"io/ioutil"
	"log"
	"net/http"
)

func Render(writer http.ResponseWriter, request *http.Request) {
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

	//log.Println("%v", path)
	if path == "" {
		panic("lose path param")
	}

	// 从view目录中寻找文件
	body := readFile("view" + path + ".html")

	_, err = writer.Write(body)

	if err != nil {
		panic(err.Error())
	}
}

func readFile(path string) []byte {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err.Error())
	}

	return body
}
