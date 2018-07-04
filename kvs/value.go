package kvs

import (
    "strconv"
    "time"
    "strings"
    "fmt"
)

const (
    _DEPTH_VALUES = 99
)

//参考了go-ini/ini中的key.go源码，做了一些默认操作的修改
var kv_delimiters []rune

func init() {
    kv_delimiters = []rune{
        ',', '|', '，', ' ', ' ',
    }
}

// KeyValue represents a key under a section.
type KeyValue struct {
    key    string
    value  string
    Delims []rune
    err    error
}

// newKey simply return a key object with given Values.
func NewKeyValue(key, val string) *KeyValue {
    return NewKeyValueByDelims(key, val, kv_delimiters)
}

// newKey simply return a key object with given Values.
func NewKeyValueByStrDelims(key, val, delims string) *KeyValue {
    delimiters := kv_delimiters
    if delims != "" {
        delimiters = []rune(delims)
    }
    return NewKeyValueByDelims(key, val, delimiters)
}

// newKey simply return a key object with given Values.
func NewKeyValueByDelims(key, val string, delims []rune) *KeyValue {
    return &KeyValue{
        key:    key,
        value:  val,
        Delims: delims,
    }
}

// ValueMapper represents a mapping function for Values, e.g. os.ExpandEnv
type ValueMapper func(string) string

// ConfName returns ConfName of key.
func (k *KeyValue) Key() string {
    return k.key
}

// Value returns raw value of key for performance purpose.
func (k *KeyValue) Value() string {
    return k.value
}

// String returns string representation of value.
func (k *KeyValue) String() string {
    return k.value
}

// Validate accepts a validate function which can
// return modifed result as key value.
func (k *KeyValue) Validate(fn func(string) string) string {
    return fn(k.String())
}

// parseBool returns the boolean value represented by the string.
//
// It accepts 1, t, T, TRUE, true, True, YES, yes, Yes, y, ON, on, On,
// 0, f, F, FALSE, false, False, NO, no, No, n, OFF, off, Off.
// Any other value returns an error.
func parseBool(str string) (value bool, err error) {
    switch str {
    case "1", "t", "T", "true", "TRUE", "True", "YES", "yes", "Yes", "y", "ON", "on", "On":
        return true, nil
    case "0", "f", "F", "false", "FALSE", "False", "NO", "no", "No", "n", "OFF", "off", "Off":
        return false, nil
    }
    return false, fmt.Errorf("parsing \"%s\": invalid syntax", str)
}

// Bool returns bool type value.
func (k *KeyValue) Bool() (bool, error) {
    return parseBool(k.String())
}

// Float64 returns float64 type value.
func (k *KeyValue) Float64() (float64, error) {
    return strconv.ParseFloat(k.String(), 64)
}

// Int returns int type value.
func (k *KeyValue) Int() (int, error) {
    return strconv.Atoi(k.String())
}

// Int64 returns int64 type value.
func (k *KeyValue) Int64() (int64, error) {
    return strconv.ParseInt(k.String(), 10, 64)
}

// Uint returns uint type valued.
func (k *KeyValue) Uint() (uint, error) {
    u, e := strconv.ParseUint(k.String(), 10, 64)
    return uint(u), e
}

// Uint64 returns uint64 type value.
func (k *KeyValue) Uint64() (uint64, error) {
    return strconv.ParseUint(k.String(), 10, 64)
}

// Duration returns time.Duration type value.
func (k *KeyValue) Duration() (time.Duration, error) {
    return time.ParseDuration(k.String())
}

// TimeFormat parses with given format and returns time.Time type value.
func (k *KeyValue) TimeFormat(format string) (time.Time, error) {
    return time.Parse(format, k.String())
}

// Time parses with RFC3339 format and returns time.Time type value.
func (k *KeyValue) Time() (time.Time, error) {
    return k.TimeFormat(time.RFC3339)
}

// MustString returns default value if key value is empty.
func (k *KeyValue) MustString(defaultVal string) string {
    val := k.String()
    if len(val) == 0 {
        k.value = defaultVal
        return defaultVal
    }
    return val
}

