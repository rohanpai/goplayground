package main

import &#34;fmt&#34;

type Monad interface {
	Bind(func(interface{}, Monad) Monad) Monad
	Return(interface{}) Monad
}

// First, let&#39;s do Maybe
type Maybe struct {
	val *interface{}
}

var None = Maybe{nil}

func Some(a interface{}) Monad {
	return Maybe{&amp;a}
}

func (m Maybe) Bind(f func(interface{}, Monad) Monad) Monad {
	if m == None {
		return None
	}
	return f(*m.val, m)
}

func (m Maybe) Return(a interface{}) Monad {
	return Some(a)
}

func (m Maybe) String() string {
	if m == None {
		return &#34;None&#34;
	}
	return fmt.Sprintf(&#34;Some(%v)&#34;, *m.val)
}

// Let&#39;s do a slice of int!
type IntSliceMonad []int

func (m IntSliceMonad) Bind(f func(interface{}, Monad) Monad) Monad {
	result := IntSliceMonad{}
	for _, val := range m {
		next := f(val, m).(IntSliceMonad)
		result = append(result, next...)
	}
	return result
}

func (m IntSliceMonad) Return(a interface{}) Monad {
	return IntSliceMonad{a.(int)}
}

func main() {
	doubleMe := func(i interface{}, m Monad) Monad {
		return m.Return(2 * i.(int))
	}

	tripleMe := func(i interface{}, m Monad) Monad {
		return m.Return(3 * i.(int))
	}

	oops := func(i interface{}, m Monad) Monad {
		return None
	}

	bothSigns := func(i interface{}, m Monad) Monad {
		return IntSliceMonad{i.(int), -i.(int)}
	}

	fmt.Println(Some(5).Bind(doubleMe).Bind(tripleMe))
	fmt.Println(Some(3).Bind(oops).Bind(doubleMe))

	var oneTwoThree = IntSliceMonad{1, 2, 3}
	fmt.Println(oneTwoThree.Bind(doubleMe).Bind(bothSigns))
}
