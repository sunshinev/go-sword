package core

import (
	"fmt"
	"io/ioutil"
	"os"
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
	// Model
	var modelFile = &FileInstance{
		FilePath:    strings.Join([]string{"model", g.TableName + ".go"}, string(os.PathSeparator)),
		FileName:    g.TableName + ".go",
		FileContent: g.createModelContent(),
	}

	g.FileList = append(g.FileList, modelFile)

	// Controller
	var controllerFile = &FileInstance{
		FilePath:    strings.Join([]string{"controller", g.TableName, g.TableName + ".go"}, string(os.PathSeparator)),
		FileName:    g.TableName + ".go",
		FileContent: g.createControllerContent(),
	}

	g.FileList = append(g.FileList, controllerFile)

	// Html
	// list.html
	var listHtmlFile = &FileInstance{
		FilePath:    strings.Join([]string{"view", g.TableName, "list.html"}, string(os.PathSeparator)),
		FileName:    "list.html",
		FileContent: g.createListHtml(),
	}

	g.FileList = append(g.FileList, listHtmlFile)
}

func (g *Generator) Generate(c *config.DbSet, table string) {
	g.Preview(c, table)

	for _, file := range g.FileList {
		var path = strings.Join([]string{g.config.RootPath, file.FilePath}, string(os.PathSeparator))
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

	// Handle route additional
	g.generateRoute()

	// Handle main.go additional
	g.generateMain()

	// Handle core.go additional
	g.generateCore()

	// Handle default.html additional
	g.generateDefaultHtml()
}

func (g *Generator) generateRoute() {
	var path = strings.Join([]string{g.config.RootPath, "route", "route.go"}, string(os.PathSeparator))
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Create new file
			err = os.MkdirAll(strings.ReplaceAll(path, "route.go", ""), 0755)
			if err != nil {
				panic(err.Error())
			}

			// Create new route.go
			newFile, err := os.Create(path)
			if err != nil {
				panic(err.Error())
			}

			routeStubPath := strings.Join([]string{"stub", "route", "route.stub"}, string(os.PathSeparator))
			_, err = newFile.Write([]byte(g.createRouteContent(routeStubPath)))
			if err != nil {
				panic(err.Error())
			}

			return
		}

		panic(err.Error())
	}

	newFile, err := os.OpenFile(path, os.O_WRONLY, 0755)
	if err != nil {
		panic(err.Error())
	}

	_, err = newFile.Write([]byte(g.createRouteContent(path)))
	if err != nil {
		panic(err.Error())
	}
}

func (g *Generator) generateDefaultHtml() {
	var path = strings.Join([]string{g.config.RootPath, "view", "layout", "default.html"}, string(os.PathSeparator))
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Create new file
			err = os.MkdirAll(strings.ReplaceAll(path, "default.html", ""), 0755)
			if err != nil {
				panic(err.Error())
			}

			// Create new route.go
			newFile, err := os.Create(path)
			if err != nil {
				panic(err.Error())
			}

			stubPath := strings.Join([]string{"stub", "layout", "default.stub"}, string(os.PathSeparator))
			_, err = newFile.Write([]byte(g.createDefaultHtml(stubPath)))
			if err != nil {
				panic(err.Error())
			}

			return
		}

		panic(err.Error())
	}

	newFile, err := os.OpenFile(path, os.O_WRONLY, 0755)
	if err != nil {
		panic(err.Error())
	}

	_, err = newFile.Write([]byte(g.createDefaultHtml(path)))
	if err != nil {
		panic(err.Error())
	}
}

func (g *Generator) generateMain() {
	var path = strings.Join([]string{g.config.RootPath, "main.go"}, string(os.PathSeparator))
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Create new file
			err = os.MkdirAll(strings.ReplaceAll(path, "main.go", ""), 0755)
			if err != nil {
				panic(err.Error())
			}

			newFile, err := os.Create(path)
			if err != nil {
				panic(err.Error())
			}
			_, err = newFile.Write([]byte(g.createMainContent()))
			if err != nil {
				panic(err.Error())
			}
		}
	}
}

func (g *Generator) generateCore() {
	var path = strings.Join([]string{g.config.RootPath, "core", "core.go"}, string(os.PathSeparator))
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Create new file
			err = os.MkdirAll(strings.ReplaceAll(path, "core.go", ""), 0755)
			if err != nil {
				panic(err.Error())
			}

			newFile, err := os.Create(path)
			if err != nil {
				panic(err.Error())
			}
			_, err = newFile.Write([]byte(g.createCoreContent()))
			if err != nil {
				panic(err.Error())
			}
		}
	}
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

	content := string(data)

	content = strings.ReplaceAll(content, "<<package_name>>", packageName)
	content = strings.ReplaceAll(content, "<<model_struct>>", modelStruct)
	content = strings.ReplaceAll(content, "<<import_model>>", importModel)

	return content
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

func (g *Generator) createRouteContent(filePath string) string {
	// Read stub
	file, err := os.Open(filePath)
	if err != nil {
		panic(err.Error())
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err.Error())
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

func (g *Generator) createDefaultHtml(filePath string) string {
	// Read stub
	file, err := os.Open(filePath)
	if err != nil {
		panic(err.Error())
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err.Error())
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
