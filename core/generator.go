package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sunshinev/go-sword/assets/resource"
	"github.com/sunshinev/go-sword/assets/stub"
	"github.com/sunshinev/go-sword/config"
	"github.com/sunshinev/go-sword/core/utils"
)

type Generator struct {
	Columns         []string
	ColumnDataTypes map[string]string
	Struc           string
	PackageName     string
	StructName      string
	TableName       string
	// Files need to generate
	FileList []*FileInstance
	config   *config.Config
	// Generate files list
	GFileList []*FileInstance
}

type FileInstance struct {
	FileName       string `json:"file_name"`
	FilePath       string `json:"file_path"`
	FileContent    string `json:"file_content"`
	FileOldContent string `json:"file_old_content"`
	IsDiff         bool   `json:"is_diff"`
	IsNew          bool   `json:"is_new"`
}

func (g Generator) Init() *Generator {
	return &Generator{
		config:          config.GlobalConfig,
		ColumnDataTypes: map[string]string{},
	}
}

// Entry
func (g *Generator) parseTable(table string) {

	columnDataTypes, err := utils.Db2struct{}.Convert(table)
	if err != nil {
		panic(err.Error())
	}

	// Set columns
	for _, r := range *columnDataTypes {
		g.Columns = append(g.Columns, r.ColumnName)
		g.ColumnDataTypes[r.ColumnName] = r.DataType
	}
	g.Columns = utils.ResortMySQLFields(&g.Columns)
	// Set TableName
	g.TableName = table
	// Set PackageName
	g.PackageName = table
	// Set StructName
	structName := strings.Replace(strings.Title(strings.Replace(table, "_", " ", -1)), " ", "", -1)
	g.StructName = structName
	// Set Content
	g.Struc = utils.Db2struct{}.FetchWholeStructFile("model", structName, table, columnDataTypes)
}

func (g *Generator) Preview(table string) {
	g.parseTable(table)

	g.gGoModFile()      // go.mod
	g.gMainFile()       // Main.go
	g.gCoreFile()       // Core.go
	g.gRouteFile()      // Route.go
	g.gModelFile()      // Model
	g.gControllerFile() // Controller
	g.gResponseFile()   // Response

	g.gHtmlDefaultFile() // default.html
	g.gHtmlListFile()    // list.htm
	g.gHtmlCreateFile()
	g.gHtmlDetailFile()
	g.gHtmlEditFile()

	for _, file := range g.FileList {
		g.handleFile(file)
	}

}

