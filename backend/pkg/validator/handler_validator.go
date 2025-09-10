package validator

import (
	"backend-service/internal/adapter/handler/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func HandleValidationError(c *gin.Context, err error, trans ut.Translator) {
	res := response.DefaultResponse{
		Success: false,
		Code:    http.StatusBadRequest,
		Data:    nil,
	}

	switch e := err.(type) {
	case validator.ValidationErrors:
		errMap := map[string][]string{}
		for _, fieldErr := range e {
			field := strings.ToLower(fieldErr.Field())
			msg := fieldErr.Translate(trans)
			errMap[field] = append(errMap[field], msg)
		}
		res.Message = errMap

	default:
		res.Message = err.Error()
	}

	c.JSON(http.StatusBadRequest, res)
}
