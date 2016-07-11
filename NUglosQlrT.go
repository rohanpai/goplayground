package main

import (
    &#34;errors&#34;
    &#34;fmt&#34;
    &#34;reflect&#34;
    &#34;sync&#34;
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
        return errors.New(&#34;Non-callable function&#34;)
    }
    if fnValue.Type().NumIn() != len(args) {
        return errors.New(&#34;Error arguments number&#34;)
    }

    argsValue := make([]reflect.Value, len(args))
    for i, arg := range args {
        argsValue[i] = reflect.ValueOf(arg)
    }

    self.list = append(self.list, &amp;invoke{fnValue, argsValue})
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
	
	destructors.add(done, &#34;I&#39;m done&#34;)
}