func (g *Generator) Generate(s *Sword, table string, files []string) {
	g.Preview(table)

	for _, file := range g.FileList {
		var path = file.FilePath
		// Files filter, only create selected file
		if !utils.IsContain(path, files) {
			continue
		}
		// Remember files need to generate
		g.GFileList = append(g.GFileList, file)

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

	// Explode resource, Recreate every time you click generate button
	g.gResourceRestore()
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
	if strings.Contains(g.Struc, "time.Time") {
		str := "package " + g.PackageName
		newStr := "package model" + `
import "time"
`
		return strings.Replace(g.Struc, str, newStr, 1)
	}

	return g.Struc
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
	importModel := g.config.RootPath + "/model"
	importResponse := g.config.RootPath + "/core/response"

	content := string(data)

	content = strings.ReplaceAll(content, "<<package_name>>", packageName)
	content = strings.ReplaceAll(content, "<<model_struct>>", modelStruct)
	content = strings.ReplaceAll(content, "<<import_model>>", importModel)
	content = strings.ReplaceAll(content, "<<import_response>>", importResponse)

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
	var fieldsType = ""

	for _, name := range g.Columns {
		columnList = columnList + fmt.Sprintf("{title:'%s', key:'%s'},\n", name, name)
		searchFields = searchFields + fmt.Sprintf("'%s',\n", name)

		fieldsType = fieldsType + fmt.Sprintf("%s:'%s',\n", name, utils.ConvertFieldsType2Js(g.ColumnDataTypes[name]))
	}

	content := string(data)

	content = strings.ReplaceAll(content, "<<table_name>>", g.TableName)
	content = strings.ReplaceAll(content, "<<js_data_column_list>>", columnList)
	content = strings.ReplaceAll(content, "<<js_data_search_fields>>", searchFields)
	content = strings.ReplaceAll(content, "<<js_data_fields_type>>", fieldsType)

	return content
}

func (g *Generator) gHtmlCreateFile() {
	var file = &FileInstance{
		FilePath:    filepath.Join(g.config.RootPath, "view", g.TableName, "create.html"),
		FileName:    "create.html",
		FileContent: g.createCreateHtml(),
	}

	g.FileList = append(g.FileList, file)
}

// Create  file content
func (g *Generator) createCreateHtml() string {
	// Read stub
	data, err := stub.Asset("stub/html/create.stub")
	if err != nil {
		panic(err.Error())
	}

	// replace
	var info = ""
	var fieldsType = ""

	for _, name := range g.Columns {
		// Create ignore `id` field
		if name == "id" || name == "created_at" || name == "updated_at" {
			continue
		}
		info = info + fmt.Sprintf("%s:'',\n", name)

		fieldsType = fieldsType + fmt.Sprintf("%s:'%s',\n", name, utils.ConvertFieldsType2Js(g.ColumnDataTypes[name]))

	}

	content := string(data)

	content = strings.ReplaceAll(content, "<<js_data_fields_type>>", fieldsType)
	content = strings.ReplaceAll(content, "<<table_name>>", g.TableName)
	content = strings.ReplaceAll(content, "<<js_data_info>>", info)

	return content
}

func (g *Generator) gHtmlEditFile() {
	var file = &FileInstance{
		FilePath:    filepath.Join(g.config.RootPath, "view", g.TableName, "edit.html"),
		FileName:    "edit.html",
		FileContent: g.createEditHtml(),
	}

	g.FileList = append(g.FileList, file)
}

// Create  file content
func (g *Generator) createEditHtml() string {
	// Read stub
	data, err := stub.Asset("stub/html/edit.stub")
	if err != nil {
		panic(err.Error())
	}

	// replace
	var info = ""
	var fieldsType = ""

	for _, name := range g.Columns {
		info = info + fmt.Sprintf("%s:'',\n", name)
		fieldsType = fieldsType + fmt.Sprintf("%s:'%s',\n", name, utils.ConvertFieldsType2Js(g.ColumnDataTypes[name]))
	}

	content := string(data)

	content = strings.ReplaceAll(content, "<<js_data_fields_type>>", fieldsType)
	content = strings.ReplaceAll(content, "<<table_name>>", g.TableName)
	content = strings.ReplaceAll(content, "<<js_data_info>>", info)

	return content
}

func (g *Generator) gHtmlDetailFile() {
	var file = &FileInstance{
		FilePath:    filepath.Join(g.config.RootPath, "view", g.TableName, "detail.html"),
		FileName:    "detail.html",
		FileContent: g.createDetailHtml(),
	}

	g.FileList = append(g.FileList, file)
}

// Create  file content
func (g *Generator) createDetailHtml() string {
	// Read stub
	data, err := stub.Asset("stub/html/detail.stub")
	if err != nil {
		panic(err.Error())
	}

	// replace
	var info = ""

	for _, name := range g.Columns {
		info = info + fmt.Sprintf("%s:'',\n", name)
	}

	content := string(data)

	content = strings.ReplaceAll(content, "<<table_name>>", g.TableName)
	content = strings.ReplaceAll(content, "<<js_data_info>>", info)

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
	http.HandleFunc("/api/%s/detail", %s.Detail(db))
	http.HandleFunc("/api/%s/create", %s.Create(db))
	http.HandleFunc("/api/%s/edit", %s.Edit(db))
	http.HandleFunc("/api/%s/batch_delete", %s.BatchDelete(db))`

	str = strings.ReplaceAll(str, "%s", g.TableName)

	var importStr = `"%s/controller/%s"`
	importStr = fmt.Sprintf(importStr, g.config.RootPath, g.TableName)

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

func (g *Generator) gGoModFile() {
	var file = &FileInstance{
		FilePath:    filepath.Join(g.config.RootPath, "go.mod"),
		FileName:    "go.mod",
		FileContent: g.createGoModContent(),
	}

	g.FileList = append(g.FileList, file)
}

func (g *Generator) gMainFile() {
	var file = &FileInstance{
		FilePath:    filepath.Join(g.config.RootPath, "main.go"),
		FileName:    "main.go",
		FileContent: g.createMainContent(),
	}

	g.FileList = append(g.FileList, file)
}

func (g *Generator) createGoModContent() string {
	// Read stub
	data, err := stub.Asset("stub/go.mod.stub")
	if err != nil {
		panic(err.Error())
	}
	content := string(data)
	content = strings.ReplaceAll(content, "<<module_name>>", g.config.RootPath)
	return content
}

func (g *Generator) createMainContent() string {
	// Read stub
	data, err := stub.Asset("stub/main.stub")
	if err != nil {
		panic(err.Error())
	}

	str := strings.Join([]string{g.config.RootPath, "core"}, "/")

	content := string(data)

	content = strings.ReplaceAll(content, "<<db_host>>", g.config.DatabaseSet.Host)
	content = strings.ReplaceAll(content, "<<db_user>>", g.config.DatabaseSet.User)
	content = strings.ReplaceAll(content, "<<db_password>>", g.config.DatabaseSet.Password)
	content = strings.ReplaceAll(content, "<<db_port>>", strconv.Itoa(g.config.DatabaseSet.Port))
	content = strings.ReplaceAll(content, "<<db_database>>", g.config.DatabaseSet.Database)
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

	str := strings.Join([]string{g.config.RootPath, "route"}, "/")

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
	var menu = fmt.Sprintf(`{icon: 'ios-people',title: '%s',name:'%s_list'},`, g.TableName, g.TableName)
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
	content = strings.ReplaceAll(content, "<<title>>", g.config.RootPath)

	return content
}

func (g *Generator) gResponseFile() {
	var file = &FileInstance{
		FilePath:    filepath.Join(g.config.RootPath, "core", "response", "response.go"),
		FileName:    "response.go",
		FileContent: g.createResponseContent(),
	}

	g.FileList = append(g.FileList, file)
}

func (g *Generator) createResponseContent() string {
	// Read stub
	data, err := stub.Asset("stub/core/response/response.stub")
	if err != nil {
		panic(err.Error())
	}

	return string(data)
}

func (g *Generator) gResourceRestore() {
	err := resource.RestoreAssets(g.config.RootPath, "resource/dist")
	if err != nil {
		panic(err.Error())
	}
}

func (g *Generator) readFile(path string) string {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ""
		}
	}

	file, err := os.Open(path)
	if err != nil {
		panic(err.Error())
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err.Error())
	}

	return string(data)
}

func (g *Generator) contentDiff(oldContent *string, newContent *string) (isNew bool, isDiff bool) {

	if *oldContent == "" {
		isNew = true
	} else {
		isNew = false
	}
	if *oldContent == *newContent {
		isDiff = false
	} else {
		// If the file is new, then do not show diff status
		if isNew {
			isDiff = false
		} else {
			isDiff = true
		}
	}

	return
}

func (g *Generator) handleFile(file *FileInstance) {
	oldContent := g.readFile(file.FilePath)
	file.FileOldContent = oldContent
	isNew, isDiff := g.contentDiff(&oldContent, &file.FileContent)
	file.IsNew = isNew
	file.IsDiff = isDiff
}
