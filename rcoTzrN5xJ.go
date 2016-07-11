// La escritura de archivos en Go sigue patrones similares a los que vimos
// anteriormente para lectura.

package main

import (
    &#34;bufio&#34;
    &#34;fmt&#34;
    &#34;io/ioutil&#34;
    &#34;os&#34;
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {

    // Aquí vemos como guardar una cadena de bytes ([]bytes) a un archivo.
    d1 := []byte(&#34;hello\ngo\n&#34;)
    err := ioutil.WriteFile(&#34;/tmp/dat1&#34;, d1, 0644)
    check(err)

    // Para escribir de manera más controlada podemos intentar abrir un
    // archivo en modo de escritura.
    f, err := os.Create(&#34;/tmp/dat2&#34;)
    check(err)

    // Es idiomático postergar la llamada a `Close` inmediatamente después de
    // abrir un archivo, para esto usamos `defer`.
    defer f.Close()

    // Es posible escribir una secuencia de `byte` (`[]byte`) usando
    // `Write`.
    d2 := []byte{115, 111, 109, 101, 10}
    n2, err := f.Write(d2)
    check(err)
    fmt.Printf(&#34;wrote %d bytes\n&#34;, n2)

    // La función `WriteString` permite escribir usando tipos `string`
    // en vez de `[]byte`.
    n3, err := f.WriteString(&#34;writes\n&#34;)
    fmt.Printf(&#34;wrote %d bytes\n&#34;, n3)

    // Usamos `Sync` para asegurarnos que las escrituras solicitadas
    // han sido ejecutadas.
    f.Sync()

    // El paquete `bufio` provee un búfer para escritura, muy útil en conjunto
    // con los búfers de lectura que vimos anteriormente.
    w := bufio.NewWriter(f)
    n4, err := w.WriteString(&#34;buffered\n&#34;)
    fmt.Printf(&#34;wrote %d bytes\n&#34;, n4)

    // Finalmente, usamos `Flush` para asegurarnos que todas las
    // operaciones en búfer han sido aplicadas al writer correspondiente.
    w.Flush()

}
