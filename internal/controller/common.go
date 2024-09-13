package controller

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"io"
	log "log/slog"
	"net/http"
)

func ReadBody(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		ErrorResponse(w, "invalid request body", http.StatusBadRequest)
		return []byte{}, nil
	}
	return body, err
}

func ParseJSONBody[T any](r *http.Request, w http.ResponseWriter) (*T, error) {
	var t T
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &t)
	if err != nil {
		return nil, err
	}

	var validate = validator.New()

	err = validate.Struct(t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func DecodeFormParams[T any](r *http.Request) (*T, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	params := new(T)
	var decoder = schema.NewDecoder()
	if err := decoder.Decode(params, r.Form); err != nil {
		return nil, err
	}

	return params, nil
}

func ParseForm(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid params", http.StatusBadRequest)
		return err
	}
	return nil
}

//func (tr *TenderController) decodeParams(w http.ResponseWriter, decoder *schema.Decoder, params interface{}) error {
//	if err := decoder.Decode(params); err != nil {
//		http.Error(w, "invalid params", http.StatusInternalServerError)
//		return err
//	}
//	return nil
//}

func HandleServiceError(w http.ResponseWriter, err error, statusCode int, errorMessage string) {
	if err != nil {
		log.Debug("err: %v", err.Error())
		ErrorResponse(w, errorMessage, statusCode)
		return
	}
}

func SendJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
