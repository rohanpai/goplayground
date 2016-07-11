// Нам часто необходимо выполнить операции над данными в коллекции,
// например выбрать все значения, удовлетворяющие какое-то условие
// или соотнести все значения в новую коллекцию с пользовательской
// функцией.

// В некоторых языках, для таких задач, идеоматически использовать [generic](http://en.wikipedia.org/wiki/Generic_programming)
// структуры данных и алгоритмы. Go не поддерживает `generics`; в Go,
// как правило, предоставляют функции коллекции если они необходимы
// конкретно для вашей программы и ваших типов данных.

// Вот несколько примеров функций коллекции для слайсов
// со строковыми значениями. Вы можете использовать эти примеры,
// чтобы сделать собственные функции. Обратите внимание,
// что в некоторых случаях, возможно, было бы более явным встроить код,
// манипулирующий коллекциями, вместо создания создания впомогательных
// функций.
package main

import &#34;strings&#34;
import &#34;fmt&#34;

// Возвращает первый индекс совпадения со строкой `t` или
// -1 если совпадение не найдено.
func Index(vs []string, t string) int {
    for i, v := range vs {
        if v == t {
            return i
        }
    }
    return -1
}

// Возвращает `true` если строка `t` находится в слайсе
func Include(vs []string, t string) bool {
    return Index(vs, t) &gt;= 0
}

// Возвращает `true` если одна из строк в слайсе
// удовлетворяет условие `f`
func Any(vs []string, f func(string) bool) bool {
    for _, v := range vs {
        if f(v) {
            return true
        }
    }
    return false
}

// Возвращает `true` если все из строк в слайсе
// удовлетворяют условие `f`
func All(vs []string, f func(string) bool) bool {
    for _, v := range vs {
        if !f(v) {
            return false
        }
    }
    return true
}

// Возвращает новый слайс, содержащий все строки в
// слайсе, которые удовлетворяют условие `f`
func Filter(vs []string, f func(string) bool) []string {
    vsf := make([]string, 0)
    for _, v := range vs {
        if f(v) {
            vsf = append(vsf, v)
        }
    }
    return vsf
}

// Возвращает новый слайс, содержащий результаты выполнения
// функции `f` с каждой строкой исходного слайса
func Map(vs []string, f func(string) string) []string {
    vsm := make([]string, len(vs))
    for i, v := range vs {
        vsm[i] = f(v)
    }
    return vsm
}

func main() {

    // Здесь мы пробуем различные фукнции коллекций.
    var strs = []string{&#34;peach&#34;, &#34;apple&#34;, &#34;pear&#34;, &#34;plum&#34;}

    fmt.Println(Index(strs, &#34;pear&#34;))

    fmt.Println(Include(strs, &#34;grape&#34;))

    fmt.Println(Any(strs, func(v string) bool {
        return strings.HasPrefix(v, &#34;p&#34;)
    }))

    fmt.Println(All(strs, func(v string) bool {
        return strings.HasPrefix(v, &#34;p&#34;)
    }))

    fmt.Println(Filter(strs, func(v string) bool {
        return strings.Contains(v, &#34;e&#34;)
    }))

    // Все примеры, приведенные выше, использовали анонимные функции,
    // но вы можете использовать именные функции корректного типа.
    fmt.Println(Map(strs, strings.ToUpper))

}
