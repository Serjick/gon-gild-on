package golden

import (
	"flag"
	"reflect"
)

// UpdateAllower is for allow/deny golden file overwrite with actual data.
type UpdateAllower func() bool

//nolint:gochecknoglobals // flags must be defined before flag.Parse() call which executed by go test framework.
var forceUpdateFlag = flag.Bool("gon-gild-on.update", false, "force golden files update with actual data")

func init() { //nolint:gochecknoinits // no other way to avoid conflicts with third-party libs using same flag
	if flag.Lookup("update") == nil {
		flag.Bool("update", false, "force golden files update with actual data")
	}
}

// NewUpdateAllowerByFlag returns UpdateAllower allowing update
// only when `-gon-gild-on.update` or `-update` flag is specified and it is not `false`.
func NewUpdateAllowerByFlag() UpdateAllower {
	return func() bool {
		if *forceUpdateFlag {
			return true
		}

		var val any
		if getter, ok := flag.Lookup("update").Value.(flag.Getter); ok {
			val = getter.Get()
		}

		return reflect.DeepEqual(val, any(true))
	}
}
