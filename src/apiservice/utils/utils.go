package utils

import (
	"apiservice/Errs"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func RenderJson(w http.ResponseWriter, status int, res interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	// We don't have to write body, If status code is 204 (No Content)
	if status == http.StatusNoContent {
		return
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("ERROR: renderJson - %q\n", err)
	}
}

func ParseJSON(w http.ResponseWriter, params io.ReadCloser, data interface{}) bool {
	if params != nil {
		defer params.Close()
	}

	err := json.NewDecoder(params).Decode(data)
	if err == nil {
		return true
	}

	e := &Errs.AppErr{
		Message: "Invalid JSON",
		Err:     err.Error(),
	}

	RenderJson(w, http.StatusBadRequest, e)
	return false
}

func TotalRecordsInPage(pno, limit, count int) int {

	skip := (pno - 1) * limit

	switch {
	case skip == 0 && limit >= count:
		return count
	case skip == 0 && limit < count:
		return limit
	case skip > 0 && limit >= (count - skip) && (count - skip) > 0:
		return count - skip
	case skip > 0 && limit > (count - skip) && (count - skip) < 0:
		return 0
	case skip > 0 && limit < (count - skip):
		return limit
	}

	return 0
}