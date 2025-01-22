package request

import (
	"net/http"
	"strconv"
)

func ParseIntQueryParam(r *http.Request, param string, defaultValue int) int {
	valueStr := r.FormValue(param)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
