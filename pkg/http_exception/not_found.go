package http_exception

type NotFound struct {
	Message  string `json:"message"`
	Resource string `json:"resource"`
}

func NewNotFound(resource string) *NotFound {
	return &NotFound{
		Message:  "resource not found",
		Resource: resource,
	}
}
