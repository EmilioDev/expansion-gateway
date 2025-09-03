package enums

type ClusterMember_Kind byte

const (
	ClusterLeader ClusterMember_Kind = iota
	ClusterFollower
)
