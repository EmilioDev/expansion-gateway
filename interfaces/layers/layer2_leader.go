package layers

type Layer2Leader interface {
	Layer2
	MarkUserAsRedirected(int64)
}
