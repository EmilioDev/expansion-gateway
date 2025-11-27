package results

type ClustersSubscriptionRequestBody struct {
	SubscriptionID      int64  // the id of this subscription
	Challenge           []byte // the challenge to check this subscription
	GatewayPath         string // the path to the gateway where the client should be redirected
	SessionEphemeralKey []byte // the public key to be used for the client to generate ths secret key
}

func (body *ClustersSubscriptionRequestBody) ToClusterUserSubscriptionResult() *ClusterUserSubscriptionResult {
	return &ClusterUserSubscriptionResult{
		SubscriptionID:      body.SubscriptionID,
		Challenge:           body.Challenge,
		GatewayPath:         body.GatewayPath,
		SessionEphemeralKey: body.SessionEphemeralKey,
	}
}
