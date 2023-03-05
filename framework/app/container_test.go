package app

import (
	"reflect"
	"testing"
)

type Stub struct {
	demoService *DemoService
}

type Dependency1 struct {
	val string
}

type DemoService struct {
	DefaultLifecycle
	val string
	dep *Dependency1
}

func (d *DemoService) RequestFor(dep *Dependency1) {
	d.dep = dep
}

func TestSetGetObject(t *testing.T) {
	app := NewAppContainer()
	app.Set(&Dependency1{val: "dependency1"})

	stub := new(Stub)
	demo := new(DemoService)
	demo.val = "hello world"

	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("not found expected")
			}
		}()
		app.Get(reflect.TypeOf(&DemoService{}))
	}()
	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("not found expected")
			}
		}()
		app.Get(reflect.TypeOf(stub.demoService))
	}()

	app.Set(demo)

	app.StartContainer()

	if v := app.Get(reflect.TypeOf(&DemoService{})); v == nil {
		t.Fatal("val expected")
	} else if v.(*DemoService).val != demo.val {
		t.Fatal("val mismatch")
	}
	if v := app.Get(reflect.TypeOf(stub.demoService)); v == nil {
		t.Fatal("val expected")
	} else if v.(*DemoService).val != demo.val {
		t.Fatal("val mismatch")
	}

	if v := app.Get(reflect.TypeOf(stub.demoService)); v == nil {
		t.Fatal("val expected")
	} else if v.(*DemoService).dep.val != "dependency1" {
		t.Fatal("val mismatch")
	}
}
