package exception

import "fmt"

type ResourceNotFound struct {
	Identifier string `json:"id"`
	Resource   string `json:"resource"`
}

func (v *ResourceNotFound) Error() string {
	return fmt.Sprintf("%s not found", v.Resource)
}
