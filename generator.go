package gosword

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sunshinev/go-sword/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sunshinev/go-sword/assets/resource"
	"github.com/sunshinev/go-sword/assets/stub"
	"github.com/sunshinev/go-sword/utils"
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

func (s Generator) Init() *Generator {
	return &Generator{
		config:          config.GlobalConfig,
		ColumnDataTypes: map[string]string{},
	}
}

// Entry
func (s *Generator) parseTable(table string) {

	columnDataTypes, err := utils.Db2struct{}.Convert(table)
	if err != nil {
		panic(err.Error())
	}

	// Set columns
	for _, r := range *columnDataTypes {
		s.Columns = append(s.Columns, r.ColumnName)
		s.ColumnDataTypes[r.ColumnName] = r.DataType
	}
	s.Columns = utils.ResortMySQLFields(&s.Columns)
	// Set TableName
	s.TableName = table
	// Set PackageName
	s.PackageName = table
	// Set StructName
	structName := strings.Replace(strings.Title(strings.Replace(table, "_", " ", -1)), " ", "", -1)
	s.StructName = structName
	// Set Content
	s.Struc = utils.Db2struct{}.FetchWholeStructFile("model", structName, table, columnDataTypes)
}

func (s *Generator) Preview(table string) {
	s.parseTable(table)

	//s.gGoModFile()      // go.mod
	//s.gMainFile()       // Main.go
	s.gCoreFile()       // Core.go
	s.gRouteFile()      // Route.go
	s.gModelFile()      // Model
	s.gControllerFile() // Controller
	s.gResponseFile()   // Response

	s.gHtmlDefaultFile() // default.html
	s.gHtmlListFile()    // list.htm
	s.gHtmlCreateFile()
	s.gHtmlDetailFile()
	s.gHtmlEditFile()

	err := s.fmtCode()
	if err != nil {

	}

	for _, file := range s.FileList {
		s.handleFile(file)
	}

}

