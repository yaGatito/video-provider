package common

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"video-provider/video-service/policy"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type serviceErrorResponse struct {
	Message string `json:"msg"`
}

// parseUrlValues parses URL query parameters.
func ParseUrlValues(query string) (url.Values, error) {
	if len(query) > policy.UrlMaxLen {
		return nil, &Error{
			Code:    http.StatusBadRequest,
			Message: "too large url",
		}
	}
	urlValues, err := url.ParseQuery(query)
	if err != nil {
		return nil, &Error{
			Err:     err,
			Code:    http.StatusBadRequest,
			Message: "unparsable url values",
			Details: query,
		}
	}
	return urlValues, nil
}

// parseIntsUrlParams parses integer parameters from the URL query.
func ParseIntsUrlParams(
	values url.Values,
	params ...string,
) ([]int32, error) {
	res := make([]int32, len(params))

	for i, param := range params {
		val, err := strconv.ParseInt(values.Get(param), 10, 32)
		if err != nil {
			return nil, &Error{
				Err:     err,
				Code:    http.StatusBadRequest,
				Message: "unparsable url param (int): " + param,
				Details: values}
		}
		res[i] = int32(val)
	}

	return res, nil
}

// parseStringsUrlParams parses string parameters from the URL query.
func ParseStringsUrlParams(
	values url.Values,
	params ...string,
) ([]string, error) {
	res := make([]string, len(params))

	for i, param := range params {
		value := values.Get(param)
		if len(value) == 0 {
			return nil, &Error{
				Code:    http.StatusBadRequest,
				Message: "empty value for " + param}
		}

		value, err := url.QueryUnescape(value)
		if err != nil {
			return nil, &Error{
				Err:     err,
				Code:    http.StatusBadRequest,
				Message: "failed to unescape url param: " + param}
		}

		res[i] = value
	}

	return res, nil
}

// pathVarHandler extracts a path variable and parses it as a UUID.
func PathVarHandler(
	r *http.Request,
	varName string,
) (uuid.UUID, error) {
	val, ok := mux.Vars(r)[varName]
	if !ok {
		return uuid.Nil, &Error{
			Code:    http.StatusBadRequest,
			Message: "path var not specified: " + varName}
	}
	res, err := uuid.Parse(val)
	if err != nil {
		return uuid.Nil, &Error{
			Err:     err,
			Code:    http.StatusBadRequest,
			Message: "unparsable ID: " + varName,
			Details: val}
	}

	return res, nil
}

// extractUrlVarString extracts and unescapes a string parameter from the URL query.
func ExtractUrlVarString(
	values url.Values,
	paramName string,
) (string, error) {
	value := values.Get(paramName)
	if len(value) == 0 {
		return "", ErrEmptyValue
	}
	value, err := url.QueryUnescape(value)
	if err != nil {
		return "", &Error{
			Err:     err,
			Code:    http.StatusBadRequest,
			Message: "failed to unescape url param: " + paramName}
	}

	return value, nil
}

// writeResponse writes a JSON response with the specified HTTP status code.
func WriteResponse(w http.ResponseWriter, log *Logger, val any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(val)
	if err != nil {
		log.Error("Error encoding response body", err)
	}
}

// writeErrorResponse writes an error response in JSON format.
func WriteErrorResponse(w http.ResponseWriter, log *Logger, vErr error) {
	w.Header().Set("Content-Type", "application/json")

	resp := serviceErrorResponse{}

	switch vErr := vErr.(type) {
	case *Error:
		log.Error(vErr.Message, vErr.Err)
		w.WriteHeader(int(vErr.Code))
		resp.Message = vErr.Message

	case validator.ValidationErrors:
		log.Debug("Validation request body error: " + vErr[0].Error())
		w.WriteHeader(http.StatusBadRequest)
		resp.Message = "invalid field: " + vErr[0].Field()

	case error:
		log.Error("Fallback error", vErr)
		w.WriteHeader(http.StatusInternalServerError)
		resp.Message = "video-provider error"
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Debug("Error encoding error response body:" + err.Error())
	}
}
