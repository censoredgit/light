package controller

import (
	"errors"
	"fmt"
	"reflect"
)

type Container struct {
	items []any
}

func (c *Container) MustSingleton(instance any) {
	err := c.Singleton(instance)
	if err != nil {
		panic(err)
	}
}

func (c *Container) Singleton(instance any) error {
	if reflect.TypeOf(instance).Kind() != reflect.Pointer {
		return errors.New("instance object must be a pointer")
	}

	if reflect.TypeOf(instance).Elem().Kind() != reflect.Struct {
		return errors.New("instance object must be a pointer to a structure")
	}

	c.items = append(c.items, instance)

	return nil
}

func (c *Container) inject(f any) any {
	inputFunction := reflect.ValueOf(f)

	if inputFunction.Kind() != reflect.Func {
		if reflect.ValueOf(&Action{}).Type().AssignableTo(inputFunction.Type()) {
			return inputFunction.Interface().(*Action)
		}

		panic("parameter must be a function that return *Action")
	}

	if inputFunction.Type().NumOut() != 1 {
		panic("parameter must have one return value")
	}

	if inputFunction.Type().Out(0).Kind() != reflect.Pointer {
		panic("parameter must have a return value as a pointer")
	}

	if !inputFunction.Type().Out(0).Elem().AssignableTo(reflect.ValueOf(&Action{}).Elem().Type()) {
		panic("parameter must have *Action return type")
	}

	newArgs := make([]reflect.Value, 0)

	var isResolvedParam bool
	inParamsCount := inputFunction.Type().NumIn()
	for index := 0; index < inParamsCount; index++ {
		isResolvedParam = false

	outLoop:
		for _, v := range c.items {
			var funcParameter reflect.Type
			containerRecord := reflect.ValueOf(v).Elem().Type()

			if inputFunction.Type().In(index).Kind() == reflect.Pointer {
				funcParameter = inputFunction.Type().In(index).Elem()
				if containerRecord.AssignableTo(funcParameter) {
					newArgs = append(newArgs, reflect.ValueOf(v))
					isResolvedParam = true
					break outLoop
				}
			} else if inputFunction.Type().In(index).Kind() == reflect.Interface {
				funcParameter = inputFunction.Type().In(index)
				if reflect.PointerTo(containerRecord).Implements(funcParameter) {
					newArgs = append(newArgs, reflect.ValueOf(v))
					isResolvedParam = true
					break outLoop
				}
			} else {
				panic(fmt.Sprintf("type (%s) must be of type pointer or interface, (%s) given",
					inputFunction.Type().In(index).Name(),
					inputFunction.Type().In(index).Kind()))
			}
		}

		if !isResolvedParam {
			switch inputFunction.Type().In(index).Kind() {
			case reflect.Pointer:
				panic(fmt.Sprintf("can not resolve param type *%s", inputFunction.Type().In(index).Elem().Name()))
			default:
				panic(fmt.Sprintf("can not resolve param type %s", inputFunction.Type().In(index).Name()))
			}
		}
	}

	if inputFunction.Type().NumIn() != len(newArgs) {
		panic(fmt.Sprint("can not resolve all arguments"))
	}

	return inputFunction.Call(newArgs)[0].Interface().(*Action)
}
