package trigger

import (
	"errors"
	"reflect"
	"sync"
	"fmt"
)

func New() Trigger {
	return &trigger{
		functionMap: make(map[string]*fn),
	}
}

type trigger struct {
	functionMap map[string]*fn

	mu sync.Mutex
}

type fn struct {
	Fn  []interface{}
	Typ string
}

func (t *trigger) On(event string, task interface{}) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	typ := reflect.ValueOf(task).Type()
	if typ.Kind() != reflect.Func {
		return errors.New("task is not a function")
	}
	if listeners, ok := t.functionMap[event]; ok {
		if listeners.Typ != typ.String() {
			panic(fmt.Sprint("could not register \"", event,"\" event listener as ", typ.String(), " previously registered as ", listeners.Typ, ))
		}
		listeners.Fn = append(listeners.Fn, task)
	} else {
		t.functionMap[event] = &fn{
			Fn:  []interface{}{task},
			Typ: typ.String(),
		}
	}
	return nil
}

func (t *trigger) Fire(event string, params ...interface{}) (error) {
	listeners, err := t.read(event, params...)
	if err != nil {
		return err
	}
	for _, f := range listeners.Functions {
		//result := f.Call(listeners.Params)
		f.Call(listeners.Params)
	}
	return nil
}

func (t *trigger) FireBackground(event string, params ...interface{}) (error) {
	listeners, err := t.read(event, params...)
	if err != nil {
		return err
	}
	for _, f := range listeners.Functions {
		//result := f.Call(listeners.Params)
		go f.Call(listeners.Params)
	}
	return nil
}

func (t *trigger) Clear(event string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.functionMap[event]; !ok {
		return errors.New("event not defined")
	}
	delete(t.functionMap, event)
	return nil
}

func (t *trigger) ClearEvents() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.functionMap = make(map[string]*fn)
}

func (t *trigger) HasEvent(event string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	_, ok := t.functionMap[event]
	return ok
}

func (t *trigger) Events() map[string]string {
	t.mu.Lock()
	defer t.mu.Unlock()
	events := make(map[string]string)
	for k, v := range t.functionMap {
		events[k] = v.Typ
	}
	return events
}

func (t *trigger) EventCount() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.functionMap)
}

type Listeners struct {
	Functions []reflect.Value
	Params    []reflect.Value
}

func (t *trigger) read(event string, params ...interface{}) (Listeners, error) {
	t.mu.Lock()
	tasks, ok := t.functionMap[event]
	t.mu.Unlock()
	var listeners = Listeners{}
	if !ok {
		return listeners, errors.New("no task found for event")
	}

	// functions
	for _, task := range tasks.Fn {
		f := reflect.ValueOf(task)
		if len(params) != f.Type().NumIn() {
			panic(fmt.Sprint("parameters count mismatched for event \"", event,"\" required ", f.Type().NumIn(), " got ", len(params), ))
		}
		listeners.Functions = append(listeners.Functions, f)
	}

	// params
	listeners.Params = make([]reflect.Value, len(params))
	for k, param := range params {
		listeners.Params[k] = reflect.ValueOf(param)
	}
	return listeners, nil
}
