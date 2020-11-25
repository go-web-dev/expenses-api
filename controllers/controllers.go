package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/steevehook/expenses-rest-api/logging"
	"github.com/steevehook/expenses-rest-api/models"
)

// routeParam fetches params from context and converts it into julienschmidt/httprouter.Params struct
func routeParam(r *http.Request, name string) string {
	ctx := r.Context()
	psCtx := ctx.Value(httprouter.ParamsKey)
	ps, ok := psCtx.(httprouter.Params)

	if !ok {
		logging.Logger.Error("could not extract params from context")
		return ""
	}
	return ps.ByName(name)
}

// parseBody parses JSON request body
func parseBody(r *http.Request, v interface{}) error {
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logging.Logger.Error("could not read request body")
		return models.InvalidJSONError{
			Message: "could not read request body",
		}
	}

	err = json.Unmarshal(bs, v)
	switch err.(type) {
	case *json.UnsupportedTypeError, *json.UnsupportedValueError:
		return models.InvalidJSONError{
			Message: err.Error(),
		}
	}
	return nil
}

// parseIDParam parses id route param and validates it
func parseIDParam(r *http.Request) (string, error) {
	id := routeParam(r, idRouteParam)
	_, err := uuid.Parse(id)
	if err != nil {
		e := models.FormatValidationError{
			Message: fmt.Sprintf("invalid uuid: %s", id),
		}
		return "", e
	}
	return id, nil
}

// parseIDsParam parses ids route param and validates it
func parseIDsParam(r *http.Request) ([]string, error) {
	// approximately max 50 UUIDs in URL
	ids := strings.Split(routeParam(r, idsRouteParam), ",")
	for _, id := range ids {
		_, err := uuid.Parse(id)
		if err != nil {
			e := models.FormatValidationError{
				Message: fmt.Sprintf("invalid uuid: %s", id),
			}
			return []string{}, e
		}
	}
	return dedupe(ids), nil
}

// dedupe parses a list of strings and removes duplicates
func dedupe(values []string) []string {
	unique, res := map[string]string{}, make([]string, 0)
	for _, v := range values {
		unique[v] = v
	}
	for _, v := range unique {
		res = append(res, v)
	}
	return res
}
