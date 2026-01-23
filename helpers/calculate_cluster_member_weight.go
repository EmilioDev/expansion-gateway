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

	return float32(math.Max(float64(messages), 1))/float32(math.Max(float64(activeSessions), 1)) + (cpuUsage+ramUsage)/2
}
