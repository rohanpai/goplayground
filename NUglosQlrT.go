package main

import (
    "errors"
    "fmt"
    "reflect"
    "sync"
)

type invoke struct {
    fn   reflect.Value
    args []reflect.Value
}

func (self *invoke) call() []reflect.Value {
    return self.fn.Call(self.args)
}

type invokeList struct {
    list  []*invoke
    mutex sync.Mutex
}

func (self *invokeList) add(fn interface{}, args ...interface{}) error {
    self.mutex.Lock()
    defer self.mutex.Unlock()

    fnValue := reflect.ValueOf(fn)
    if fnValue.Kind() != reflect.Func {
        return errors.New("Non-callable function")
    }
    if fnValue.Type().NumIn() != len(args) {
        return errors.New("Error arguments number")
    }

    argsValue := make([]reflect.Value, len(args))
    for i, arg := range args {
        argsValue[i] = reflect.ValueOf(arg)
    }

    self.list = append(self.list, &invoke{fnValue, argsValue})
    return nil
}

func (self *invokeList) call() (result [][]reflect.Value) {
    self.mutex.Lock()
    defer self.mutex.Unlock()

    result = make([][]reflect.Value, len(self.list))
    for i, f := range self.list {
        result[i] = f.call()
    }
    return
}

func main() {
	// Put invokeList into some non-main package, new destructors in init(),
	// use destructors.add() anywhere you need,
	// and defer destructors.call() in main.
	
	destructors := new(invokeList)
	defer destructors.call()
	
	done := func(msg string) {
		fmt.Println(msg)
	}
	
	destructors.add(done, "I'm done")
}