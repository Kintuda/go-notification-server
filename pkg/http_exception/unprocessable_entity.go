package http_exception

import (
	"fmt"

	"github.com/Kintuda/notification-server/pkg/exception"
	"github.com/go-playground/validator"
)

type UnprocessableEntity struct {
	Message string                        `json:"message"`
	Errors  []exception.ValidationDetails `json:"errors"`
}

func NewUnprocessableEntity() *UnprocessableEntity {
	return &UnprocessableEntity{
		Message: "validation error",
		Errors:  make([]exception.ValidationDetails, 0),
	}
}

func (u *UnprocessableEntity) AddError(vl exception.ValidationDetails) {
	u.Errors = append(u.Errors, vl)
}

func (u *UnprocessableEntity) AddErrorFromField(f validator.FieldError) {
	var vl exception.ValidationDetails
	vl.Field = exception.LowerFirstChar(f.Field())
	vl.Constraint = f.ActualTag()
	vl.Value = fmt.Sprintf("%v", f.Value())
	vl.Description = exception.FormatMessage(f.ActualTag())
	u.Errors = append(u.Errors, vl)
}
