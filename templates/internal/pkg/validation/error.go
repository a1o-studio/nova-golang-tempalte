package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field  string `json:"field"`  // 字段名
	Reason string `json:"reason"` // 错误信息
}

func Descriptive(errs validator.ValidationErrors) []ValidationError {
	res := []ValidationError{}

	for _, f := range errs {
		err := f.ActualTag()
		if f.Param() != "" {
			err = fmt.Sprintf("%s=%s", err, f.Param())
		}
		res = append(res, ValidationError{Field: f.Field(), Reason: err})
	}

	return res
}
