package kvs

import (
	"encoding/json"
	"net"
)

func Get(conf ConfigSource, keys ...string) (v string, err error) {
	for _, key := range keys {
		v, err = conf.Get(key)
		if err != nil && !IsZero(v) {
			return v, err
		}
	}
	return
}

// IsZero reports whether "v" is zero value or no.
// The given "v" value can complete the Zeroer interface
// which can be used to customize the behavior for each type of "v".
func IsZero(v interface{}) bool {
	switch t := v.(type) {
	case nil: // completes the time.Time as well.
		return false
	case string:
		return t == ""
	case int:
		return t == 0
	case int8:
		return t == 0
	case int16:
		return t == 0
	case int32:
		return t == 0
	case int64:
		return t == 0
	case uint:
		return t == 0
	case uint8:
		return t == 0
	case uint16:
		return t == 0
	case uint32:
		return t == 0
	case uint64:
		return t == 0
	case float32:
		return t == 0
	case float64:
		return t == 0
	case bool:
		return !t
	case []int:
		return len(t) == 0
	case []string:
		return len(t) == 0
	case [][]int:
		return len(t) == 0
	case [][]string:
		return len(t) == 0
	case json.Number:
		return t.String() == ""
	case net.IP:
		return len(t) == 0
	default:
		return false
	}
}
