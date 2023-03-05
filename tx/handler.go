package tx

type Handler interface {
	Handle(*Tx) error
}
