package routes

import (
	"net/http"

	"github.com/smallbatch-apps/earnsmart-api/controllers"
)

func RegisterSettingRoutes(mux *http.ServeMux, controller *controllers.SettingController) {
	mux.HandleFunc("GET /settings", controller.ListSettings)
	mux.HandleFunc("POST /settings", controller.AddSetting)
	mux.HandleFunc("GET /settings/{id}", controller.GetSetting)
	mux.HandleFunc("PATCH /settings/{id}", controller.EditSetting)
	mux.HandleFunc("DELETE /settings/{id}", controller.DeleteSetting)
}
