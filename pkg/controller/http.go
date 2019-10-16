package controller

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	HTTPHeaderAccept = "Accept"
	HTTPHeaderContentType = "Content-Type"

	MIMETypeJSON = "application/json"
	MIMETypeJSONProblem = "application/problem+json"
	MIMETYPEFormURLEncoded = "application/x-www-form-urlencoded"
	MIMETypePlainText = "text/plain"
	MimeTypeEverything = "*/*"

)

func hasMIMETYPE(header, mimeType string) bool {
	if header == "" {
		return false
	}
	for _, v := range strings.Split(header, ",") {
		if t, _, err := mime.ParseMediaType(v); err != nil {
			return false
		} else if t == mimeType || t == MimeTypeEverything {
			return true
		}
	}
	return false
}


func sendResponse(w http.ResponseWriter, r *http.Request, message string) {
	if hasMIMETYPE(r.Header.Get(HTTPHeaderAccept), MIMETypeJSON) {
		w.Header().Set(HTTPHeaderContentType, MIMETypeJSON)
		_ = json.NewEncoder(w).Encode(struct{ Message string }{message})
	} else {
		w.Header().Set(HTTPHeaderContentType, MIMETypePlainText)
		fmt.Fprint(w, message)
	}
}


func handleError(w http.ResponseWriter, r *http.Request, err *TypedError) {
	log.Errorf("error when processing request: %s", err.Error())
	switch err.kind {
	case ErrUnsupportedMediaType:
		w.WriteHeader(http.StatusUnsupportedMediaType)
	case ErrParsingForm:
		w.WriteHeader(http.StatusBadRequest)
	case ErrDecodingJSON:
		fallthrough
	case ErrMissingDateOfBirth:
		fallthrough
	case ErrInvalidDateOfBirthFormat:
		fallthrough
	case ErrMissingName:
		encodeErrorMessage(w, r, err, http.StatusUnprocessableEntity)
	case ErrMissingRecord:
		encodeErrorMessage(w, r, err, http.StatusNotFound)
	case ErrDatabaseError:
		fallthrough
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func encodeErrorMessage(w http.ResponseWriter, r *http.Request, err error, status int) {
	if hasMIMETYPE(r.Header.Get(HTTPHeaderAccept), MIMETypeJSON) {
		if err != nil {
			w.Header().Set(HTTPHeaderContentType, MIMETypeJSONProblem)
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(struct{ Error string }{err.Error()})
		}
	} else {
		w.Header().Set(HTTPHeaderContentType, MIMETypePlainText)
		w.WriteHeader(status)
		fmt.Fprint(w, err.Error())
	}
}
