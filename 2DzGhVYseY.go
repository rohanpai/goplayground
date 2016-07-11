package main

import (
	&#34;fmt&#34;
)

type MenuComponent interface {
	Name() string
	Description() string
	Price() float32
	Print()

	Add(MenuComponent)
	Remove(int)
	Child(int) MenuComponent
}

type MenuItem struct {
	name        string
	description string
	price       float32
}

func NewMenuItem(name, description string, price float32) MenuComponent {
	return &amp;MenuItem{
		name:        name,
		description: description,
		price:       price,
	}
}

func (m *MenuItem) Name() string {
	return m.name
}

func (m *MenuItem) Description() string {
	return m.description
}

func (m *MenuItem) Price() float32 {
	return m.price
}

func (m *MenuItem) Print() {
	fmt.Printf(&#34;  %s, ￥%.2f\n&#34;, m.name, m.price)
	fmt.Printf(&#34;    -- %s\n&#34;, m.description)
}

func (m *MenuItem) Add(MenuComponent) {
	panic(&#34;not implement&#34;)
}

func (m *MenuItem) Remove(int) {
	panic(&#34;not implement&#34;)
}

func (m *MenuItem) Child(int) MenuComponent {
	panic(&#34;not implement&#34;)
}

type Menu struct {
	name        string
	description string
	children    []MenuComponent
}

func NewMenu(name, description string) MenuComponent {
	return &amp;Menu{
		name:        name,
		description: description,
	}
}

func (m *Menu) Name() string {
	return m.name
}

func (m *Menu) Description() string {
	return m.description
}

func (m *Menu) Price() (price float32) {
	for _, v := range m.children {
		price &#43;= v.Price()
	}
	return
}

func (m *Menu) Print() {
	fmt.Printf(&#34;%s, %s, ￥%.2f\n&#34;, m.name, m.description, m.Price())
	fmt.Println(&#34;------------------------&#34;)
	for _, v := range m.children {
		v.Print()
	}
	fmt.Println()
}

func (m *Menu) Add(c MenuComponent) {
	m.children = append(m.children, c)
}

func (m *Menu) Remove(idx int) {
	m.children = append(m.children[:idx], m.children[idx&#43;1:]...)
}

func (m *Menu) Child(idx int) MenuComponent {
	return m.children[idx]
}

func main() {
	menu1 := NewMenu(&#34;培根鸡腿燕麦堡套餐&#34;, &#34;供应时间：09:15--22:44&#34;)
	menu1.Add(NewMenuItem(&#34;主食&#34;, &#34;培根鸡腿燕麦堡1个&#34;, 11.5))
	menu1.Add(NewMenuItem(&#34;小吃&#34;, &#34;玉米沙拉1份&#34;, 5.0))
	menu1.Add(NewMenuItem(&#34;饮料&#34;, &#34;九珍果汁饮料1杯&#34;, 6.5))

	menu2 := NewMenu(&#34;奥尔良烤鸡腿饭套餐&#34;, &#34;供应时间：09:15--22:44&#34;)
	menu2.Add(NewMenuItem(&#34;主食&#34;, &#34;新奥尔良烤鸡腿饭1份&#34;, 15.0))
	menu2.Add(NewMenuItem(&#34;小吃&#34;, &#34;新奥尔良烤翅2块&#34;, 11.0))
	menu2.Add(NewMenuItem(&#34;饮料&#34;, &#34;芙蓉荟蔬汤1份&#34;, 4.5))

	all := NewMenu(&#34;超值午餐&#34;, &#34;周一至周五有售&#34;)
	all.Add(menu1)
	all.Add(menu2)

	all.Print()
}
