package clustering

type ClusteringServer struct {
	ClusterNode
}

func (cluster ClusteringServer) IsServer() bool {
	return true
}
