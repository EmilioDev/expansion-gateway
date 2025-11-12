package results

type ClustersSubscriptionRequestBody struct {
	SubscriptionID int64   // the id of this subscription
	Challenge      []int32 // the challenge to check this subscription
	GatewayPath    string  // the path to the gateway where the client should be redirected
}
