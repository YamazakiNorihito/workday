package validator

import (
	"context"
	"fmt"
	"strings"

	"github.com/YamazakiNorihito/workday/cmd/rss/lambda/api/shared/validation_error"
	"github.com/go-playground/validator"
)

func Validate(ctx context.Context, s interface{}) error {
	validate := validator.New()
	errMap := make(map[string]string)

	if err := validate.StructCtx(ctx, s); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := err.Field()
			tag := err.Tag()
			param := err.Param()

			var values string
			switch tag {
			case "min", "max":
				values = fmt.Sprintf("value must be %s %s", tag, param)
			case "oneof":
				values = fmt.Sprintf("value must be one of [%s]", strings.ReplaceAll(param, " ", ", "))
			default:
				if param == "" {
					values = "invalid value"
				} else {
					values = param
				}
			}

			message := fmt.Sprintf("%s is %s: %s", fieldName, tag, values)
			errMap[fieldName] = message
		}
		if len(errMap) > 0 {
			return validation_error.New(errMap)
		}
	}
	return nil
}