// MustBool always returns value without error,
// it returns false if error occurs.
func (k *KeyValue) MustBool(defaultVal ...bool) bool {
    val, err := k.Bool()
    if len(defaultVal) > 0 && err != nil {
        k.value = strconv.FormatBool(defaultVal[0])
        return defaultVal[0]
    }
    return val
}

// MustFloat64 always returns value without error,
// it returns 0.0 if error occurs.
func (k *KeyValue) MustFloat64(defaultVal ...float64) float64 {
    val, err := k.Float64()
    if len(defaultVal) > 0 && err != nil {
        k.value = strconv.FormatFloat(defaultVal[0], 'f', -1, 64)
        return defaultVal[0]
    }
    return val
}

// MustInt always returns value without error,
// it returns 0 if error occurs.
func (k *KeyValue) MustInt(defaultVal ...int) int {
    val, err := k.Int()
    if len(defaultVal) > 0 && err != nil {
        k.value = strconv.FormatInt(int64(defaultVal[0]), 10)
        return defaultVal[0]
    }
    return val
}

// MustInt64 always returns value without error,
// it returns 0 if error occurs.
func (k *KeyValue) MustInt64(defaultVal ...int64) int64 {
    val, err := k.Int64()
    if len(defaultVal) > 0 && err != nil {
        k.value = strconv.FormatInt(defaultVal[0], 10)
        return defaultVal[0]
    }
    return val
}

// MustUint always returns value without error,
// it returns 0 if error occurs.
func (k *KeyValue) MustUint(defaultVal ...uint) uint {
    val, err := k.Uint()
    if len(defaultVal) > 0 && err != nil {
        k.value = strconv.FormatUint(uint64(defaultVal[0]), 10)
        return defaultVal[0]
    }
    return val
}

// MustUint64 always returns value without error,
// it returns 0 if error occurs.
func (k *KeyValue) MustUint64(defaultVal ...uint64) uint64 {
    val, err := k.Uint64()
    if len(defaultVal) > 0 && err != nil {
        k.value = strconv.FormatUint(defaultVal[0], 10)
        return defaultVal[0]
    }
    return val
}

// MustDuration always returns value without error,
// it returns zero value if error occurs.
func (k *KeyValue) MustDuration(defaultVal ...time.Duration) time.Duration {
    val, err := k.Duration()
    if len(defaultVal) > 0 && err != nil {
        k.value = defaultVal[0].String()
        return defaultVal[0]
    }
    return val
}

// MustTimeFormat always parses with given format and returns value without error,
// it returns zero value if error occurs.
func (k *KeyValue) MustTimeFormat(format string, defaultVal ...time.Time) time.Time {
    val, err := k.TimeFormat(format)
    if len(defaultVal) > 0 && err != nil {
        k.value = defaultVal[0].Format(format)
        return defaultVal[0]
    }
    return val
}

// MustTime always parses with RFC3339 format and returns value without error,
// it returns zero value if error occurs.
func (k *KeyValue) MustTime(defaultVal ...time.Time) time.Time {
    return k.MustTimeFormat(time.RFC3339, defaultVal...)
}

// In always returns value without error,
// it returns default value if error occurs or doesn't fit into candidates.
func (k *KeyValue) In(defaultVal string, candidates []string) string {
    val := k.String()
    for _, cand := range candidates {
        if val == cand {
            return val
        }
    }
    return defaultVal
}

// InFloat64 always returns value without error,
// it returns default value if error occurs or doesn't fit into candidates.
func (k *KeyValue) InFloat64(defaultVal float64, candidates []float64) float64 {
    val := k.MustFloat64()
    for _, cand := range candidates {
        if val == cand {
            return val
        }
    }
    return defaultVal
}

// InInt always returns value without error,
// it returns default value if error occurs or doesn't fit into candidates.
func (k *KeyValue) InInt(defaultVal int, candidates []int) int {
    val := k.MustInt()
    for _, cand := range candidates {
        if val == cand {
            return val
        }
    }
    return defaultVal
}

