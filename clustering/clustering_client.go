package clustering

type ClusteringClient struct {
	ClusterNode
}

func (cluster ClusteringClient) IsServer() bool {
	return true
}
