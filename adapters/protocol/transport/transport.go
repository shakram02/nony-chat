package transport

type Transport[T any] interface {
	Read() (p T, err error)
	Write(p T) (err error)
	Close() error
}
