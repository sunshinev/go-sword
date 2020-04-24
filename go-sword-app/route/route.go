// Sword will check if route file is created, if existed, Sword will modify it
// If you want to recreated the route,you should delete the file,and then use Sword generate again, or you can copy from the `stub/route/route.stub` file

// Do not modify the notes `----Route-begin----` or `----Route-end----` or `----Import----`

package route

import (
	"go-sword/controller/render"
	"net/http"

	"github.com/jinzhu/gorm"
	"go-sword/go-sword-app/controller/goadmin_permissions"
	"go-sword/go-sword-app/controller/goadmin_role_users"
	"go-sword/go-sword-app/controller/goadmin_site"
	// ----Import----
)

func Register(db *gorm.DB) {
	// Static file
	http.Handle("/go_sword_public/", http.StripPrefix("/go_sword_public/", http.FileServer(http.Dir("resource/web/base/dist"))))

	// Default index.html
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("resource/web/base/dist"))))

	// Render Vue html component
	http.HandleFunc("/render", render.Render)

	// ----Route-begin----

	// Route tag goadmin_role_users
	http.HandleFunc("/api/goadmin_role_users/list", goadmin_role_users.List(db))
	http.HandleFunc("/api/goadmin_role_users/delete", goadmin_role_users.Delete(db))
	http.HandleFunc("/api/goadmin_role_users/detail", goadmin_role_users.Detail(db))

	// Route tag goadmin_site
	http.HandleFunc("/api/goadmin_site/list", goadmin_site.List(db))
	http.HandleFunc("/api/goadmin_site/delete", goadmin_site.Delete(db))
	http.HandleFunc("/api/goadmin_site/detail", goadmin_site.Detail(db))

	// Route tag goadmin_permissions
	http.HandleFunc("/api/goadmin_permissions/list", goadmin_permissions.List(db))
	http.HandleFunc("/api/goadmin_permissions/delete", goadmin_permissions.Delete(db))
	http.HandleFunc("/api/goadmin_permissions/detail", goadmin_permissions.Detail(db))
	// ----Route-end----
}
