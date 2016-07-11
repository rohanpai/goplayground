// Los _Slices_ son un tipo de datos en Go que proporcionan
// una interfaz más poderosa a las secuencias que los arreglos.

package main

import &#34;fmt&#34;

func main() {

    // A comparación de los arreglos, los slices son solo del tipo
    // de los elementos que contienen (no del numero de elementos).
    // Para crear un slice de tamaño cero, se usa la sentencia `make`.
    // En este ejemplo creamos un slice de `string`s de tamaño `3`
    // (inicializado con valores cero).
    s := make([]string, 3)
    fmt.Println(&#34;emp:&#34;, s)

    // Podemos establecer y obtener valores just como con los arreglos.
    s[0] = &#34;a&#34;
    s[1] = &#34;b&#34;
    s[2] = &#34;c&#34;
    fmt.Println(&#34;set:&#34;, s)
    fmt.Println(&#34;get:&#34;, s[2])

    // `len` regresa el tamaño del slice.
    fmt.Println(&#34;len:&#34;, len(s))

    // Aparte de estas operaciones básicas, los slices
    // soportan muchas mas que los hacen más funcionales
    // que los arreglos. Una de ellas es `append`, la que
    // regresa un slice que contiene uno o mas valores nuevos.
    // Nota que necesitamos asignar el valor de regreso de
    // append tal como lo haríamos con el valor de un slice nuevo.
    s = append(s, &#34;d&#34;)
    s = append(s, &#34;e&#34;, &#34;f&#34;)
    fmt.Println(&#34;apd:&#34;, s)

    // Los Slices pueden ser copiados utilizando `copy`.
    // Aquí creamos un slice vacío `c` del mismo tamaño que
    // `s` y copiamos el contenido de `s` a `c`.
    c := make([]string, len(s))
    copy(c, s)
    fmt.Println(&#34;cpy:&#34;, c)

    // Los Slices soportan un operador de rango con la sintaxis
    // `slice[low:high]`. Por ejemplo, esto regresa un slice
    // de los elementos `s[2]`, `s[3]`, y `s[4]`.
    l := s[2:5]
    fmt.Println(&#34;sl1:&#34;, l)

    // Esto regresa los elementos hasta antes de `s[5]`.
    l = s[:5]
    fmt.Println(&#34;sl2:&#34;, l)

    // y esto regresa los elementos desde `s[2]`.
    l = s[2:]
    fmt.Println(&#34;sl3:&#34;, l)

    // Podemos declarar e inicializar una variable para el slice
    // en una sola línea también.
    t := []string{&#34;g&#34;, &#34;h&#34;, &#34;i&#34;}
    fmt.Println(&#34;dcl:&#34;, t)

    // Los slices pueden ser compuestos de estructuras multi dimensionales.
    // A diferencia de los arreglos, el tamaño de los slices interiores
    // puede variar.
    twoD := make([][]int, 3)
    for i := 0; i &lt; 3; i&#43;&#43; {
        innerLen := i &#43; 1
        twoD[i] = make([]int, innerLen)
        for j := 0; j &lt; innerLen; j&#43;&#43; {
            twoD[i][j] = i &#43; j
        }
    }
    fmt.Println(&#34;2d: &#34;, twoD)
}
