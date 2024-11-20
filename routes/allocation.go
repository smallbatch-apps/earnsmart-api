package routes

import (
	"net/http"

	"github.com/smallbatch-apps/earnsmart-api/controllers"
)

func RegisterAllocationRoutes(mux *http.ServeMux, controller *controllers.AllocationController) {
	mux.HandleFunc("GET /allocations", controller.ListAllocations)
	mux.HandleFunc("POST /allocations", controller.CreateAllocation)
	mux.HandleFunc("PATCH /allocations/{id}", controller.UpdateAllocation)
}
