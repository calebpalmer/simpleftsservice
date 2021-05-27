package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Status struct {
	Status string `json:"status"`
}

type GetStatusHandler struct{}

func (h *GetStatusHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/status" {
		http.NotFound(w, req)
		return
	}

	resp, err := json.Marshal(Status{"ok"})
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(resp))
}

func RegisterStatusHandlers(router *mux.Router) error {
	router.Handle("/status", &GetStatusHandler{}).Methods("GET")
	return nil
}
