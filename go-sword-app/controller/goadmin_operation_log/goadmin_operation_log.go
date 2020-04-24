package goadmin_operation_log

import (
	"encoding/json"
	"go-sword/go-sword-app/model/goadmin_operation_log"
	"go-sword/response"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
)

func List(db *gorm.DB) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		sqls := []goadmin_operation_log.GoadminOperationLog{}
		db.Find(&sqls)

		ret := response.Ret{
			Code: http.StatusOK,
			Data: response.List{
				List: &sqls,
			},
		}

		d, err := json.Marshal(&ret)
		if err != nil {
			log.Printf("%v", err)
			return
		}

		_, err = writer.Write(d)
		if err != nil {
			log.Printf("%v", err)
			return
		}
	}
}

func Delete(db *gorm.DB) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {

		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Printf("%v", err)
			return
		}

		data := make(map[string]int)
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Printf("%v", err)
			return
		}

		sql := goadmin_operation_log.GoadminOperationLog{
			ID: data["id"],
		}

		db.Delete(&sql)

		ret := response.Ret{
			Code: http.StatusOK,
			Msg:  "success",
		}

		d, err := json.Marshal(&ret)
		if err != nil {
			log.Printf("%v", err)
			return
		}

		_, err = writer.Write(d)
		if err != nil {
			log.Printf("%v", err)
			return
		}
	}
}

func Create(db *gorm.DB) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Printf("%v", err)
			return
		}

		data := goadmin_operation_log.GoadminOperationLog{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Printf("%v", err)
			return
		}

		db.Create(&data)

		ret := response.Ret{
			Code: http.StatusOK,
			Msg:  "success",
		}

		d, err := json.Marshal(&ret)
		if err != nil {
			log.Printf("%v", err)
			return
		}

		_, err = writer.Write(d)
		if err != nil {
			log.Printf("%v", err)
			return
		}
	}
}

func Edit(db *gorm.DB) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Printf("%v", err)
			return
		}

		data := goadmin_operation_log.GoadminOperationLog{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Printf("%v", err)
			return
		}

		db.First(&data, data.ID)
		db.Save(&data)

		ret := response.Ret{
			Code: http.StatusOK,
			Msg:  "success",
			Data: &data,
		}

		d, err := json.Marshal(&ret)
		if err != nil {
			log.Printf("%v", err)
			return
		}

		_, err = writer.Write(d)
		if err != nil {
			log.Printf("%v", err)
			return
		}
	}
}

func Detail(db *gorm.DB) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Printf("%v", err)
			return
		}

		log.Printf("%s", body)

		params := make(map[string]int)
		err = json.Unmarshal(body, &params)
		if err != nil {
			log.Printf("%v", err)
			return
		}

		data := goadmin_operation_log.GoadminOperationLog{}

		db.First(&data, params["id"])

		ret := response.Ret{
			Code: http.StatusOK,
			Msg:  "success",
			Data: &data,
		}

		d, err := json.Marshal(&ret)
		if err != nil {
			log.Printf("%v", err)
			return
		}

		_, err = writer.Write(d)
		if err != nil {
			log.Printf("%v", err)
			return
		}
	}
}
