package request

import (
	"net/http"
	"strconv"
	"strings"
)

func ParseIntQueryParam(r *http.Request, param string, defaultValue int) int {
	valueStr := r.FormValue(param)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func ParseCommaSeparatedQueryParamIds(r *http.Request, param string) []int {
	idsStr := r.FormValue(param)
	ids := strings.Split(idsStr, ",")
	var result []int
	for _, id := range ids {
		parsedId, err := strconv.Atoi(id)
		if err != nil {
			continue
		}
		result = append(result, parsedId)
	}
	return result
}
