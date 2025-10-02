package controllers

import "expansion-gateway/clustering"

type Layer2Follower struct {
	*Layer2Core
	clusterServer *clustering.ClusteringFollower
}
