package internal

type InvalidParam struct {
	S string
}

func (e *InvalidParam) Error() string {
	return e.S
}
