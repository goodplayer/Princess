package app

// Lifecycle support a bunch of APIs required for lifecycle
// Certain initialization order will be automatically calculated according to the dependency injection.
// Don't use circle dependency, otherwise it will cause unexpected behavior.
type Lifecycle interface {
	PrepareInitialization() error
	Initialize() error
	PrepareDestruction() error
	Destroy() error
}

type DefaultLifecycle struct {
}

func (d DefaultLifecycle) PrepareInitialization() error {
	return nil
}

func (d DefaultLifecycle) Initialize() error {
	return nil
}

func (d DefaultLifecycle) PrepareDestruction() error {
	return nil
}

func (d DefaultLifecycle) Destroy() error {
	return nil
}
