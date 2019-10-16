package controller

import (
	"bytes"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

type fakeModel struct {
}

func (f *fakeModel) RetrieveBirthdayMessage(name string) (message string, err error) {
	switch name {
	case "":
		fallthrough
	case "Max Mustermann":
		return "", pgx.ErrNoRows
	case "Erika Mustermann":
		return "", pgx.ErrTxClosed
	case "Phil Connors":
		return "Happy birthday, Phil Connors!", nil
	default:
		return fmt.Sprintf("Hello, %s! Your birthday is in 100 days", name), nil
	}

	return "", nil
}

func (f *fakeModel) StoreBirthday(name string, dateOfBirth time.Time) error {
	switch name {
	case "":
		fallthrough
	case "Erika Mustermann":
		return pgx.ErrTxClosed
	default:
		return nil
	}
}


func _e(text string) string{
	return fmt.Sprintf(`{"Error":"%s"}`, text)
}

func _m(text string) string {
	return fmt.Sprintf(`{"Message":"%s"}`, text)
}


func TestTellBirthday(t *testing.T) {
	ctl := NewController(&fakeModel{})
	router := httprouter.New()
	router.Handler(http.MethodGet, "/hello", http.HandlerFunc(ctl.TellBirthday))
	log.SetLevel(log.FatalLevel)

	var tests = []struct{
		param string
		code int
		body string
	}{
		{"?name=Alex", http.StatusOK,
			_m("Hello, Alex! Your birthday is in 100 days")},

		{"?name="+url.QueryEscape("Max Mustermann"),
			http.StatusNotFound, _e("user not found: Max Mustermann")},

		{"?name="+url.QueryEscape("Erika Mustermann"),
			http.StatusInternalServerError, ""},

		{"?name="+url.QueryEscape("Phil Connors"),
			http.StatusOK, _m("Happy birthday, Phil Connors!")},

		{"?name=",
			http.StatusUnprocessableEntity, _e("missing name")},
	}

	for _ ,tt := range tests {
		req, err := http.NewRequest("GET", fmt.Sprintf("/hello%s", tt.param), nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set(HTTPHeaderAccept, MIMETypeJSON)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		if status := rr.Code; status != tt.code {
			t.Errorf("handler returned wrong status code: got %v, want %v", status, tt.code)
		}
		// JSON Encoder adds a trailing newline
		if body := strings.TrimSpace(rr.Body.String()); body != tt.body {
			t.Errorf("handler returned unexpected body: got '%v', want '%v'", body, tt.body)
		}
	}
}

func TestRecordBirthday(t *testing.T) {
	ctl := NewController(&fakeModel{})
	router := httprouter.New()
	router.Handler(http.MethodPut, "/hello/:name", http.HandlerFunc(ctl.RecordBirthday))
	log.SetLevel(log.FatalLevel)

	var tests = []struct{
		url string
		key string
		value string
		code int
		body string
	}{
		{ "/"+url.PathEscape("Erika Mustermann"),"dateOfBirth",
			"1982-12-13", http.StatusInternalServerError, "",
		},
		{ "/" + url.PathEscape("Max Mustermann"), "dateOfBirth",
			"1982-01-01", http.StatusCreated, "",
		},
		{"/" + url.PathEscape("John Doe"), "birthDate",
			"1982-03-03", http.StatusUnprocessableEntity, _e("missing date of birth"),
		},
		{
			"/" + url.PathEscape("Joan Doe"), "dateOfBirth",
			"2019-02-29", http.StatusUnprocessableEntity,
			_e("invalid data of birth: 2019-02-29, should be in a format YYYY-MM-DD"),
		},
		{
			"/" + url.PathEscape("Jan Jansen"), "dateOfBirth",
			"2019-29-01", http.StatusUnprocessableEntity,
			_e("invalid data of birth: 2019-29-01, should be in a format YYYY-MM-DD"),
		},
		{
			"/" + url.PathEscape("Nguyễn Văn A"), "dateOfBirth",
			"0001-01-01", http.StatusCreated,
			"",
		},


	}

	for _, tt := range tests {
		data := []byte(fmt.Sprintf(`{"%s":"%s"}`, tt.key, tt.value))
		urlPart := fmt.Sprintf("/hello%s", tt.url)
		req, err := http.NewRequest("PUT", urlPart, bytes.NewBuffer(data))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set(HTTPHeaderAccept, MIMETypeJSON)
		req.Header.Set(HTTPHeaderContentType, MIMETypeJSON)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != tt.code {
			t.Errorf("handler returned wrong status code: got %v, want %v", status, tt.code)
		}
		// JSON Encoder adds a trailing newline
		if body := strings.TrimSpace(rr.Body.String()); body != tt.body {
			t.Errorf("handler returned unexpected body: got '%v', want '%v'", body, tt.body)
		}
	}
}