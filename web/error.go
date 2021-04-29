package web

type Error struct {
	Message  string   `json:"message"`
	Field    string   `json:"field,omitempty"`
	Children []*Error `json:"children,omitempty"`
}

func (err Error) Error() string {
	return err.Message
}

func (err *Error) AddChild(child *Error) {
	err.Children = append(err.Children, child)
}

func NewError(err error) *Error {
	switch v := err.(type) {
	case Error:
		return &v
	case *Error:
		return v
	default:
		return &Error{
			Message: v.Error(),
		}
	}
}
