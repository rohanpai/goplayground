package main

import &#34;fmt&#34;
import &#34;sort&#34;

type User struct{
    Name string
    Balance float32
    MonthOutlay float32
}

type SortAdapter struct {
    len func() int
    less func(i, j int) bool
    swap func(i, j int)
}

func (s SortAdapter) Len() int{
    return s.len()
}

func (s SortAdapter) Less(i, j int) bool{
    return s.less(i, j)
}

func (s SortAdapter) Swap(i, j int) {
    s.swap(i, j)
}

func main() {
        users := []User{{&#34;Петя&#34;, 100, 1000}, {&#34;Коля&#34;, 50, 2000}, {&#34;Катя&#34;, 900, 500}, {&#34;Маша&#34;, 200, 0}, {&#34;Вася&#34;, 400, 700}}

	sort.Sort(SortAdapter{func() int {
    		return len(users)
	}, func(i, j int) bool {
    		return users[i].Balance &lt; users[j].Balance
	}, func(i, j int) {
    		users[i], users[j] = users[j], users[i]
	}})

	fmt.Println(&#34;Sorted by balance:\n&#34;, users)

	sort.Sort(SortAdapter{func() int {
    		return len(users)
	}, func(i, j int) bool {
    		return users[i].MonthOutlay &lt; users[j].MonthOutlay
	}, func(i, j int) {
    		users[i], users[j] = users[j], users[i]
	}})

	fmt.Println(&#34;Sorted by month outlay:\n&#34;, users)
}