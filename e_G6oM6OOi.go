package main

import (
	&#34;container/list&#34;
	&#34;fmt&#34;
)
/////////////语境/////////////////////

type 比较器 interface {
	比一比(interface{}, interface{}) int
}

/////////////代码/////////////////////

func 归并排序(待排序列表 *list.List, 某个比较器 比较器) {

	长度 := 待排序列表.Len()
	if 长度 &lt; 2 {
		return
	}
	
	左边儿那半个, 右边儿那半个 := 拆开(待排序列表)
	
	归并排序(左边儿那半个, 某个比较器)
	归并排序(右边儿那半个, 某个比较器)
	归并(左边儿那半个, 右边儿那半个, 待排序列表, 某个比较器)
}

func 归并(表1 *list.List, 表2 *list.List, 原表 *list.List, 某个比较器 比较器) {
	for 表1.Len() != 0 &amp;&amp; 表2.Len() != 0 {
		if 某个比较器.比一比(表1.Front().Value, 表2.Front().Value) &lt;= 0 {
			原表.PushBack(表1.Remove(表1.Front()))
		} else {
			原表.PushBack(表2.Remove(表2.Front()))
		}
	}
	for 表1.Len() != 0 {
		原表.PushBack(表1.Remove(表1.Front()))
	}
	for 表2.Len() != 0 {
		原表.PushBack(表2.Remove(表2.Front()))
	}
}

func 拆开(待拆表 *list.List) (表1 *list.List, 表2 *list.List) {
	表1, 表2 = list.New(), list.New()
	长度的一半 := 待拆表.Len()/2
	
	for i := 0; i &lt; 长度的一半; i&#43;&#43; {
		表1.PushBack(待拆表.Remove(待拆表.Front()))
	}
	for 待拆表.Len() != 0 {
		表2.PushBack(待拆表.Remove(待拆表.Front()))
	}
	return
}

////////////比较器/////////

type 字典序比较器 struct{}

func (p 字典序比较器) 比一比(i1, i2 interface{}) int {
	if i1.(string) &lt; i2.(string) {
		return -1
	}
	if i1.(string) &gt; i2.(string) {
		return 1
	}
	return 0
}

type 字符长度比较器 struct{}

func (p 字符长度比较器) 比一比(i1, i2 interface{}) int {
	if len(i1.(string)) &lt; len(i2.(string)) {
		return -1
	}
	if len(i1.(string)) &gt; len(i2.(string)) {
		return 1
	}
	return 0
}
////////////////////////
func main() {
	表 := list.New()
	表.PushBack(&#34;hello&#34;)
	表.PushBack(&#34;namaste&#34;)
	表.PushBack(&#34;aloha&#34;)
	表.PushBack(&#34;你好&#34;)
	//var i 字符长度比较器
	var i 字典序比较器
	归并排序(表, i)
	遍历(表)
}

func 遍历(列表 *list.List) {
	f := &#34;%#v &#34;
	fmt.Print(&#34;[ &#34;)
	for i := 列表.Front(); i != nil; i = i.Next() {
		fmt.Printf(f, i.Value)
	}
	fmt.Println(&#34;]&#34;)
}
