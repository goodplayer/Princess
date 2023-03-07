package app

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
)

const (
	_DefaultMethodNameForDependencyInjection = "DependOn"
)

// ApplicationContainer support simple object container
// If object provide 'DependOn' method with several parameters as dependencies,
// then DI will be effective and dependencies will be provided for this method.
type ApplicationContainer struct {
	sync                          sync.RWMutex
	once                          sync.Once
	store                         map[reflect.Type]interface{}
	pendingDependencyInjectObject map[reflect.Type]dependencyInjectionEntry

	initialized int32 // 0 - uninitialized, 1 - initialized, other - other states

	//TODO support lifecycle
}

type dependencyInjectionEntry struct {
	argTypes []reflect.Type
}

func NewAppContainer() *ApplicationContainer {
	return &ApplicationContainer{
		store:                         map[reflect.Type]interface{}{},
		pendingDependencyInjectObject: map[reflect.Type]dependencyInjectionEntry{},
	}
}

func (a *ApplicationContainer) Get(t reflect.Type) interface{} {
	// should block operations on uninitialized container
	a.ensureInitialized()
	a.sync.RLock()
	defer a.sync.RUnlock()

	v, ok := a.store[t]
	if !ok {
		panic(errors.New("object not found"))
	}
	return v
}

// Set will register an object to the container
// In order to use dependency injection, the correct type(struct or *struct) of object should be used,
// for both 'DependOn' receiver and provided objects.
func (a *ApplicationContainer) Set(obj interface{}) {
	// should block operations on uninitialized container
	a.ensureUninitialized()
	a.sync.Lock()
	defer a.sync.Unlock()

	if obj == nil {
		panic(errors.New("object is nil"))
	}
	// object duplicated report
	if _, ok := a.store[reflect.TypeOf(obj)]; ok {
		panic(errors.New("duplicated type registered:" + fmt.Sprint(reflect.TypeOf(obj))))
	}
	a.store[reflect.TypeOf(obj)] = obj
	// prepare DI
	a.simpleDependencyInjection(obj)
}

func (a *ApplicationContainer) simpleDependencyInjection(obj interface{}) {
	t := reflect.TypeOf(obj)
	method, ok := t.MethodByName(_DefaultMethodNameForDependencyInjection)
	if ok {
		if method.Type.NumIn() == 1 {
			// by default, struct will be the first parameter
			// so there won't be required dependency
			return
		}
		var argTypes []reflect.Type
		for i := 1; i < method.Type.NumIn(); i++ {
			argTypes = append(argTypes, method.Type.In(i))
		}
		a.pendingDependencyInjectObject[t] = dependencyInjectionEntry{
			argTypes: argTypes,
		}
	}
}

func (a *ApplicationContainer) StartContainer() {
	a.once.Do(func() {
		// do dependency injection
		{
			// 1. check
			for t, entry := range a.pendingDependencyInjectObject {
				objects := make([]reflect.Value, 0, len(entry.argTypes))
				for _, argType := range entry.argTypes {
					obj, ok := a.store[argType]
					if !ok {
						panic(errors.New("cannot find type:" + fmt.Sprint(argType) + " for object:" + fmt.Sprint(t)))
					}
					objects = append(objects, reflect.ValueOf(obj))
				}
				if obj, ok := a.store[t]; !ok {
					panic(errors.New("cannot find object:" + fmt.Sprint(t)))
				} else {
					// 2. invoke if dependencies are prepared
					reflect.ValueOf(obj).MethodByName(_DefaultMethodNameForDependencyInjection).Call(objects)
				}
			}
		}

		//TODO lifecycle management
		//_, ok := reflect.ValueOf(obj).Interface().(Lifecycle)
	})
	// finally mark as initialized
	atomic.CompareAndSwapInt32(&a.initialized, 0, 1)
}

func (a *ApplicationContainer) isInitialized() bool {
	return atomic.LoadInt32(&a.initialized) == 1
}

func (a *ApplicationContainer) ensureInitialized() {
	if !a.isInitialized() {
		panic(errors.New("container is not in initialized state"))
	}
}

func (a *ApplicationContainer) isUninitialized() bool {
	return atomic.LoadInt32(&a.initialized) == 0
}

func (a *ApplicationContainer) ensureUninitialized() {
	if !a.isUninitialized() {
		panic(errors.New("container is not in uninitialized state"))
	}
}

func (a *ApplicationContainer) WaitingShutdown() error {
	//TODO implement blocking on container completely shutting down
	return nil
}
