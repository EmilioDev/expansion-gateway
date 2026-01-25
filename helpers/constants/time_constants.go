package constants

import "time"

const CLUSTER_REQUEST_TIMEOUT time.Duration = time.Second * 25      // timeout of each request to another cluster member
const CLUSTER_MEMBER_INTERVAL_BETWEEN_CONNECTIONS = time.Second * 5 // interval between each connection attempt
