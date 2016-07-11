// _Defer_ se usa para asegurar que una función es
// llamada posteriormente durante la ejecución del
// programa, generalmente con propósitos de limpieza.
// `defer` se usa regularmente donde en otros lenguajes
// se utilizaría `ensure` y `finally`

package main

import &#34;fmt&#34;
import &#34;os&#34;

// Supongamos que queremos crear un archivo, escribir
// en él y luego cerrarlo al terminar. Así es como lo
// haríamos utilizando `defer`
func main() {

    // Inmediatamente después de obtener el objeto archivo
    // con `createFile`, diferimos el cierre del archivo con
    // `closeFile`. Esto se ejecutará al término de la función
    // contenedora (`main`), después de que `writeFile`
    // terminó de ejecutarse.
    f := createFile(&#34;/tmp/defer.txt&#34;)
    defer closeFile(f)
    writeFile(f)
}

func createFile(p string) *os.File {
    fmt.Println(&#34;crear&#34;)
    f, err := os.Create(p)
    if err != nil {
        panic(err)
    }
    return f
}

func writeFile(f *os.File) {
    fmt.Println(&#34;escribir&#34;)
    fmt.Fprintln(f, &#34;data&#34;)

}

func closeFile(f *os.File) {
    fmt.Println(&#34;cerrar&#34;)
    f.Close()
}