func (s *Generator) Generate(table string, files []string) {
	s.Preview(table)

	for _, file := range s.FileList {
		var path = file.FilePath
		// Files filter, only create selected file
		if !utils.IsContain(path, files) {
			continue
		}
		// Remember files need to generate
		s.GFileList = append(s.GFileList, file)

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
	s.gResourceRestore()
}

func (s *Generator) gModelFile() {
	content := s.createModelContent()
	file := &FileInstance{
		FilePath:    filepath.Join(s.config.RootPath, "model", s.TableName+".go"),
		FileName:    s.TableName + ".go",
		FileContent: content,
	}

	s.FileList = append(s.FileList, file)
}

// Create model file content
func (s *Generator) createModelContent() string {
	// Modify s.Struc & add import
	if strings.Contains(s.Struc, "time.Time") {
		str := "package " + s.PackageName
		newStr := "package model" + `
import "time"
`
		return strings.Replace(s.Struc, str, newStr, 1)
	}

	return s.Struc
}

func (s *Generator) gControllerFile() {
	var file = &FileInstance{
		FilePath:    filepath.Join(s.config.RootPath, "controller", s.TableName, s.TableName+".go"),
		FileName:    s.TableName + ".go",
		FileContent: s.createControllerContent(),
	}

	s.FileList = append(s.FileList, file)
}

// Create model file content
func (s *Generator) createControllerContent() string {
	// Read stub
	data, err := stub.Asset("stub/controller/controller.stub")
	if err != nil {
		panic(err.Error())
	}

	// replace
	packageName := s.TableName
	modelStruct := "model." + s.StructName
	importModel := fmt.Sprintf("%v/%v/%v", s.config.ModuleName, s.config.RootPath, "model")
	importResponse := fmt.Sprintf("%v/%v/%v", s.config.ModuleName, s.config.RootPath, "core/response")

	content := string(data)

	content = strings.ReplaceAll(content, "<<package_name>>", packageName)
	content = strings.ReplaceAll(content, "<<model_struct>>", modelStruct)
	content = strings.ReplaceAll(content, "<<import_model>>", importModel)
	content = strings.ReplaceAll(content, "<<import_response>>", importResponse)

	return content
}

func (s *Generator) gHtmlListFile() {
	var file = &FileInstance{
		FilePath:    filepath.Join(s.config.RootPath, "view", s.TableName, "list.html"),
		FileName:    "list.html",
		FileContent: s.createListHtml(),
	}

	s.FileList = append(s.FileList, file)
}

// Create model file content
func (s *Generator) createListHtml() string {
	// Read stub
	data, err := stub.Asset("stub/html/list.stub")
	if err != nil {
		panic(err.Error())
	}

	// replace
	var columnList = ""
	var searchFields = ""
	var fieldsType = ""

	for _, name := range s.Columns {
		columnList = columnList + fmt.Sprintf("{title:'%s', key:'%s'},\n", name, name)
		searchFields = searchFields + fmt.Sprintf("'%s',\n", name)

		fieldsType = fieldsType + fmt.Sprintf("%s:'%s',\n", name, utils.ConvertFieldsType2Js(s.ColumnDataTypes[name]))
	}

	content := string(data)

	content = strings.ReplaceAll(content, "<<table_name>>", s.TableName)
	content = strings.ReplaceAll(content, "<<js_data_column_list>>", columnList)
	content = strings.ReplaceAll(content, "<<js_data_search_fields>>", searchFields)
	content = strings.ReplaceAll(content, "<<js_data_fields_type>>", fieldsType)

	return content
}

func (s *Generator) gHtmlCreateFile() {
	var file = &FileInstance{
		FilePath:    filepath.Join(s.config.RootPath, "view", s.TableName, "create.html"),
		FileName:    "create.html",
		FileContent: s.createCreateHtml(),
	}

	s.FileList = append(s.FileList, file)
}

// Create  file content
func (s *Generator) createCreateHtml() string {
	// Read stub
	data, err := stub.Asset("stub/html/create.stub")
	if err != nil {
		panic(err.Error())
	}

	// replace
	var info = ""
	var fieldsType = ""

	for _, name := range s.Columns {
		// Create ignore `id` field
		if name == "id" || name == "created_at" || name == "updated_at" {
			continue
		}
		info = info + fmt.Sprintf("%s:'',\n", name)

		fieldsType = fieldsType + fmt.Sprintf("%s:'%s',\n", name, utils.ConvertFieldsType2Js(s.ColumnDataTypes[name]))

	}

	content := string(data)

	content = strings.ReplaceAll(content, "<<js_data_fields_type>>", fieldsType)
	content = strings.ReplaceAll(content, "<<table_name>>", s.TableName)
	content = strings.ReplaceAll(content, "<<js_data_info>>", info)

	return content
}

func (s *Generator) gHtmlEditFile() {
	var file = &FileInstance{
		FilePath:    filepath.Join(s.config.RootPath, "view", s.TableName, "edit.html"),
		FileName:    "edit.html",
		FileContent: s.createEditHtml(),
	}

	s.FileList = append(s.FileList, file)
}

// Create  file content
func (s *Generator) createEditHtml() string {
	// Read stub
	data, err := stub.Asset("stub/html/edit.stub")
	if err != nil {
		panic(err.Error())
	}

	// replace
	var info = ""
	var fieldsType = ""

	for _, name := range s.Columns {
		info = info + fmt.Sprintf("%s:'',\n", name)
		fieldsType = fieldsType + fmt.Sprintf("%s:'%s',\n", name, utils.ConvertFieldsType2Js(s.ColumnDataTypes[name]))
	}

	content := string(data)

	content = strings.ReplaceAll(content, "<<js_data_fields_type>>", fieldsType)
	content = strings.ReplaceAll(content, "<<table_name>>", s.TableName)
	content = strings.ReplaceAll(content, "<<js_data_info>>", info)

	return content
}

func (s *Generator) gHtmlDetailFile() {
	var file = &FileInstance{
		FilePath:    filepath.Join(s.config.RootPath, "view", s.TableName, "detail.html"),
		FileName:    "detail.html",
		FileContent: s.createDetailHtml(),
	}

	s.FileList = append(s.FileList, file)
}

// Create  file content
func (s *Generator) createDetailHtml() string {
	// Read stub
	data, err := stub.Asset("stub/html/detail.stub")
	if err != nil {
		panic(err.Error())
	}

	// replace
	var info = ""

	for _, name := range s.Columns {
		info = info + fmt.Sprintf("%s:'',\n", name)
	}

	content := string(data)

	content = strings.ReplaceAll(content, "<<table_name>>", s.TableName)
	content = strings.ReplaceAll(content, "<<js_data_info>>", info)

	return content
}

func (s *Generator) gRouteFile() {
	var file = &FileInstance{
		FilePath:    filepath.Join(s.config.RootPath, "route", "route.go"),
		FileName:    "route.go",
		FileContent: s.getRouteContent(),
	}

	s.FileList = append(s.FileList, file)
}

func (s *Generator) getRouteContent() string {
	var path = filepath.Join(s.config.RootPath, "route", "route.go")

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return s.createRouteContent("")
		}

		panic(err.Error())
	}

	return s.createRouteContent(path)
}

