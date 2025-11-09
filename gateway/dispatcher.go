package gateway

import (
	"encoding/json"
	"log"
	"reflect"
	"sync"
)

type handlerEntry struct {
	fn     reflect.Value
	argTyp reflect.Type
}

type Dispatcher struct {
	mu        sync.Mutex
	listeners map[string][]handlerEntry
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		listeners: make(map[string][]handlerEntry),
	}
}

func (d *Dispatcher) On(event string, handler any) {
	val := reflect.ValueOf(handler)
	typ := val.Type()

	if typ.Kind() != reflect.Func || typ.NumIn() != 1 {
		panic("handler must be a function with exactly one argument")
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	entry := handlerEntry{
		fn:     val,
		argTyp: typ.In(0),
	}
	d.listeners[event] = append(d.listeners[event], entry)
}

func (d *Dispatcher) dispatch(event string, data any) {
	d.mu.Lock()
	list := d.listeners[event]
	d.mu.Unlock()

	for _, entry := range list {
		argPtr := reflect.New(entry.argTyp)
		raw, _ := json.Marshal(data)
		if err := json.Unmarshal(raw, argPtr.Interface()); err != nil {
			log.Printf("[Dispatcher] decode error for %s: %v", event, err)
			continue
		}

		entry.fn.Call([]reflect.Value{argPtr.Elem()})
	}
}
