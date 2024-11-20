package routes

import (
	"net/http"

	"github.com/smallbatch-apps/earnsmart-api/controllers"
)

func RegisterUserRoutes(mux *http.ServeMux, controller *controllers.UserController) {

	mux.HandleFunc("GET /user", controller.GetUser)
	mux.HandleFunc("DELETE /user/me", controller.DeleteUser)
	mux.HandleFunc("PATCH /user/me", controller.UpdateUser)

	mux.HandleFunc("POST /auth/logout", controller.LogOut)
}
