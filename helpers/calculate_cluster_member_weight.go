package helpers

import "math"

func CalculateClusterMemberWeight(
	messages int64,
	activeSessions int32,
	cpuUsage,
	ramUsage float32,
	healthy bool,
) float32 {
	if !healthy {
		return math.MaxFloat32
	}

	return float32(messages)/float32(activeSessions) + (cpuUsage+ramUsage)/2
}
