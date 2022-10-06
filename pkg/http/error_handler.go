package http

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/Kintuda/notification-server/pkg/exception"
	"github.com/Kintuda/notification-server/pkg/http_exception"

	"github.com/sirupsen/logrus"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
			if errors.Is(err, io.EOF) {
				c.JSON(http.StatusBadRequest, gin.H{"message": "request body is empty"})
				c.Abort()
				return
			}

			switch err.Err.(type) {
			case validator.ValidationErrors:
				validationErrors := http_exception.NewUnprocessableEntity()
				errs := err.Err.(validator.ValidationErrors)

				for _, err := range errs {
					validationErrors.AddErrorFromField(err)
				}

				c.JSON(http.StatusUnprocessableEntity, validationErrors)
				c.Abort()
				return
			case *exception.ResourceNotFound:
				errs := err.Err.(*exception.ResourceNotFound)
				notFound := http_exception.NewNotFound(errs.Identifier)
				c.JSON(http.StatusNotFound, notFound)
				c.Abort()
			default:
				logrus.Error("ErrorHandler: internal server error", err.Err)
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				c.Abort()
			}
		}
	}
}
