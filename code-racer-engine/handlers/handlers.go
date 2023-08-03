package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type codeExecRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

type codeExecResponse struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

func CodeExecHandler(w http.ResponseWriter, r *http.Request) {
	var rBody codeExecRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &rBody)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := &codeExecResponse{Language: rBody.Language, Code: rBody.Code}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err)
	}
}
