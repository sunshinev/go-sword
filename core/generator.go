package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sunshinev/go-sword/assets/stub"

	"github.com/sunshinev/go-sword/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sunshinev/db2struct"
)

type Generator struct {
	Columns     []string
	Struc       []byte
	PackageName string
	StructName  string
	TableName   string
	// Files need to generate
	FileList []*FileInstance
	config   *config.Config
}

type FileInstance struct {
	FileName    string
	FilePath    string
	FileContent string
	Diff        string
	IsDiff      bool
	IsNew       bool
}

func (g *Generator) Init(c *config.Config) *Generator {
	g.config = c
	return g
}

// Entry
func (g *Generator) parseTable(c *config.DbSet, table string) {
	// Use db2struct (https://github.com/sunshinev/db2struct) Forked (https://github.com/Shelnutt2/db2struct) and modify Generate func
	columnDataTypes, err := db2struct.GetColumnsFromMysqlTable(c.User, c.Password, c.Host, c.Port, c.Database, table)
	if err != nil {
		panic(err.Error())
	}

	structName := strings.Replace(table, "_", " ", -1)
	structName = strings.Title(structName)
	structName = strings.Replace(structName, " ", "", -1)

	// Set columns
	for name := range *columnDataTypes {
		g.Columns = append(g.Columns, name)
	}

	struc, err := db2struct.Generate(*columnDataTypes, table, structName, table, true, true, false)
	if err != nil {
		panic(err.Error())
	}

	// Set TableName
	g.TableName = table
	// Set PackageName
	g.PackageName = table
	// Set StructName
	g.StructName = structName
	// Set Content
	g.Struc = struc
}

func (g *Generator) Preview(c *config.DbSet, table string) {
	g.parseTable(c, table)
	// Main.go
	g.gMainFile()
	// Core.go
	g.gCoreFile()
	// Route.go
	g.gRouteFile()
	// Model
	g.gModelFile()
	// Controller
	g.gControllerFile()
	// Html
	// default.html
	g.gHtmlDefaultFile()
	// list.html
	g.gHtmlListFile()

}

func (g *Generator) Generate(c *config.DbSet, table string) {
	g.Preview(c, table)

	for _, file := range g.FileList {
		var path = file.FilePath
		_, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				// Create new file
				err = os.MkdirAll(strings.ReplaceAll(path, file.FileName, ""), 0755)
				if err != nil {
					panic(err.Error())
				}
			}
		}

		newFile, err := os.Create(path)
		if err != nil {
			panic(err.Error())
		}
		_, err = newFile.Write([]byte(file.FileContent))
		if err != nil {
			panic(err.Error())
		}
	}
}

func (g *Generator) gModelFile() {
	content := g.createModelContent()
	file := &FileInstance{
		FilePath:    filepath.Join(g.config.RootPath, "model", g.TableName+".go"),
		FileName:    g.TableName + ".go",
		FileContent: content,
	}

	g.FileList = append(g.FileList, file)
}

// Create model file content
func (g *Generator) createModelContent() string {
	// Modify g.Struc & add import
	if strings.Contains(string(g.Struc), "time.Time") {
		str := "package " + g.PackageName
		newStr := "package model" + `
import "time"
`
		return strings.Replace(string(g.Struc), str, newStr, 1)
	}

	return string(g.Struc)
}

func (g *Generator) gControllerFile() {
	var file = &FileInstance{
		FilePath:    filepath.Join(g.config.RootPath, "controller", g.TableName, g.TableName+".go"),
		FileName:    g.TableName + ".go",
		FileContent: g.createControllerContent(),
	}

	g.FileList = append(g.FileList, file)
}

// Create model file content
func (g *Generator) createControllerContent() string {
	// Read stub
	data, err := stub.Asset("stub/controller/controller.stub")
	if err != nil {
		panic(err.Error())
	}

	// replace
	packageName := g.TableName
	modelStruct := "model." + g.StructName
	importModel := g.config.ModuleName + "/" + g.config.RootPath + "/model"
	importCore := g.config.ModuleName + "/" + g.config.RootPath + "/core"

	content := string(data)

	content = strings.ReplaceAll(content, "<<package_name>>", packageName)
	content = strings.ReplaceAll(content, "<<model_struct>>", modelStruct)
	content = strings.ReplaceAll(content, "<<import_model>>", importModel)
	content = strings.ReplaceAll(content, "<<import_core>>", importCore)

	return content
}

func (g *Generator) gHtmlListFile() {
	var file = &FileInstance{
		FilePath:    filepath.Join(g.config.RootPath, "view", g.TableName, "list.html"),
		FileName:    "list.html",
		FileContent: g.createListHtml(),
	}

	g.FileList = append(g.FileList, file)
}

// Create model file content
func (g *Generator) createListHtml() string {
	// Read stub
	data, err := stub.Asset("stub/html/list.stub")
	if err != nil {
		panic(err.Error())
	}

	// replace
	var columnList = ""
	var searchFields = ""

	for _, name := range g.Columns {
		columnList = columnList + fmt.Sprintf("{title:'%s', key:'%s'},\n", name, name)
		searchFields = searchFields + fmt.Sprintf("'%s',\n", name)
	}

	content := string(data)

	content = strings.ReplaceAll(content, "<<table_name>>", g.TableName)
	content = strings.ReplaceAll(content, "<<js_data_column_list>>", columnList)
	content = strings.ReplaceAll(content, "<<js_data_search_fields>>", searchFields)

	return content
}

func (g *Generator) gRouteFile() {
	var file = &FileInstance{
		FilePath:    filepath.Join(g.config.RootPath, "route", "route.go"),
		FileName:    "route.go",
		FileContent: g.getRouteContent(),
	}

	g.FileList = append(g.FileList, file)
}

