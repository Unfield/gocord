package gateway

import (
	"encoding/json"
	"log"
	"reflect"
	"sync"
)

type TypedDispatcher struct {
	mu         sync.RWMutex
	listeners  map[EventName][]reflect.Value
	eventTypes map[EventName]reflect.Type
}

func NewTypedDispatcher() *TypedDispatcher {
	return &TypedDispatcher{
		listeners:  make(map[EventName][]reflect.Value),
		eventTypes: make(map[EventName]reflect.Type),
	}
}

func RegisterType[T any](d *TypedDispatcher, name EventName) {
	var t T
	d.eventTypes[name] = reflect.TypeOf(t)
}

func On[T any](d *TypedDispatcher, name EventName, fn func(T)) {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		panic("On requires a function")
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	d.listeners[name] = append(d.listeners[name], reflect.ValueOf(fn))
	d.eventTypes[name] = reflect.TypeOf(*new(T))
}

func (d *TypedDispatcher) dispatch(name EventName, raw any) {
	d.mu.RLock()
	handlers := d.listeners[name]
	t, ok := d.eventTypes[name]
	d.mu.RUnlock()

	if !ok || len(handlers) == 0 {
		return
	}

	b, _ := json.Marshal(raw)
	ptr := reflect.New(t)
	if err := json.Unmarshal(b, ptr.Interface()); err != nil {
		log.Printf("[Dispatcher] decode error for %s: %v", name, err)
		return
	}
	val := ptr.Elem().Interface()

	for _, h := range handlers {
		h.Call([]reflect.Value{reflect.ValueOf(val)})
	}
}
