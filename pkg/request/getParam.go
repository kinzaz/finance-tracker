package request

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func GetParam[T any](r *http.Request, param string) (T, error) {
	var zero T
	stringParam := r.PathValue(param)

	if stringParam == "" {
		return zero, fmt.Errorf("параметр %v не найден", param)
	}

	var result any

	switch any(*new(T)).(type) {
	case int:
		parsed, err := strconv.Atoi(stringParam)
		if err != nil {
			return zero, fmt.Errorf("некорректное значение %v", param)
		}
		result = parsed

	case uint:
		parsed, err := strconv.ParseUint(stringParam, 10, 32)
		if err != nil {
			return zero, fmt.Errorf("некорректное значение %v", param)
		}
		result = uint(parsed)

	case string:
		result = stringParam

	default:
		return zero, errors.New("неподдерживаемый тип")
	}

	finalResult, ok := result.(T)
	if !ok {
		return zero, errors.New("не удалось привести к нужному типу")
	}

	return finalResult, nil
}
