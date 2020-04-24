package route

import (
	"go-sword/controller/render"
	"go-sword/gosword/controller/user"
	"net/http"

	"github.com/jinzhu/gorm"
)

func Register(db *gorm.DB) {
	// 静态路由
	http.Handle("/go_sword_public/", http.StripPrefix("/public/", http.FileServer(http.Dir("resource/web/base/dist"))))

	http.HandleFunc("/admin/render", render.Render)
	http.Handle("/admin/", http.StripPrefix("/admin/", http.FileServer(http.Dir("resource/web/base/dist"))))

	http.HandleFunc("/api/admin/users/list", user.List(db))
	http.HandleFunc("/api/admin/users/delete", user.Delete(db))
	http.HandleFunc("/api/admin/users/detail", user.Detail(db))
	//http.HandleFunc("/api/admin/users/batch_delete", user.List(db))
}
