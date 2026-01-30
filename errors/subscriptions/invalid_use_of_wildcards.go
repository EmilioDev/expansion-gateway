package subscriptions

import (
	"expansion-gateway/interfaces/errorinfo"
	"fmt"
)

type InvalidUseOfWildcards struct {
	*SubscriptionError
}

func (err *InvalidUseOfWildcards) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return err
}

func (err *InvalidUseOfWildcards) Error() string {
	return fmt.Sprintf("Error with subscription %s, it has wildcards used in a wrong way", err.subscription)
}

func CreateInvalidUseOfWildcardsError(file, subscription string, index uint16) *InvalidUseOfWildcards {
	return &InvalidUseOfWildcards{
		SubscriptionError: CreateSubscriptionError(file, "invalid use of wildcards in subscription", subscription, index, 25),
	}
}