func (g *Generator) getRouteContent() string {
	var path = filepath.Join(g.config.RootPath, "route", "route.go")

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return g.createRouteContent("")
		}

		panic(err.Error())
	}

	return g.createRouteContent(path)
}

func (g *Generator) createRouteContent(path string) string {

	var data = []byte{}
	var err error

	// Read bytes from stub
	if path == "" {
		data, err = stub.Asset("stub/route/route.stub")
		if err != nil {
			panic(err.Error())
		}
	} else {
		// Open created route.go from rootPath
		file, err := os.Open(path)
		if err != nil {
			panic(err.Error())
		}

		data, err = ioutil.ReadAll(file)
		if err != nil {
			panic(err.Error())
		}
	}

	// replace
	var str = `
	// Route tag %s
	http.HandleFunc("/api/%s/list", %s.List(db))
	http.HandleFunc("/api/%s/delete", %s.Delete(db))
	http.HandleFunc("/api/%s/detail", %s.Detail(db))`

	str = strings.ReplaceAll(str, "%s", g.TableName)

	var importStr = `"%s/%s/controller/%s"`
	importStr = fmt.Sprintf(importStr, g.config.ModuleName, g.config.RootPath, g.TableName)

	// Check if content repeated,if true then ignore replace
	content := string(data)
	if !strings.Contains(content, str) {
		content = strings.ReplaceAll(content, "// ----Route-end----", str+`
	// ----Route-end----`)
	}
	if !strings.Contains(content, importStr) {
		content = strings.ReplaceAll(content, "// ----Import----", importStr+`
	// ----Import----`)
	}

	return content
}

func (g *Generator) gMainFile() {
	var file = &FileInstance{
		FilePath:    filepath.Join(g.config.RootPath, "main.go"),
		FileName:    "main.go",
		FileContent: g.createMainContent(),
	}

	g.FileList = append(g.FileList, file)
}

func (g *Generator) createMainContent() string {
	// Read stub
	data, err := stub.Asset("stub/main.stub")
	if err != nil {
		panic(err.Error())
	}

	str := strings.Join([]string{g.config.ModuleName, g.config.RootPath, "core"}, "/")

	content := string(data)

	content = strings.ReplaceAll(content, "<<import_core>>", str)
	return content
}

func (g *Generator) gCoreFile() {
	var file = &FileInstance{
		FilePath:    filepath.Join(g.config.RootPath, "core", "core.go"),
		FileName:    "core.go",
		FileContent: g.createCoreContent(),
	}

	g.FileList = append(g.FileList, file)
}
func (g *Generator) createCoreContent() string {
	// Read stub
	data, err := stub.Asset("stub/core/core.stub")
	if err != nil {
		panic(err.Error())
	}

	str := strings.Join([]string{g.config.ModuleName, g.config.RootPath, "route"}, "/")

	content := string(data)

	content = strings.ReplaceAll(content, "<<import_route>>", str)
	return content
}

func (g *Generator) gHtmlDefaultFile() {
	var file = &FileInstance{
		FilePath:    filepath.Join(g.config.RootPath, "view", "layout", "default.html"),
		FileName:    "default.html",
		FileContent: g.getHtmlDefaultFile(),
	}

	g.FileList = append(g.FileList, file)
}

func (g *Generator) getHtmlDefaultFile() string {
	var path = filepath.Join(g.config.RootPath, "view", "layout", "default.html")

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return g.createDefaultHtml("")
		}

		panic(err.Error())
	}

	return g.createDefaultHtml(path)
}

func (g *Generator) createDefaultHtml(path string) string {
	var data = []byte{}
	var err error

	// Read bytes from stub
	if path == "" {
		data, err = stub.Asset("stub/layout/default.stub")
		if err != nil {
			panic(err.Error())
		}
	} else {
		// Open created route.go from rootPath
		file, err := os.Open(path)
		if err != nil {
			panic(err.Error())
		}

		data, err = ioutil.ReadAll(file)
		if err != nil {
			panic(err.Error())
		}
	}

	// replace
	var menu = fmt.Sprintf(`{icon: 'ios-people',title: '%s_list',name:'%s_list'},`, g.TableName, g.TableName)
	var routeSets = []string{"list", "create", "detail", "edit"}
	var route = ""

	for _, set := range routeSets {
		route += strings.ReplaceAll(strings.ReplaceAll(`{
                    name: 'user_list',
                    path: '/user/list',
                    url: '/render?path=/user/list'
                },`, "user", g.TableName), "list", set)
	}

	var defaultRoute = g.TableName + "_list"

	// Check if content repeated,if true then ignore replace
	content := string(data)
	if !strings.Contains(content, menu) {
		content = strings.ReplaceAll(content, "// ----Menus-Add-----", menu+`
	                    // ----Menus-Add-----`)
	}
	if !strings.Contains(content, route) {
		content = strings.ReplaceAll(content, "// ----Routes-Add-----", route+`
                // ----Routes-Add-----`)
	}

	content = strings.ReplaceAll(content, "<<default_route>>", defaultRoute)

	return content
}

//
//// Read created file content
//func (g *Generator) readCreatedFile(path string) string {
//	_, err := os.Stat(path)
//	if err != nil {
//		if os.IsNotExist(err) {
//			return ""
//		}
//
//		panic(err.Error())
//	}
//
//	fileInfo, err := os.OpenFile(path, os.O_WRONLY, 0755)
//	if err != nil {
//		panic(err.Error())
//	}
//
//	defer fileInfo.Close()
//
//	data, err := ioutil.ReadAll(fileInfo)
//
//	return string(data)
//}