// InInt64 always returns value without error,
// it returns default value if error occurs or doesn't fit into candidates.
func (k *KeyValue) InInt64(defaultVal int64, candidates []int64) int64 {
    val := k.MustInt64()
    for _, cand := range candidates {
        if val == cand {
            return val
        }
    }
    return defaultVal
}

// InUint always returns value without error,
// it returns default value if error occurs or doesn't fit into candidates.
func (k *KeyValue) InUint(defaultVal uint, candidates []uint) uint {
    val := k.MustUint()
    for _, cand := range candidates {
        if val == cand {
            return val
        }
    }
    return defaultVal
}

// InUint64 always returns value without error,
// it returns default value if error occurs or doesn't fit into candidates.
func (k *KeyValue) InUint64(defaultVal uint64, candidates []uint64) uint64 {
    val := k.MustUint64()
    for _, cand := range candidates {
        if val == cand {
            return val
        }
    }
    return defaultVal
}

// InTimeFormat always parses with given format and returns value without error,
// it returns default value if error occurs or doesn't fit into candidates.
func (k *KeyValue) InTimeFormat(format string, defaultVal time.Time, candidates []time.Time) time.Time {
    val := k.MustTimeFormat(format)
    for _, cand := range candidates {
        if val == cand {
            return val
        }
    }
    return defaultVal
}

// InTime always parses with RFC3339 format and returns value without error,
// it returns default value if error occurs or doesn't fit into candidates.
func (k *KeyValue) InTime(defaultVal time.Time, candidates []time.Time) time.Time {
    return k.InTimeFormat(time.RFC3339, defaultVal, candidates)
}

// RangeFloat64 checks if value is in given range inclusively,
// and returns default value if it's not.
func (k *KeyValue) RangeFloat64(defaultVal, min, max float64) float64 {
    val := k.MustFloat64()
    if val < min || val > max {
        return defaultVal
    }
    return val
}

// RangeInt checks if value is in given range inclusively,
// and returns default value if it's not.
func (k *KeyValue) RangeInt(defaultVal, min, max int) int {
    val := k.MustInt()
    if val < min || val > max {
        return defaultVal
    }
    return val
}

// RangeInt64 checks if value is in given range inclusively,
// and returns default value if it's not.
func (k *KeyValue) RangeInt64(defaultVal, min, max int64) int64 {
    val := k.MustInt64()
    if val < min || val > max {
        return defaultVal
    }
    return val
}

// RangeTimeFormat checks if value with given format is in given range inclusively,
// and returns default value if it's not.
func (k *KeyValue) RangeTimeFormat(format string, defaultVal, min, max time.Time) time.Time {
    val := k.MustTimeFormat(format)
    if val.Unix() < min.Unix() || val.Unix() > max.Unix() {
        return defaultVal
    }
    return val
}

// RangeTime checks if value with RFC3339 format is in given range inclusively,
// and returns default value if it's not.
func (k *KeyValue) RangeTime(defaultVal, min, max time.Time) time.Time {
    return k.RangeTimeFormat(time.RFC3339, defaultVal, min, max)
}

// Strings returns list of string divided by given delimiter.
func (k *KeyValue) Strings() []string {
    str := k.String()
    if len(str) == 0 {
        return []string{}
    }

    //vals := strings.Split(str, k.Delims)
    vals := strings.FieldsFunc(str, k.split)
    for i := range vals {
        // vals[i] = k.transformValue(strings.TrimSpace(vals[i]))
        vals[i] = strings.TrimSpace(vals[i])
    }
    return vals
}

func (k *KeyValue) split(r rune) bool {
    for _, delm := range k.Delims {
        if r == delm {
            return true
        }

    }
    return false
}

// Float64s returns list of float64 divided by given delimiter. Any invalid input will be treated as zero value.
func (k *KeyValue) Float64s() []float64 {
    vals, _ := k.parseFloat64s(k.Strings(), true, false)
    return vals
}

// Ints returns list of int divided by given delimiter. Any invalid input will be treated as zero value.
func (k *KeyValue) Ints() []int {
    vals, _ := k.parseInts(k.Strings(), true, false)
    return vals
}

