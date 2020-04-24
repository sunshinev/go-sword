package model

import (
	"go-sword/config"
	"io/ioutil"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sunshinev/db2struct"
)

type ModelCreate struct {
	Columns     []string
	Struc       []byte
	PackageName string
	StructName  string
	TableName   string
	// Files need to generate
	FileList []*FileInstance
}

type FileInstance struct {
	FileName    string
	FilePath    string
	FileContent []byte
}

// Entry
func (m *ModelCreate) ParseTable(c *config.Db, table string) {

	// Use db2struct (https://github.com/Shelnutt2/db2struct)
	columnDataTypes, err := db2struct.GetColumnsFromMysqlTable(c.User, c.Password, c.Host, c.Port, c.Database, table)
	if err != nil {
		panic(err.Error())
	}

	structName := strings.Replace(table, "_", " ", -1)
	structName = strings.Title(structName)
	structName = strings.Replace(structName, " ", "", -1)

	// Set columns
	for name := range *columnDataTypes {
		m.Columns = append(m.Columns, name)
	}

	struc, err := db2struct.Generate(*columnDataTypes, table, structName, table, true, true, false)
	if err != nil {
		panic(err.Error())
	}

	// Set TableName
	m.TableName = table
	// Set PackageName
	m.PackageName = table
	// Set StructName
	m.StructName = structName
	// Set Content
	m.Struc = struc
}

func (m *ModelCreate) Preview() {
	// Model
	var modelFile = &FileInstance{
		FilePath:    strings.Join([]string{"model", m.TableName, "model.go"}, string(os.PathSeparator)),
		FileName:    "model.go",
		FileContent: m.createModelContent(),
	}

	m.FileList = append(m.FileList, modelFile)

	// Controller
	var controllerFile = &FileInstance{
		FilePath:    strings.Join([]string{"controller", m.TableName, m.TableName + ".go"}, string(os.PathSeparator)),
		FileName:    m.TableName + ".go",
		FileContent: m.createControllerContent(),
	}

	m.FileList = append(m.FileList, controllerFile)

	// Html
	// create.html
	var listHtmlFile = &FileInstance{
		FilePath:    strings.Join([]string{"view", m.TableName, "create.html"}, string(os.PathSeparator)),
		FileName:    "create.html",
		FileContent: m.createListHtml(),
	}

	m.FileList = append(m.FileList, listHtmlFile)
}

// Create model file content
func (m *ModelCreate) createModelContent() []byte {
	// Modify m.Struc & add import
	if strings.Contains(string(m.Struc), "time.Time") {
		str := "package " + m.PackageName
		newStr := str + `
import "time"
`
		return []byte(strings.Replace(string(m.Struc), str, newStr, 1))
	}

	return m.Struc
}

// Create model file content
func (m *ModelCreate) createControllerContent() []byte {
	// Read stub
	file, err := os.Open("stub/controller/controller.stub")
	if err != nil {
		panic(err.Error())
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err.Error())
	}

	// replace
	packageName := m.TableName
	modelStruct := m.TableName + "." + m.StructName
	importModel := "go-sword/gosword/model/" + m.TableName

	content := string(data)

	content = strings.ReplaceAll(content, "<<package_name>>", packageName)
	content = strings.ReplaceAll(content, "<<model_struct>>", modelStruct)
	content = strings.ReplaceAll(content, "<<import_model>>", importModel)

	return []byte(content)
}

// Create model file content
func (m *ModelCreate) createListHtml() []byte {
	// Read stub
	file, err := os.Open("stub/controller/controller.stub")
	if err != nil {
		panic(err.Error())
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err.Error())
	}

	// replace
	packageName := m.TableName
	modelStruct := m.TableName + "." + m.StructName
	importModel := "go-sword/gosword/model/" + m.TableName

	content := string(data)

	content = strings.ReplaceAll(content, "<<package_name>>", packageName)
	content = strings.ReplaceAll(content, "<<model_struct>>", modelStruct)
	content = strings.ReplaceAll(content, "<<import_model>>", importModel)

	return []byte(content)
}