func (s *Generator) createRouteContent(path string) string {

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

	str = strings.ReplaceAll(str, "%s", s.TableName)

	var importStr = `"%s/controller/%s"`
	importStr = fmt.Sprintf(importStr, strings.Join([]string{s.config.ModuleName, s.config.RootPath}, "/"), s.TableName)

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

func (s *Generator) gGoModFile() {
	var file = &FileInstance{
		FilePath:    filepath.Join(s.config.RootPath, "go.mod"),
		FileName:    "go.mod",
		FileContent: s.createGoModContent(),
	}

	s.FileList = append(s.FileList, file)
}

func (s *Generator) gMainFile() {
	var file = &FileInstance{
		FilePath:    filepath.Join(s.config.RootPath, "main.go"),
		FileName:    "main.go",
		FileContent: s.createMainContent(),
	}

	s.FileList = append(s.FileList, file)
}

func (s *Generator) createGoModContent() string {
	// Read stub
	data, err := stub.Asset("stub/go.mod.stub")
	if err != nil {
		panic(err.Error())
	}
	content := string(data)
	content = strings.ReplaceAll(content, "<<module_name>>", s.config.ModuleName)
	return content
}

func (s *Generator) createMainContent() string {
	// Read stub
	data, err := stub.Asset("stub/main.stub")
	if err != nil {
		panic(err.Error())
	}

	str := strings.Join([]string{s.config.ModuleName, s.config.RootPath, "core"}, "/")

	content := string(data)

	content = strings.ReplaceAll(content, "<<db_host>>", s.config.DatabaseSet.Host)
	content = strings.ReplaceAll(content, "<<db_user>>", s.config.DatabaseSet.User)
	content = strings.ReplaceAll(content, "<<db_password>>", s.config.DatabaseSet.Password)
	content = strings.ReplaceAll(content, "<<db_port>>", strconv.Itoa(s.config.DatabaseSet.Port))
	content = strings.ReplaceAll(content, "<<db_database>>", s.config.DatabaseSet.Database)
	content = strings.ReplaceAll(content, "<<import_core>>", str)
	return content
}

func (s *Generator) gCoreFile() {

	var file = &FileInstance{
		FilePath:    filepath.Join(s.config.RootPath, "core", "core.go"),
		FileName:    "core.go",
		FileContent: s.createCoreContent(),
	}

	s.FileList = append(s.FileList, file)
}
func (s *Generator) createCoreContent() string {
	// Read stub
	data, err := stub.Asset("stub/core/core.stub")
	if err != nil {
		panic(err.Error())
	}

	str := strings.Join([]string{s.config.ModuleName, s.config.RootPath, "route"}, "/")

	content := string(data)

	content = strings.ReplaceAll(content, "<<import_route>>", str)
	return content
}

func (s *Generator) gHtmlDefaultFile() {
	var file = &FileInstance{
		FilePath:    filepath.Join(s.config.RootPath, "view", "layout", "default.html"),
		FileName:    "default.html",
		FileContent: s.getHtmlDefaultFile(),
	}

	s.FileList = append(s.FileList, file)
}

func (s *Generator) getHtmlDefaultFile() string {
	var path = filepath.Join(s.config.RootPath, "view", "layout", "default.html")

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return s.createDefaultHtml("")
		}

		panic(err.Error())
	}

	return s.createDefaultHtml(path)
}

