package auth

type InvalidParam struct {
	s string
}

func (e *InvalidParam) Error() string {
	return e.s
}