// Int64s returns list of int64 divided by given delimiter. Any invalid input will be treated as zero value.
func (k *KeyValue) Int64s() []int64 {
    vals, _ := k.parseInt64s(k.Strings(), true, false)
    return vals
}

// Int64s returns list of int64 divided by given delimiter. Any invalid input will be treated as zero value.
func (k *KeyValue) Durations() []time.Duration {
    vals, _ := k.parseDurations(k.Strings(), true, false)
    return vals
}

// Uints returns list of uint divided by given delimiter. Any invalid input will be treated as zero value.
func (k *KeyValue) Uints() []uint {
    vals, _ := k.parseUints(k.Strings(), true, false)
    return vals
}

// Uint64s returns list of uint64 divided by given delimiter. Any invalid input will be treated as zero value.
func (k *KeyValue) Uint64s() []uint64 {
    vals, _ := k.parseUint64s(k.Strings(), true, false)
    return vals
}

// TimesFormat parses with given format and returns list of time.Time divided by given delimiter.
// Any invalid input will be treated as zero value (0001-01-01 00:00:00 +0000 UTC).
func (k *KeyValue) TimesFormat(format string) []time.Time {
    vals, _ := k.parseTimesFormat(format, k.Strings(), true, false)
    return vals
}

// Times parses with RFC3339 format and returns list of time.Time divided by given delimiter.
// Any invalid input will be treated as zero value (0001-01-01 00:00:00 +0000 UTC).
func (k *KeyValue) Times() []time.Time {
    return k.TimesFormat(time.RFC3339)
}

// ValidFloat64s returns list of float64 divided by given delimiter. If some value is not float, then
// it will not be included to result list.
func (k *KeyValue) ValidFloat64s() []float64 {
    vals, _ := k.parseFloat64s(k.Strings(), false, false)
    return vals
}

// ValidInts returns list of int divided by given delimiter. If some value is not integer, then it will
// not be included to result list.
func (k *KeyValue) ValidInts() []int {
    vals, _ := k.parseInts(k.Strings(), false, false)
    return vals
}

// ValidInt64s returns list of int64 divided by given delimiter. If some value is not 64-bit integer,
// then it will not be included to result list.
func (k *KeyValue) ValidInt64s() []int64 {
    vals, _ := k.parseInt64s(k.Strings(), false, false)
    return vals
}

// ValidUints returns list of uint divided by given delimiter. If some value is not unsigned integer,
// then it will not be included to result list.
func (k *KeyValue) ValidUints() []uint {
    vals, _ := k.parseUints(k.Strings(), false, false)
    return vals
}

// ValidUint64s returns list of uint64 divided by given delimiter. If some value is not 64-bit unsigned
// integer, then it will not be included to result list.
func (k *KeyValue) ValidUint64s() []uint64 {
    vals, _ := k.parseUint64s(k.Strings(), false, false)
    return vals
}

// ValidTimesFormat parses with given format and returns list of time.Time divided by given delimiter.
func (k *KeyValue) ValidTimesFormat(format string) []time.Time {
    vals, _ := k.parseTimesFormat(format, k.Strings(), false, false)
    return vals
}

// ValidTimes parses with RFC3339 format and returns list of time.Time divided by given delimiter.
func (k *KeyValue) ValidTimes() []time.Time {
    return k.ValidTimesFormat(time.RFC3339)
}

// StrictFloat64s returns list of float64 divided by given delimiter or error on first invalid input.
func (k *KeyValue) StrictFloat64s() ([]float64, error) {
    return k.parseFloat64s(k.Strings(), false, true)
}

// StrictInts returns list of int divided by given delimiter or error on first invalid input.
func (k *KeyValue) StrictInts() ([]int, error) {
    return k.parseInts(k.Strings(), false, true)
}

// StrictInt64s returns list of int64 divided by given delimiter or error on first invalid input.
func (k *KeyValue) StrictInt64s() ([]int64, error) {
    return k.parseInt64s(k.Strings(), false, true)
}

// StrictUints returns list of uint divided by given delimiter or error on first invalid input.
func (k *KeyValue) StrictUints() ([]uint, error) {
    return k.parseUints(k.Strings(), false, true)
}

