package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/smallbatch-apps/earnsmart-api/errs"
	"github.com/smallbatch-apps/earnsmart-api/middleware"
	"github.com/smallbatch-apps/earnsmart-api/models"
	"github.com/smallbatch-apps/earnsmart-api/services"
	"github.com/smallbatch-apps/earnsmart-api/utils"
)

type AllocationController struct {
	services *services.Services
}

func NewAllocationController(services *services.Services) *AllocationController {
	return &AllocationController{services}
}

func (c *AllocationController) ListAllocations(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		errs.UnauthorisedError(w, "Unable to get user context")
		return
	}
	allocations, err := c.services.Allocation.GetAllocations(userID)
	if err != nil {
		utils.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	utils.RespondOk(w, "allocations", allocations)
}

func (c *AllocationController) UpdateAllocation(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.RespondError(w, err, http.StatusUnauthorized)
		return
	}

	idString := r.PathValue("id")
	id, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		utils.RespondError(w, err, http.StatusBadRequest)
		return
	}

	var payload models.AllocationPlan
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		errs.InvalidPayloadError(w)
		return
	}

	allocation, err := c.services.Allocation.GetAllocation(id)
	if err != nil {
		utils.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	if !allocation.IsOwner(userID) {
		utils.RespondError(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	updatedAllocation, err := c.services.Allocation.UpdateAllocation(&allocation, payload)
	if err != nil {
		utils.RespondError(w, err, http.StatusInternalServerError)
	}
	utils.RespondOk(w, "allocation", updatedAllocation)
}

func (c *AllocationController) CreateAllocation(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.RespondError(w, err, http.StatusUnauthorized)
		return
	}

	var payload models.AllocationPlan
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		errs.InvalidPayloadError(w)
		return
	}
	payload.UserID = userID

	allocation, err := c.services.Allocation.CreateAllocation(&payload)
	if err != nil {
		utils.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	utils.RespondOk(w, "allocation", allocation)
}
