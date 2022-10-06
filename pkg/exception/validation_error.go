package exception

import "fmt"

type ValidationError struct {
	Message string              `json:"message"`
	Details []ValidationDetails `json:"details"`
}

type ValidationDetails struct {
	Field       string `json:"field"`
	Value       string `json:"value"`
	Constraint  string `json:"constraint"`
	Description string `json:"description"`
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("%s: %v", v.Message, v.Details)
}