// StrictUint64s returns list of uint64 divided by given delimiter or error on first invalid input.
func (k *KeyValue) StrictUint64s() ([]uint64, error) {
    return k.parseUint64s(k.Strings(), false, true)
}

// StrictTimesFormat parses with given format and returns list of time.Time divided by given delimiter
// or error on first invalid input.
func (k *KeyValue) StrictTimesFormat(format string) ([]time.Time, error) {
    return k.parseTimesFormat(format, k.Strings(), false, true)
}

// StrictTimes parses with RFC3339 format and returns list of time.Time divided by given delimiter
// or error on first invalid input.
func (k *KeyValue) StrictTimes() ([]time.Time, error) {
    return k.StrictTimesFormat(time.RFC3339)
}

// parseFloat64s transforms strings to float64s.
func (k *KeyValue) parseFloat64s(strs []string, addInvalid, returnOnInvalid bool) ([]float64, error) {
    vals := make([]float64, 0, len(strs))
    for _, str := range strs {
        val, err := strconv.ParseFloat(str, 64)
        if err != nil && returnOnInvalid {
            return nil, err
        }
        if err == nil || addInvalid {
            vals = append(vals, val)
        }
    }
    return vals, nil
}

// parseInts transforms strings to ints.
func (k *KeyValue) parseInts(strs []string, addInvalid, returnOnInvalid bool) ([]int, error) {
    vals := make([]int, 0, len(strs))
    for _, str := range strs {
        val, err := strconv.Atoi(str)
        if err != nil && returnOnInvalid {
            return nil, err
        }
        if err == nil || addInvalid {
            vals = append(vals, val)
        }
    }
    return vals, nil
}

// parseInt64s transforms strings to int64s.
func (k *KeyValue) parseDurations(strs []string, addInvalid, returnOnInvalid bool) ([]time.Duration, error) {
    vals := make([]time.Duration, 0, len(strs))
    for _, str := range strs {
        val, err := time.ParseDuration(str)
        if err != nil && returnOnInvalid {
            return nil, err
        }
        if err == nil || addInvalid {
            vals = append(vals, val)
        }
    }
    return vals, nil
}

// parseInt64s transforms strings to int64s.
func (k *KeyValue) parseInt64s(strs []string, addInvalid, returnOnInvalid bool) ([]int64, error) {
    vals := make([]int64, 0, len(strs))
    for _, str := range strs {
        val, err := strconv.ParseInt(str, 10, 64)
        if err != nil && returnOnInvalid {
            return nil, err
        }
        if err == nil || addInvalid {
            vals = append(vals, val)
        }
    }
    return vals, nil
}

// parseUints transforms strings to uints.
func (k *KeyValue) parseUints(strs []string, addInvalid, returnOnInvalid bool) ([]uint, error) {
    vals := make([]uint, 0, len(strs))
    for _, str := range strs {
        val, err := strconv.ParseUint(str, 10, 0)
        if err != nil && returnOnInvalid {
            return nil, err
        }
        if err == nil || addInvalid {
            vals = append(vals, uint(val))
        }
    }
    return vals, nil
}

// parseUint64s transforms strings to uint64s.
func (k *KeyValue) parseUint64s(strs []string, addInvalid, returnOnInvalid bool) ([]uint64, error) {
    vals := make([]uint64, 0, len(strs))
    for _, str := range strs {
        val, err := strconv.ParseUint(str, 10, 64)
        if err != nil && returnOnInvalid {
            return nil, err
        }
        if err == nil || addInvalid {
            vals = append(vals, val)
        }
    }
    return vals, nil
}

// parseTimesFormat transforms strings to times in given format.
func (k *KeyValue) parseTimesFormat(format string, strs []string, addInvalid, returnOnInvalid bool) ([]time.Time, error) {
    vals := make([]time.Time, 0, len(strs))
    for _, str := range strs {
        val, err := time.Parse(format, str)
        if err != nil && returnOnInvalid {
            return nil, err
        }
        if err == nil || addInvalid {
            vals = append(vals, val)
        }
    }
    return vals, nil
}
