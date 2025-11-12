package results

type ClusterUserSubscriptionResult struct {
	Challenge      []byte // the challenge to be sent to the client
	SubscriptionID int64  // the id of the subscription to the user
	GatewayPath    string // the path to access the new gateway
}
