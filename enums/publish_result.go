package enums

type PublishResult byte

const (
	SUCCEED PublishResult = iota
	USER_NOT_REGISTERED_IN_SUBSCRIPTION
	INVALID_KEY
	ENCRYPTION_FAILED
	OUTDATED
)
