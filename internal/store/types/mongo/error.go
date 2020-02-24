package mongo

type OIDTypeError struct {
	s string
}

func (e *OIDTypeError) Error() string {
	return "Invalid OID"
}

// StoreEmpty implementation of store being empty
type StoreEmpty struct {
	s string
}

func (e *StoreEmpty) Error() string {
	return "Store is nil"
}

type connectError struct {
	s string
}

func (e *connectError) Error() string {
	return e.s
}
