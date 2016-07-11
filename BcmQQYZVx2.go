package main

//author: www@kaibb.com
//2013-11-19 17:43:19
import (
	&#34;fmt&#34;
)

func main() {
	var a string
	a = &#34;类型判断,只支持interface{}&#34;
	// 如果写成 b ,ok:=a.(string) ;
	// 将会报 invalid type assertion: a.(string) (non-interface type string on left)
	// 因为类型判断只支持interface{}类型
	// 这儿强制转换 interface{}(a)
	if b, ok := interface{}(a).(string); ok {
		fmt.Println(b)
	} else {
		fmt.Println(&#34;Not String&#34;)
	}
}
