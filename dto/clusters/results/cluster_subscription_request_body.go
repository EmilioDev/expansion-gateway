package results

type ClustersSubscriptionRequestBody struct {
	SubscriptionID int64   // the id of this subscription
	Challenge      []int32 // the challenge to check this subscription
}
