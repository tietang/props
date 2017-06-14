package props

import (
	"strconv"
	"errors"
)

type MapProperties struct {
	values map[string]string
}

func NewMapProperties() *MapProperties {
	p := &MapProperties{
		values: make(map[string]string),
	}
	return p
}

func NewMapPropertiesByMap(kv map[string]string) *MapProperties {
	p := &MapProperties{
		values: kv,
	}
	return p
}

//--get key/value

// Get retrieves the value of a property. If the property does not exist, an
// empty string will be returned.
func (p *MapProperties) Get(key string) (string, error) {
	if v, ok := p.values[key]; ok {
		return v, nil
	}
	return "", errors.New("not exists for key: " + key)
}

// GetDefault retrieves the value of a property. If the property does not
// exist, then the default value will be returned.
func (p *MapProperties) GetDefault(key, defVal string) string {
	if v, ok := p.values[key]; ok {
		return v
	}
	return defVal
}

func (p *MapProperties) GetInt(key string) (int, error) {
	v := p.values[key]

	if v, err := strconv.Atoi(v); err == nil {
		return v, nil
	}
	return 0, errors.New("not exists for key: " + key)
}

func (p *MapProperties) GetIntDefault(key string, defVal int) int {
	if v, ok := p.values[key]; ok {
		if v, err := strconv.Atoi(v); err == nil {
			return v
		}
	}
	return defVal
}
func (p *MapProperties) GetBool(key string) (bool, error) {
	v := p.values[key]

	if v, err := strconv.ParseBool(v); err == nil {
		return v, nil
	}
	return false, errors.New("not exists for key: " + key)
}

func (p *MapProperties) GetBoolDefault(key string, defVal bool) bool {
	if v, ok := p.values[key]; ok {
		if v, err := strconv.ParseBool(v); err == nil {
			return v
		}
	}
	return defVal
}

func (p *MapProperties) GetFloat64(key string) (float64, error) {
	v := p.values[key]

	if v, err := strconv.ParseFloat(v, 64); err == nil {
		return v, nil
	}
	return 0.0, errors.New("not exists for key: " + key)
}

func (p *MapProperties) GetFloat64Default(key string, defVal float64) float64 {
	if v, ok := p.values[key]; ok {
		if v, err := strconv.ParseFloat(v, 64); err == nil {
			return v
		}
	}
	return defVal
}

// Names returns the keys for all properties in the set.
func (p *MapProperties) Keys() []string {
	keys := make([]string, 0, len(p.values))
	for k, _ := range p.values {
		keys = append(keys, k)
	}
	return keys
}

// Set adds or changes the value of a property.
func (p *MapProperties) Set(key, val string) {
	p.values[key] = val
}
func (p *MapProperties) SetAll(values map[string]string) {
	for k, v := range values {
		p.values[k] = v
	}

}

// Clear removes all key-value pairs.
func (p *MapProperties) Clear() {
	p.values = make(map[string]string)
}
