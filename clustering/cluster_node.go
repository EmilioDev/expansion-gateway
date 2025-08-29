package clustering

import "expansion-gateway/config"

type ClusterNode struct {
	grpcCurrentServerPath string
}

func CreateBaseClusterNode(conf *config.Configuration) ClusterNode {
	return ClusterNode{
		grpcCurrentServerPath: conf.GetGrpcCurrentServerPath(),
	}
}
