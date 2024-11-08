package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/smallbatch-apps/earnsmart-api/models"
	"github.com/smallbatch-apps/earnsmart-api/schema"
	"github.com/smallbatch-apps/earnsmart-api/services"
)

type UserController struct {
	services *services.Services
}

type LogInPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUserController(services *services.Services) *UserController {
	return &UserController{services}
}

func (c *UserController) AddUser(w http.ResponseWriter, r *http.Request) {
	var payload schema.NewUserPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var user = models.User{
		Email:    payload.Email,
		Password: payload.Password,
		Name:     payload.Name,
	}
	err := c.services.User.CreateUser(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := schema.NewUserResponse{User: user, Status: "ok"}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "getting a user\n")
}

func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "updating a user\n")
}

func (c *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "deleting a user\n")
}

func (c *UserController) LogIn(w http.ResponseWriter, r *http.Request) {
	var payload LogInPayload
	error_string := "Invalid request payload"
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := c.services.User.FindUserByEmail(payload.Email)

	if err != nil {
		http.Error(w, error_string, http.StatusBadRequest)
		return
	}

	err = user.ComparePassword(payload.Password)

	if err != nil {
		http.Error(w, error_string, http.StatusBadRequest)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":    user.ID,
			"email": user.Email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		http.Error(w, "Error signing token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "bearer "+tokenString)

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (c *UserController) LogOut(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Logging out\n")
}
