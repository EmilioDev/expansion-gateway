package dispatchers

type Reciver[T any] interface {
	GetShard(index int) <-chan T // returns the shard at the given index
	ShardCount() int             // number of shards this receiver has
	TotalPending() int           // number of messages still pending to be processed
}
