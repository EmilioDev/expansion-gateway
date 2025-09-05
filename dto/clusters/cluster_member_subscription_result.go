package clusters

type ClusterMemberSubscriptionResult struct {
	ServerID       int64 // the id of this server in the cluster
	HealthyTimeout int64 // the time after a check that the server allows before declaring this member as unhealthy
}
