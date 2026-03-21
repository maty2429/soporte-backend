package handlers

import (
	stdjson "encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func init() {
	// Registra una función que le dice al validador que use el nombre del tag
	// json (o form) en lugar del nombre del campo Go.
	// Efecto: los mensajes de error muestran "nombre" en vez de "Nombre".
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "" || name == "-" {
				name = strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
			}
			if name == "-" {
				return ""
			}
			return name
		})
	}
}

// formatValidationError convierte errores de binding/validación de gin
// a mensajes legibles sin exponer nombres internos de structs ni tags.
func formatValidationError(err error) string {
	// JSON malformado
	var syntaxErr *stdjson.SyntaxError
	if errors.As(err, &syntaxErr) {
		return "request body contains malformed JSON"
	}

	// Tipo incorrecto (ej: string donde se espera número)
	var typeErr *stdjson.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		return fmt.Sprintf("field '%s' has an invalid type", typeErr.Field)
	}

	// Errores de validación (required, min, max, oneof, etc.)
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		msgs := make([]string, 0, len(ve))
		for _, fe := range ve {
			msgs = append(msgs, fieldErrorMsg(fe))
		}
		return strings.Join(msgs, "; ")
	}

	return "invalid request"
}

func fieldErrorMsg(fe validator.FieldError) string {
	field := fe.Field()
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("'%s' is required", field)
	case "gt":
		return fmt.Sprintf("'%s' must be greater than %s", field, fe.Param())
	case "gte":
		return fmt.Sprintf("'%s' must be greater than or equal to %s", field, fe.Param())
	case "lt":
		return fmt.Sprintf("'%s' must be less than %s", field, fe.Param())
	case "lte":
		return fmt.Sprintf("'%s' must be less than or equal to %s", field, fe.Param())
	case "min":
		return fmt.Sprintf("'%s' must have at least %s characters", field, fe.Param())
	case "max":
		return fmt.Sprintf("'%s' must have at most %s characters", field, fe.Param())
	case "oneof":
		return fmt.Sprintf("'%s' must be one of: %s", field, strings.ReplaceAll(fe.Param(), " ", ", "))
	case "email":
		return fmt.Sprintf("'%s' must be a valid email", field)
	case "url":
		return fmt.Sprintf("'%s' must be a valid URL", field)
	default:
		return fmt.Sprintf("'%s' is invalid (%s)", field, fe.Tag())
	}
}
