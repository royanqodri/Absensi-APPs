package util

import (
	"strings"

	"github.com/grafana/pyroscope-go"
)

// Helper convert string profiles → []pyroscope.ProfileType
func ParseProfiles(s string) []pyroscope.ProfileType {
	var out []pyroscope.ProfileType
	for _, p := range strings.Split(s, ",") {
		switch strings.TrimSpace(p) {
		case "cpu":
			out = append(out, pyroscope.ProfileCPU)
		case "alloc":
			out = append(out,
				pyroscope.ProfileAllocObjects,
				pyroscope.ProfileAllocSpace,
			)
		case "inuse":
			out = append(out,
				pyroscope.ProfileInuseObjects,
				pyroscope.ProfileInuseSpace,
			)
		case "goroutine":
			out = append(out, pyroscope.ProfileGoroutines)
		case "mutex":
			out = append(out,
				pyroscope.ProfileMutexCount,
				pyroscope.ProfileMutexDuration,
			)
		}
	}
	return out
}
