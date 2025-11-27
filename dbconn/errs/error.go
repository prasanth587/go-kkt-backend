package errs

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-pg/pg"
)

// Error is to capture error
type Error struct {
	Code    int    `json:"errCode"`
	Message string `json:"message"`
	Err     error  `json:"error"`
	Module  string `json:"-"`
	IsDbErr bool   `json:"-"`
}

func (e Error) Error() string {
	e.Set()
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

// Stack will return the actual error
func (e Error) Stack() string {
	if e.Err == nil {
		return ""
	}
	log.Println("stack", e.Err.Error())
	return e.Err.Error()
}

// SetDbError handle db error set message for that
func (e *Error) SetDbError() {
	switch e.Err {
	case pg.ErrNoRows:
		e.Code = http.StatusNotFound
		e.Message = fmt.Sprintf("%s not found", e.Module)
		return
	}

	dbErr, ok := (e.Err).(pg.Error)
	if !ok {
		e.Code = http.StatusBadRequest
		if e.Message == "" {
			e.Message = fmt.Sprintf("DB Error: %s", e.Stack())
		}
		return
	}

	switch dbErr.Field('C') {
	case "23505":
		e.Code = http.StatusBadRequest
		e.Message = dbErr.Field('D')
	default:
		e.Code = http.StatusBadRequest
		if e.Message == "" {
			e.Message = fmt.Sprintf("DB Error: %s", e.Stack())
		}
	}
}

// Set set the code& respective message if not available
func (e *Error) Set() {

	if e.IsDbErr {
		e.SetDbError()
		return
	}

	if e.Code == 0 {
		e.Code = http.StatusBadRequest
	}

	if e.Message != "" {
		return
	}

	switch e.Code {
	case 404:
		e.Message = fmt.Sprintf("%s not found", e.Module)
	default:
		e.Message = fmt.Sprintf("%d: Error", e.Code)
	}
}

// MarshalJSON implements json marshaller
func (e Error) MarshalJSON() ([]byte, error) {
	type Alias Error
	return json.Marshal(&struct {
		Err string `json:"error"`
		Alias
	}{
		Err:   e.Stack(),
		Alias: (Alias)(e),
	})
}
