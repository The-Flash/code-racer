package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type ExecFiles struct {
	Path string `json:"path"`
	Body string `json:"body"`
}

type codeExecRequest struct {
	Language   string      `json:"language"`
	Entrypoint string      `json:"entrypoint"`
	Files      []ExecFiles `json:"files"`
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
	io.WriteString(w, fmt.Sprintf("%s %s %s", rBody.Language, rBody.Entrypoint, rBody.Files))
	//	response := &codeExecResponse{Language: rBody.Language, Code: rBody.Code}
	//	err = json.NewEncoder(w).Encode(response)

	//	if err != nil {
	//		log.Println(err)
	//	}
}
