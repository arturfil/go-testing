package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, wrap ...string) error {
    var out []byte
    if len(wrap) > 0 {
        wrapper := make(map[string]interface{})
        wrapper[wrap[0]] = data
        jsonBytes, err := json.Marshal(wrapper)
        if err != nil {
            return err
        }
        out = jsonBytes
    } else {
        jsonBytes, err := json.Marshal(data)
        if err != nil {
            return err
        }
        out = jsonBytes
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)

    _, err := w.Write(out)
    if err != nil {
        return err
    }
    return nil
}

func (app *application) errorJSON(w http.ResponseWriter, err error, status ...int) {
    statusCode := http.StatusBadRequest
    if len(status) > 0 {
        statusCode = status[0]
    }

    type jsonError struct {
        Message string `json:"message"` 
    }

    theError := jsonError {
        Message: err.Error(),
    }

    _ = app.writeJSON(w, statusCode, theError, "error")
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
    maxBytes := 1024 * 1024 // one megabyte
    r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
    dec := json.NewDecoder(r.Body)
    dec.DisallowUnknownFields()

    err := dec.Decode(data)
    if err != nil {
        return err 
    }

    err = dec.Decode(&struct{}{})
    if err != io.EOF {
        return errors.New("body must only contain a single JSON value")
    }
    return nil
}

