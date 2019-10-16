package controller

import (
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/alexeyklyukin/rev4/pkg/db"
)

type birthDayData struct {
	DateOfBirth string `json:"dateOfBirth"`
}

type Controller struct {
	db db.Model
}

func NewController(db db.Model) *Controller {
	return &Controller{db}
}

func (ctl *Controller) RecordBirthday(w http.ResponseWriter, r *http.Request) {
	var (
		dateOfBirth string
		err error
		dateOfBirthParsed time.Time
	)

	err = r.ParseForm()
	if err != nil {
		handleError(w, r, NewTypedError(ErrParsingForm, fmt.Sprintf("couldn't parse input: %v", err)))
		return
	}

	contentType := r.Header.Get(HTTPHeaderContentType)
	if hasMIMETYPE(contentType, MIMETypeJSON) {
		decoder := json.NewDecoder(r.Body)
		var data birthDayData
		err = decoder.Decode(&data)
		if err != nil {
			handleError(w, r, NewTypedError(ErrDecodingJSON,fmt.Sprintf("could not decode JSON data: %v", err)))
			return
		}
		dateOfBirth = data.DateOfBirth
	} else if contentType == "" || hasMIMETYPE(contentType, MIMETYPEFormURLEncoded) {
		dateOfBirth = r.Form.Get("dateOfBirth")
	} else {
		handleError(w, r, NewTypedError(ErrUnsupportedMediaType,
			       fmt.Sprintf("unsupported Content-Type: %s", contentType)))
		return
	}

	if dateOfBirth == "" {
		handleError(w, r, NewTypedError(ErrMissingDateOfBirth, "missing date of birth"))
		return
	}
	dateOfBirthParsed, err = time.Parse("2006-01-02", dateOfBirth)
	if err != nil {
		handleError(w, r, NewTypedError(ErrInvalidDateOfBirthFormat,
			       fmt.Sprintf("invalid data of birth: %s, should be in a format YYYY-MM-DD", dateOfBirth)))
		return
	}
	ps := httprouter.ParamsFromContext(r.Context())
	name := ps.ByName("name")
	if name == "" {
		handleError(w, r, NewTypedError(ErrMissingName, "missing name"))
		return
	}

	err = ctl.db.StoreBirthday(name, dateOfBirthParsed)
	if err != nil {
		handleError(w, r, NewTypedError(ErrDatabaseError, fmt.Sprintf("database error: %v", err)))
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (ctl *Controller) TellBirthday(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		handleError(w, r, NewTypedError(ErrParsingForm, fmt.Sprintf("couldn't parse input: %v", err)))
		return
	}
	name := r.Form.Get("name")
	if name == "" {
		handleError(w, r, NewTypedError(ErrMissingName, "missing name"))
		return
	}
	msg, err := ctl.db.RetrieveBirthdayMessage(name)
	if err != nil {
		if err != pgx.ErrNoRows {
			handleError(w, r, NewTypedError(ErrDatabaseError, fmt.Sprintf("database error: %v", err)))
		} else {
			handleError(w, r, NewTypedError(ErrMissingRecord, fmt.Sprintf("user not found: %s", name)))
		}
		return
	}
	sendResponse(w, r, msg)
}
