package errs

import (
	"log"
	"net/http"
)

func UnauthorisedError(w http.ResponseWriter, internalMessage string) {
	http.Error(w, "Unauthorised access", http.StatusUnauthorized)
	log.Printf("unauthorized access: %s", internalMessage)
}

func InternalError(w http.ResponseWriter, internalMessage string, message string) {
	http.Error(w, message, http.StatusInternalServerError)
	log.Printf("unauthorized access: %s", internalMessage)
}

func InvalidPayloadError(w http.ResponseWriter) {
	http.Error(w, "Unable to parse payload", http.StatusBadRequest)
	log.Print("Unable to parse payload")
}

func UserTokenNotValidError(w http.ResponseWriter) {
	http.Error(w, "Unauthorised access", http.StatusUnauthorized)
	log.Printf("User token is not valid")
}
