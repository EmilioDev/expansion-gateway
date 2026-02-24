package subscriptions

import (
	"expansion-gateway/interfaces/errorinfo"
	"fmt"
)

type InvalidCharactersInSubscription struct {
	*SubscriptionError
}

func (err *InvalidCharactersInSubscription) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return err
}

func (err *InvalidCharactersInSubscription) Error() string {
	return fmt.Sprintf("Error with subscription %s, it has the '@' or '.' characters, which are not allowed", err.subscription)
}

func CreateUseOfInvalidCharactersInSubscriptionError(file, subscription string, index uint16) *InvalidCharactersInSubscription {
	return &InvalidCharactersInSubscription{
		SubscriptionError: CreateSubscriptionError(file, "subscription with invalid characters", subscription, index, 26),
	}
}
