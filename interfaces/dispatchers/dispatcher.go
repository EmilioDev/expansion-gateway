package dispatchers

type Dispatcher[T any] interface {
	Dispatch(pkt T)
}
