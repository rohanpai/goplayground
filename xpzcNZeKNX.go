// Go ofrece soporte incluido para codificar y descifrar JSON,
// incluyendo desde tipos de datos incorporados hasta
// tipos de datos personalizados.

package main

import &#34;encoding/json&#34;
import &#34;fmt&#34;
import &#34;os&#34;

// Vamos a usar estas dos estructuras para demostrar cifrado y
// descifrado de los tipos personalizados mostrados a continuación.
type Respuesta1 struct {
    Pagina int
    Frutas []string
}
type Respuesta2 struct {
    Pagina int      `json:&#34;pagina&#34;`
    Frutas []string `json:&#34;frutas&#34;`
}

func main() {

    // Primero, vamos a hechar un vistazo al cifrado basico de tipos de datos
    // a cadenas de JSON. Aqui hay algunos ejemplos para valores atómicos.
    bolB, _ := json.Marshal(true)
    fmt.Println(string(bolB))

    intB, _ := json.Marshal(1)
    fmt.Println(string(intB))

    fltB, _ := json.Marshal(2.34)
    fmt.Println(string(fltB))

    strB, _ := json.Marshal(&#34;gopher&#34;)
    fmt.Println(string(strB))

    // Y aqui hay algunos para mapas y porciones(slices), donde
    // se cifra a cadenas de JSON y objetos como se esperaría.
    slcD := []string{&#34;manzana&#34;, &#34;durazno&#34;, &#34;pera&#34;}
    slcB, _ := json.Marshal(slcD)
    fmt.Println(string(slcB))

    mapD := map[string]int{&#34;manzana&#34;: 5, &#34;lechuga&#34;: 7}
    mapB, _ := json.Marshal(mapD)
    fmt.Println(string(mapB))

    // El paquete JSON puede cifrar automaticamente tus
    // tipos de datos personalizados. Solo incluirá campos
    // exportados en el cifrado de salida y por defecto
    // va a utilizar esos nombres como las llaves del JSON.
    res1D := &amp;Respuesta1{
        Pagina: 1,
        Frutas: []string{&#34;manzana&#34;, &#34;durazno&#34;, &#34;pera&#34;}}
    res1B, _ := json.Marshal(res1D)
    fmt.Println(string(res1B))

    // Puedes usar etiquetas en las declaraciones de campos
    // de estructuras para personalizar los nombres de las llaves
    // del JSON cifrado. Mira la definición previa de `Respuesta2`
    // para ver un ejemplo de esas etiquetas.
    res2D := &amp;Respuesta2{
        Pagina: 1,
        Frutas: []string{&#34;manzana&#34;, &#34;durazno&#34;, &#34;pera&#34;}}
    res2B, _ := json.Marshal(res2D)
    fmt.Println(string(res2B))

    // Ahora echemos un vistazo al decifrado de datos de JSON
    // a valores de Go. Aqui hay un ejemplo para una estructura
    // genérica de datos.
    byt := []byte(`{&#34;num&#34;:6.13,&#34;strs&#34;:[&#34;a&#34;,&#34;b&#34;]}`)

    // Necesitamos proveer una variable donde el paquete
    // JSON pueda colocar los datos decifrados. Este
    // `map[string]interface{}` va a contener un mapa de
    // cadenas para tipos de datos arbitrarios.
    var dat map[string]interface{}

    // Aqui está el decifrado real, y una verificación
    // por errores asociados.
    if err := json.Unmarshal(byt, &amp;dat); err != nil {
        panic(err)
    }
    fmt.Println(dat)

    // A fin de usar los valores en el mapa decifrado,
    // necesitaremos emitirlos a su tipo apropiado.
    // Por ejemplo, aqui emitimos el valor en `num` al
    // tipo esperado `float64`.
    num := dat[&#34;num&#34;].(float64)
    fmt.Println(num)

    // Accesar a datos anidados necesita series de
    // emisiones.
    strs := dat[&#34;strs&#34;].([]interface{})
    str1 := strs[0].(string)
    fmt.Println(str1)

    // También podemos decifrar JSON a tipos de datos personalizados.
    // Esto tiene la ventaja de añadir seguridad adicional en el tipo
    // para nuestros programas y eliminar la necesidad de afirmar el
    // tipo al accesar los datos decifrados.
    str := `{&#34;pagina&#34;: 1, &#34;frutas&#34;: [&#34;manzana&#34;, &#34;pera&#34;]}`
    res := &amp;Respuesta2{}
    json.Unmarshal([]byte(str), &amp;res)
    fmt.Println(res)
    fmt.Println(res.Frutas[0])

    // En los ejemplos previos siempre usamos bytes y cadenas
    // como intermediarios entre los datos y la represación del JSON
    // de acuerdo al estandar. Además podemos correr cifrados JSON
    // directamente a `os.Writer`s como `os.Stdout` o incluso
    // en respuestas del cuerpo de HTTP.
    enc := json.NewEncoder(os.Stdout)
    d := map[string]int{&#34;manzana&#34;: 5, &#34;lechuga&#34;: 7}
    enc.Encode(d)
}
