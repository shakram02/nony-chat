package adapter

type Adapter[THigher any, TLower any] interface {
	CanReceive() bool
	Receive([]TLower) THigher
	Send(THigher) []TLower
}
