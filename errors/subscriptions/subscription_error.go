package subscriptions

import (
	"expansion-gateway/errors"
	"expansion-gateway/interfaces/errorinfo"
	"fmt"
)

type SubscriptionError struct {
	errors.BaseError
	subscription string
}

func (err *SubscriptionError) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return err
}

func (err *SubscriptionError) Error() string {
	return fmt.Sprintf("Error with subscription %s", err.subscription)
}

func CreateSubscriptionError(file, description, subscription string, index uint16, errorCode byte) *SubscriptionError {
	return &SubscriptionError{
		BaseError:    errors.CreateBaseError(file, description, index, errorCode),
		subscription: subscription,
	}
}