func (s *Generator) createDefaultHtml(path string) string {
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
	var menu = fmt.Sprintf(`{icon: 'ios-people',title: '%s',name:'%s_list'},`, s.TableName, s.TableName)
	var routeSets = []string{"list", "create", "detail", "edit"}
	var route = ""

	for _, set := range routeSets {
		route += strings.ReplaceAll(strings.ReplaceAll(`{
                    name: 'user_list',
                    path: '/user/list',
                    url: '/render?path=/user/list'
                },`, "user", s.TableName), "list", set)
	}

	var defaultRoute = s.TableName + "_list"

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
	content = strings.ReplaceAll(content, "<<title>>", s.config.RootPath)

	return content
}

func (s *Generator) gResponseFile() {
	var file = &FileInstance{
		FilePath:    filepath.Join(s.config.RootPath, "core", "response", "response.go"),
		FileName:    "response.go",
		FileContent: s.createResponseContent(),
	}

	s.FileList = append(s.FileList, file)
}

func (s *Generator) createResponseContent() string {
	// Read stub
	data, err := stub.Asset("stub/core/response/response.stub")
	if err != nil {
		panic(err.Error())
	}

	return string(data)
}

func (s *Generator) gResourceRestore() {
	err := resource.RestoreAssets(s.config.RootPath, "resource/dist")
	if err != nil {
		panic(err.Error())
	}
}

func (s *Generator) readFile(path string) string {
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

func (s *Generator) contentDiff(oldContent *string, newContent *string) (isNew bool, isDiff bool) {

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

func (s *Generator) handleFile(file *FileInstance) {
	oldContent := s.readFile(file.FilePath)
	file.FileOldContent = oldContent
	isNew, isDiff := s.contentDiff(&oldContent, &file.FileContent)
	file.IsNew = isNew
	file.IsDiff = isDiff
}

// 讲代码进行gofmt格式化
func (s *Generator) fmtCode() error {
	if len(s.FileList) == 0 {
		return nil
	}

	// 创建临时目录
	dir, err := ioutil.TempDir("/tmp", "temp")
	if err != nil {
		log.Printf("create tempdir %v err %v", dir, err)
		return err
	}
	// 处理完移除临时文件
	defer os.RemoveAll(dir)

	for _, v := range s.FileList {
		if !strings.HasSuffix(v.FileName, ".go") {
			continue
		}

		filePath := strings.Replace(v.FilePath, "/", "_", -1)
		file := fmt.Sprintf("%v/%v", dir, filePath)

		err := ioutil.WriteFile(file, []byte(v.FileContent), 0644)
		if err != nil {
			log.Printf("fmt code write temp file %v err %v", file, err)
			return err
		}
	}

	err = s.gofmt(dir)
	if err != nil {
		return err
	}

	for _, v := range s.FileList {
		if !strings.HasSuffix(v.FileName, ".go") {
			continue
		}

		filePath := strings.Replace(v.FilePath, "/", "_", -1)
		file := fmt.Sprintf("%v/%v", dir, filePath)

		b, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		v.FileContent = string(b)
	}

	return nil
}

func (s *Generator) gofmt(dir string) error {
	c := exec.Command("gofmt", "-w", dir)
	err := c.Run()
	if err != nil {
		log.Printf("gofmt err %v", err)
		return err
	}

	return nil
}
