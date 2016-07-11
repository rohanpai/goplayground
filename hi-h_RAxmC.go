// La derivación con `if` y `else` en Go
// se hace de manera directa.

package main

import &#34;fmt&#34;

func main() {

    // Ejemplo básico.
    if 7%2 == 0 {
        fmt.Println(&#34;7 es par&#34;)
    } else {
        fmt.Println(&#34;7 es impar&#34;)
    }

    // Puedes utilizar un `if` sin un else.
    if 8%4 == 0 {
        fmt.Println(&#34;8 es divisible entre 4&#34;)
    }

    // Los condicionales pueden ser precedidos por
    // una declaración; cualquier
    // variable declarada en dicha declaración estará disponible
    // en todas las derivaciones.
    if num := 9; num &lt; 0 {
        fmt.Println(num, &#34;es negativo&#34;)
    } else if num &lt; 10 {
        fmt.Println(num, &#34;tiene 1 digito&#34;)
    } else {
        fmt.Println(num, &#34;tiene multiples digitos&#34;)
    }
}

// Nota que no necesitas los paréntesis alrededor de las
// condiciones en Go, pero las llaves {} si son obligatorias.
