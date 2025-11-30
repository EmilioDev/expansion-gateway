package constants

import "time"

const CLUSTER_REQUEST_TIMEOUT time.Duration = time.Millisecond * 200 // timeout of each request to another cluster member
