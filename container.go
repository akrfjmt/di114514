// Package di114514 provides an dependency injection container.
package di114514

import (
	"errors"
	"reflect"
)

type factoryWithContainerFunction func(ContainerInterface) interface{}
type factoryFunction func() interface{}

// ContainerInterface はコンテナのインターフェース
type ContainerInterface interface {
	Define(name string, value interface{}) error
	GetInstance(name string) interface{}
	NewInstance(name string) interface{}
}

type container struct {
	instances                     map[string]interface{}
	factoryWithContainerFunctions map[string]factoryWithContainerFunction
	factoryFunctions              map[string]factoryFunction
}

// NewContainer はコンテナのコンストラクタ
func NewContainer() ContainerInterface {
	return &container{
		instances:                     make(map[string]interface{}),
		factoryWithContainerFunctions: make(map[string]factoryWithContainerFunction),
		factoryFunctions:              make(map[string]factoryFunction),
	}
}

func (c *container) Define(name string, value interface{}) error {
	functionReflValue := reflect.ValueOf(value)
	functionReflType := functionReflValue.Type()
	if functionReflType.Kind() != reflect.Func {
		c.instances[name] = value
		return nil
	}

	numIn := functionReflType.NumIn()
	numOut := functionReflType.NumOut()

	if numOut != 1 {
		return errors.New("factory function's return value must be single argument")
	}

	if functionReflType.Out(0).Kind() != reflect.Interface {
		return errors.New("factory function's return type must be interface")
	}

	if numIn == 0 {
		c.factoryFunctions[name] = func() interface{} {
			return functionReflValue.Call([]reflect.Value{})[0].Interface()
		}
		return nil
	}

	if numIn == 1 {
		targetType := reflect.TypeOf((*ContainerInterface)(nil)).Elem()
		if functionReflType.In(0).Implements(targetType) {
			c.factoryWithContainerFunctions[name] = func(c ContainerInterface) interface{} {
				return functionReflValue.Call([]reflect.Value{reflect.ValueOf(c)})[0].Interface()
			}
			return nil
		}
		return errors.New("factory function requires ContainerInterface type parameter")
	}

	return errors.New("factory function's return type must be interface")
}

// GetInstance gets an instance from container
func (c *container) GetInstance(name string) interface{} {
	if instance, ok := c.instances[name]; ok {
		return instance
	}

	if factoryWithContainerFunction, ok := c.factoryWithContainerFunctions[name]; ok {
		instance := factoryWithContainerFunction(c)

		c.instances[name] = instance
		return instance
	}

	if factoryFunction, ok := c.factoryFunctions[name]; ok {
		instance := factoryFunction()

		c.instances[name] = instance
		return instance
	}

	return nil
}

// NewInstance creates a new instance
func (c *container) NewInstance(name string) interface{} {
	if factoryWithContainerFunction, ok := c.factoryWithContainerFunctions[name]; ok {
		return factoryWithContainerFunction(c)
	}

	if factoryFunction, ok := c.factoryFunctions[name]; ok {
		return factoryFunction()
	}

	return nil
}
