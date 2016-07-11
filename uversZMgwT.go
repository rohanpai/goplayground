package main
 
import (
	&#34;fmt&#34;
	&#34;reflect&#34;
	&#34;sync&#34;
)
 
type EventBus struct {
	handlers map[reflect.Type][]reflect.Value
	lock     sync.RWMutex
}
 
func New() *EventBus {
	return &amp;EventBus{
		make(map[reflect.Type][]reflect.Value),
		sync.RWMutex{},
	}
}
 
func (bus *EventBus) addHandler(t reflect.Type, fn reflect.Value) {
	bus.lock.Lock()
	bus.lock.Unlock()
	handlers, ok := bus.handlers[t]
	if !ok {
		handlers = make([]reflect.Value, 0)
	}
	bus.handlers[t] = append(handlers, fn)
}
 
func (bus *EventBus) RegisterHandler(fn interface{}, forTypes ...interface{}) {
	v := reflect.ValueOf(fn)
	def := v.Type()
 
	// the message handler must have a single parameter
	if def.NumIn() != 1 {
		panic(&#34;Handler must have a single argument&#34;)
	}
	// find out the handler argument type
	argument := def.In(0)
 
	// check wether we can convert the types into the argument
	for _, typ := range forTypes {
		t := reflect.TypeOf(typ)
		if !t.ConvertibleTo(argument) {
			panic(fmt.Sprintf(&#34;Handler argument %v is not compatible with type %v&#34;, argument, t))
		}
		bus.addHandler(t, v)
	}
 
	// if we aren&#39;t specific, we just handle the specified message
	if len(forTypes) == 0 {
		bus.addHandler(argument, v)
	}
}
 
func (bus *EventBus) Publish(ev interface{}) error {
	bus.lock.RLock()
	defer bus.lock.RUnlock()
 
	t := reflect.TypeOf(ev)
 
	handlers, ok := bus.handlers[t]
	if !ok {
		return nil
	}
 
	args := [...]reflect.Value{reflect.ValueOf(ev)}
	for _, fn := range handlers {
		fn.Call(args[:])
	}
	return nil
}

/// example


type SomeEvent struct {
  A, B string
}
 
func SomeEventHandler(ev SomeEvent) {
  fmt.Printf(&#34;%s =&gt; %s\n&#34;, ev.A, ev.B)
}
 
func main(){
  bus := New()
  bus.RegisterHandler(SomeEventHandler)
  bus.Publish(SomeEvent{&#34;A&#34;, &#34;B&#34;})
}

