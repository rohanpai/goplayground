package main

import (
	&#34;fmt&#34;
)

func main() {
	var sortArray = []int{1, 3, 41, 24, 76, 11, 45, 3, 3, 64, 21, 69, 19, 36}
	fmt.Println(sortArray)
	qsort(sortArray, 0, len(sortArray)-1)
	fmt.Println(sortArray)
}

func qsort(array []int, low, high int) {
	if low &lt; high {
		m := partition(array, low, high)
		// fmt.Println(m)
		qsort(array, low, m-1)
		qsort(array, m&#43;1, high)
	}
}

func partition(array []int, low, high int) int {
	key := array[low]
	tmpLow := low
	tmpHigh := high
	for {
		//查找小于等于key的元素，该元素的位置一定是tmpLow到high之间，因为array[tmpLow]及左边元素小于等于key，不会越界
		for array[tmpHigh] &gt; key {
			tmpHigh--
		}
		//找到大于key的元素，该元素的位置一定是low到tmpHigh&#43;1之间。因为array[tmpHigh&#43;1]必定大于key
		for array[tmpLow] &lt;= key &amp;&amp; tmpLow &lt; tmpHigh {
			tmpLow&#43;&#43;
		}

		if tmpLow &gt;= tmpHigh {
			break
		}
		// swap(array[tmpLow], array[tmpHigh])
		array[tmpLow], array[tmpHigh] = array[tmpHigh], array[tmpLow]
		fmt.Println(array)
	}
	array[tmpLow], array[low] = array[low], array[tmpLow]
	return tmpLow
}
