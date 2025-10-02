package layers

import (
	"expansion-gateway/dto/clusters/results"
	"expansion-gateway/enums"
	"expansion-gateway/interfaces/errorinfo"
)

type Layer2Follower interface {
	Layer2
	GenerateUserSubscription(
		userID,
		requestedSessionID int64,
		clientType enums.ClientType,
		clientVersion byte,
		encryption enums.EncryptionAlgorithm,
		protocolVersion enums.ProtocolVersion,
		sessionResume bool) (*results.ClustersSubscriptionRequestBody, errorinfo.GatewayError)
}
