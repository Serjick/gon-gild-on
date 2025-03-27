package golden

import (
	"flag"
)

// UpdateAllower is for allow/deny golden file overwrite with actual data.
type UpdateAllower func() bool

//nolint:gochecknoglobals // flags must be defined before flag.Parse() call which executed by go test framework.
var forceUpdateFlag = flag.Bool("update", false, "force golden files update with actual data")

// NewUpdateAllowerByFlag returns UpdateAllower who allow update only if `-update` flag is specified and not `false`.
func NewUpdateAllowerByFlag() UpdateAllower {
	return func() bool {
		return forceUpdateFlag != nil && *forceUpdateFlag
	}
}
