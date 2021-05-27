package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// writes an internal server error from an error.
func writeInternalServerError(w http.ResponseWriter, err error) {
	log.Print(err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	msg, _ := json.Marshal(map[string]string{"details": fmt.Sprintf("%s", err)})
	fmt.Fprint(w, string(msg))
}
