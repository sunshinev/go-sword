package route

import (
	"go-sword/controller/render"
	"go-sword/gosword/controller/user"
	"net/http"

	"github.com/jinzhu/gorm"
)

func Register(db *gorm.DB) {
	// 静态路由
	http.Handle("/go_sword_public/", http.StripPrefix("/go_sword_public/", http.FileServer(http.Dir("resource/web/base/dist"))))
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("resource/web/base/dist"))))
	http.HandleFunc("/render", render.Render)

	http.HandleFunc("/api/user/list", user.List(db))
	http.HandleFunc("/api/user/delete", user.Delete(db))
	http.HandleFunc("/api/user/detail", user.Detail(db))
}
