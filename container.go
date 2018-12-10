// Package di114514 provides an dependency injection container.
package di114514

import (
	"errors"
	"reflect"
)

type FactoryWithContainerFunction func(ContainerInterface) (interface{})
type FactoryFunction func() (interface{})

type ContainerInterface interface {
	Define(name string, value interface{}) error
	GetInstance(name string) (interface{})
	NewInstance(name string) (interface{})
}

type container struct {
	instances                             map[string]interface{}
	factoryWithContainerFunctions         map[string]FactoryWithContainerFunction
	factoryFunctions                      map[string]FactoryFunction
}

func NewContainer() ContainerInterface {
	return &container{
		instances:                             make(map[string]interface{}),
		factoryWithContainerFunctions:         make(map[string]FactoryWithContainerFunction),
		factoryFunctions:                      make(map[string]FactoryFunction),
	}
}

func (this *container) Define(name string, value interface{}) error {
	functionReflValue := reflect.ValueOf(value);
	functionReflType := functionReflValue.Type()
	if functionReflType.Kind() != reflect.Func {
		this.instances[name] = value
		return nil
	}

	numIn := functionReflType.NumIn()
	numOut := functionReflType.NumOut()

	if numOut != 1 {
		return errors.New("factory function's return value must be single argument.")
	}

	if functionReflType.Out(0).Kind() != reflect.Interface {
		return errors.New("factory function's return type must be interface.")
	}

	if numIn == 0 {
		this.factoryFunctions[name] = func() interface{} {
			return functionReflValue.Call([]reflect.Value{})[0].Interface()
		}
		return nil
	}

	if numIn == 1 {
		targetType := reflect.TypeOf((*ContainerInterface)(nil)).Elem()
		if functionReflType.In(0).Implements(targetType) {
			this.factoryWithContainerFunctions[name] = func(c ContainerInterface) interface{} {
				return functionReflValue.Call([]reflect.Value{reflect.ValueOf(c)})[0].Interface()
			}
			return nil
		}
		return errors.New("factory function requires ContainerInterface type parameter.")
	}

	return errors.New("factory function's return type must be interface.")
}

// GetInstance gets an instance from container
func (this *container) GetInstance(name string) interface{} {
	if instance, ok := this.instances[name]; ok {
		return instance;
	}

	if factoryWithContainerFunction, ok := this.factoryWithContainerFunctions[name]; ok {
		instance := factoryWithContainerFunction(this)

		this.instances[name] = instance;
		return instance;
	}

	if factoryFunction, ok := this.factoryFunctions[name]; ok {
		instance := factoryFunction()

		this.instances[name] = instance;
		return instance;
	}

	return nil
}

// NewInstance creates a new instance
func (this *container) NewInstance(name string) interface{} {
	if factoryWithContainerFunction, ok := this.factoryWithContainerFunctions[name]; ok {
		return factoryWithContainerFunction(this)
	}

	if factoryFunction, ok := this.factoryFunctions[name]; ok {
		return factoryFunction()
	}

	return nil
}
