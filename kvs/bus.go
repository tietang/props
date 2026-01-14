package kvs

import "github.com/asaskevich/EventBus"

const (
	ON_CHANGE = "on:change"
)

var bus = EventBus.New()

func init() {
	bus.SubscribeAsync(ON_CHANGE, onChange, false)
}

func onChange(listeners []func(k, v string), k, v string) {
	for _, listener := range listeners {
		listener(k, v)
	}
}

func CheckAndPublishChange(values map[string]string, key, val string, onChanges map[string][]func(k, v string)) {
	if onChanges == nil {
		return
	}
	listeners, found := onChanges[key]
	if !found {
		return
	}
	v, found := values[key]
	if found && v != val {
		bus.Publish(ON_CHANGE, listeners, key, val)
	}
}

func PublishChange(isChanged bool, key, val string, onChanges map[string][]func(k, v string)) {
	if onChanges == nil {
		return
	}
	listeners, found := onChanges[key]
	if !found {
		return
	}
	if isChanged {
		bus.Publish(ON_CHANGE, listeners, key, val)
	}
}
