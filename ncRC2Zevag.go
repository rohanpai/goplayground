// Comentario de una sola línea
/* Comentario
   multilínea */

// La cláusula `package` aparece al comienzo de cada archivo fuente.
// `main` es un nombre especial que declara un ejecutable en vez de una
// biblioteca.
package main

// La instrucción `import` declara los paquetes de bibliotecas referidos
// en este archivo.
import (
	&#34;fmt&#34;      // Un paquete en la biblioteca estándar de Go.
	m &#34;math&#34;   // Biblioteca de matemáticas con alias local m.
	&#34;strconv&#34;  // Conversiones de cadenas.
)

// Definición de una función. `main` es especial. Es el punto de entrada
// para el ejecutable. Te guste o no, Go utiliza llaves.
func main() {
	// Println imprime una línea a stdout.
	// Cualificalo con el nombre del paquete, fmt.
	fmt.Println(&#34;¡Hola mundo!&#34;)

	// Llama a otra función de este paquete.
	másAlláDelHola()
}

// Las funciones llevan parámetros entre paréntesis.
// Si no hay parámetros, los paréntesis siguen siendo obligatorios.
func másAlláDelHola() {
	var x int // Declaración de una variable.
 	          // Las variables se deben declarar antes de utilizarlas.
	x = 3     // Asignación de variable.
	// Declaración &#34;corta&#34; con := para inferir el tipo, declarar y asignar.
	y := 4
	suma, producto := aprendeMúltiple(x, y) // La función devuelve dos
											// valores.
	fmt.Println(&#34;suma:&#34;, suma, &#34;producto:&#34;, producto) // Simple salida.
	aprendeTipos()                        // &lt; y minutes, ¡aprende más!
}

// Las funciones pueden tener parámetros y (¡múltiples!) valores de
// retorno.
func aprendeMúltiple(x, y int) (suma, producto int) {
    return x &#43; y, x * y // Devuelve dos valores.
}

// Algunos tipos incorporados y literales.
func aprendeTipos() {
	// La declaración corta suele darte lo que quieres.
	s := &#34;¡Aprende Go!&#34; // tipo cadena.
	s2 := `Un tipo cadena &#34;puro&#34; puede incluir
saltos de línea.` // mismo tipo cadena

	// Literal no ASCII. Los ficheros fuente de Go son UTF-8.
	g := &#39;Σ&#39; // Tipo rune, un alias de int32, alberga un carácter unicode.
	f := 3.14195 // float64, el estándar IEEE-754 de coma flotante 64-bit.
	c := 3 &#43; 4i  // complex128, representado internamente por dos float64.
	// Sintaxis Var con iniciadores.
	var u uint = 7 // Sin signo, pero la implementación depende del tamaño
	               // como en int.
	var pi float32 = 22. / 7

	// Sintáxis de conversión con una declaración corta.
	n := byte(&#39;\n&#39;) // byte es un alias para uint8.

	// Los Arreglos tienen un tamaño fijo a la hora de compilar.
	var a4 [4]int           // Un arreglo de 4 ints, iniciados a 0.
	a3 := [...]int{3, 1, 5} // Un arreglo iniciado con un tamaño fijo de tres
							// elementos, con valores 3, 1 y 5.
	// Los Sectores tienen tamaño dinámico. Los arreglos y sectores tienen
	// sus ventajas y desventajas pero los casos de uso para los sectores
	// son más comunes.
	s3 := []int{4, 5, 9}     // Comparar con a3. No hay puntos suspensivos.
	s4 := make([]int, 4)     // Asigna sectores de 4 ints, iniciados a 0.
	var d2 [][]float64       // Solo declaración, sin asignación.
	bs := []byte(&#34;a sector&#34;) // Sintaxis de conversión de tipo.
	// Debido a que son dinámicos, los sectores pueden crecer bajo demanda.
	// Para añadir elementos a un sector, se utiliza la función incorporada
	// append().
	// El primer argumento es el sector al que se está anexando. Comúnmente,
	// la variable del arreglo se actualiza en su lugar, como en el 
	// siguiente ejemplo.
	sec := []int{1, 2 , 3}      // El resultado es un sector de longitud 3.
	sec = append(sec, 4, 5, 6)  // Añade 3 elementos. El sector ahora tiene una
								// longitud de 6.
	fmt.Println(sec) // El sector actualizado ahora es [1 2 3 4 5 6]
	// Para anexar otro sector, en lugar de la lista de elementos atómicos
	// podemos pasar una referencia a un sector o un sector literal como
	// este, con elipsis al final, lo que significa tomar un sector y
	// desempacar sus elementos, añadiéndolos al sector sec.
	sec = append(sec, []int{7, 8, 9} ...) // El segundo argumento es un
										  // sector literal.
	fmt.Println(sec)  // El sector actualizado ahora es [1 2 3 4 5 6 7 8 9]
	p, q := aprendeMemoria() // Declara p, q para ser un tipo puntero a
							 // int.
	fmt.Println(*p, *q)      // * sigue un puntero. Esto imprime dos ints.

	// Los Mapas son arreglos asociativos dinámicos, como los hash o
	// diccionarios de otros lenguajes.
	m := map[string]int{&#34;tres&#34;: 3, &#34;cuatro&#34;: 4}
	m[&#34;uno&#34;] = 1

	// Las variables no utilizadas en Go producen error.
	// El guión bajo permite &#34;utilizar&#34; una variable, pero descartar su
	// valor.
	_, _, _, _, _, _, _, _, _ = s2, g, f, u, pi, n, a3, s4, bs
	// Esto cuenta como utilización de variables.
	fmt.Println(s, c, a4, s3, d2, m)

	aprendeControlDeFlujo() // Vuelta al flujo.
}

// Es posible, a diferencia de muchos otros lenguajes tener valores de
// retorno con nombre en las funciones.
// Asignar un nombre al tipo que se devuelve en la línea de declaración de
// la función nos permite volver fácilmente desde múltiples puntos en una
// función, así como sólo utilizar la palabra clave `return`, sin nada
// más.
func aprendeRetornosNombrados(x, y int) (z int) {
	z = x * y
	return // aquí z es implícito, porque lo nombramos antes.
}

// Go posee recolector de basura. Tiene punteros pero no aritmética de
// punteros. Puedes cometer errores con un puntero nil, pero no
// incrementando un puntero.
func aprendeMemoria() (p, q *int) {
	// Los valores de retorno nombrados q y p tienen un tipo puntero
	// a int.
	p = new(int) // Función incorporada que reserva memoria.
	// La asignación de int se inicia a 0, p ya no es nil.
	s := make([]int, 20) // Reserva 20 ints en un solo bloque de memoria.
	s[3] = 7             // Asigna uno de ellos.
	r := -2              // Declara otra variable local.
	return &amp;s[3], &amp;r     // &amp; toma la dirección de un objeto.
}

func cálculoCaro() float64 {
	return m.Exp(10)
}

func aprendeControlDeFlujo() {
	// La declaración If requiere llaves, pero no paréntesis.
	if true {
		fmt.Println(&#34;ya relatado&#34;)
	}
	// El formato está estandarizado por la orden &#34;go fmt.&#34;
	if false {
		// Abadejo.
	} else {
		// Relamido.
	}
	// Utiliza switch preferentemente para if encadenados.
	x := 42.0
	switch x {
	case 0:
	case 1:
	case 42:
		// Los cases no se mezclan, no requieren de &#34;break&#34;.
	case 43:
		// No llega.
	}
	// Como if, for no utiliza paréntesis tampoco.
	// Variables declaradas en for e if son locales a su ámbito.
	for x := 0; x &lt; 3; x&#43;&#43; { // &#43;&#43; es una instrucción.
		fmt.Println(&#34;iteración&#34;, x)
	}
	// aquí x == 42.

	// For es la única instrucción de bucle en Go, pero tiene formas
	// alternativas.
	for { // Bucle infinito.
		break    // ¡Solo bromeaba!
		continue // No llega.
	}

	// Puedes usar `range` para iterar en un arreglo, un sector, una
	// cadena, un mapa o un canal.
	// `range` devuelve o bien, un canal o de uno a dos values (arreglo,
	// sector, cadena y mapa).
	for clave, valor := range map[string]int{&#34;uno&#34;: 1, &#34;dos&#34;: 2, &#34;tres&#34;: 3} {
		// por cada par en el mapa, imprime la clave y el valor
		fmt.Printf(&#34;clave=%s, valor=%d\n&#34;, clave, valor)
	}

	// Como en for, := en una instrucción if significa declarar y asignar
	// primero, luego comprobar y &gt; x.
	if y := cálculoCaro(); y &gt; x {
		x = y
    }
	// Las funciones literales son &#34;cierres&#34;.
	granX := func() bool {
		return x &gt; 100 // Referencia a x declarada encima de la instrucción
							// switch.
	}
	fmt.Println(&#34;granX:&#34;, granX()) // cierto (la última vez asignamos
											 // 1e6 a x).
	x /= 1.3e3                   // Esto hace a x == 1300
	fmt.Println(&#34;granX:&#34;, granX()) // Ahora es falso.

	// Es más las funciones literales se pueden definir y llamar en línea,
	// actuando como un argumento para la función, siempre y cuando:
	// a) la función literal sea llamada inmediatamente (),
	// b) el tipo del resultado sea del tipo esperado del argumento
	fmt.Println(&#34;Suma dos números &#43; doble: &#34;,
		func(a, b int) int {
			return (a &#43; b) * 2
		}(10, 2)) // Llamada con argumentos 10 y 2
	// =&gt; Suma dos números &#43; doble: 24

	// Cuando lo necesites, te encantará.
	goto encanto
encanto:

	aprendeFunciónFábrica() // func devolviendo func es divertido(3)(3)
	aprendeADiferir()       // Un rápido desvío a una importante palabra clave.
	aprendeInterfaces()     // ¡Buen material dentro de poco!
}

func aprendeFunciónFábrica() {
	// Las dos siguientes son equivalentes, la segunda es más práctica
	fmt.Println(instrucciónFábrica(&#34;día&#34;)(&#34;Un bello&#34;, &#34;de verano&#34;))

	d := instrucciónFábrica(&#34;atardecer&#34;)
	fmt.Println(d(&#34;Un maravilloso&#34;, &#34;de verano&#34;))
	fmt.Println(d(&#34;Un hermoso&#34;, &#34;de verano&#34;))
}

// Los decoradores son comunes en otros languajes. Lo mismo se puede hacer
// en Go con funciónes literales que aceptan argumentos.
func instrucciónFábrica(micadena string) func(antes, después string) string {
	return func(antes, después string) string {
		return fmt.Sprintf(&#34;¡%s %s %s!&#34;, antes, micadena, después) // nueva cadena
	}
}

func aprendeADiferir() (ok bool) {
	// las instrucciones diferidas se ejecutan justo antes de que la
	// función regrese.
	defer fmt.Println(&#34;las instrucciones diferidas se ejecutan en orden inverso (PEPS).&#34;)
	defer fmt.Println(&#34;\nEsta línea se imprime primero debido a que&#34;)
	// Defer se usa comunmente para cerrar un fichero, por lo que la
	// función que cierra el fichero se mantiene cerca de la función que lo
	// abrió.
	return true
}

// Define Stringer como un tipo interfaz con un método, String.
type Stringer interface {
	String() string
}

// Define par como una estructura con dos campos int, x e y.
type par struct {
	x, y int
}

// Define un método en el tipo par. Par ahora implementa a Stringer.
func (p par) String() string { // p se conoce como el &#34;receptor&#34;
	// Sprintf es otra función pública del paquete fmt.
	// La sintaxis con punto se refiere a los campos de p.
	return fmt.Sprintf(&#34;(%d, %d)&#34;, p.x, p.y)
}

func aprendeInterfaces() {
	// La sintaxis de llaves es una &#34;estructura literal&#34;. Evalúa a una
	// estructura iniciada. La sintaxis := declara e inicia p a esta
	// estructura.
	p := par{3, 4}
	fmt.Println(p.String()) // Llama al método String de p, de tipo par.
	var i Stringer          // Declara i como interfaz de tipo Stringer.
	i = p                   // Válido porque par implementa Stringer.
	// Llama al metodo String de i, de tipo Stringer. Misma salida que
	// arriba.
	fmt.Println(i.String())

	// Las funciones en el paquete fmt llaman al método String para
	// consultar un objeto por una representación imprimible de si
	// mismo.
	fmt.Println(p) // Salida igual que arriba. Println llama al método
				   // String.
	fmt.Println(i) // Salida igual que arriba.
	aprendeNúmeroVariableDeParámetros(&#34;¡gran&#34;, &#34;aprendizaje&#34;, &#34;aquí!&#34;)
}

// Las funciones pueden tener número variable de argumentos.
func aprendeNúmeroVariableDeParámetros(misCadenas ...interface{}) {
	// Itera en cada valor de los argumentos variables.
	// El espacio en blanco aquí omite el índice del argumento arreglo.
	for _, parámetro := range misCadenas {
		fmt.Println(&#34;parámetro:&#34;, parámetro)
	}

	// Pasa el valor de múltiples variables como parámetro variadic.
	fmt.Println(&#34;parámetros:&#34;, fmt.Sprintln(misCadenas...))
	aprendeManejoDeError()
}

func aprendeManejoDeError() {
	// &#34;, ok&#34; forma utilizada para saber si algo funcionó o no.
	m := map[int]string{3: &#34;tres&#34;, 4: &#34;cuatro&#34;}
	if x, ok := m[1]; !ok { // ok será falso porque 1 no está en el mapa.
		fmt.Println(&#34;nada allí&#34;)
	} else {
		fmt.Print(x) // x sería el valor, si estuviera en el mapa.
	}
	// Un valor de error comunica más información sobre el problema aparte
	// de &#34;ok&#34;.
	if _, err := strconv.Atoi(&#34;no-int&#34;); err != nil { // _ descarta el
																	  // valor
		// Imprime &#34;strconv.ParseInt: parsing &#34;no-int&#34;: invalid syntax&#34;.
		fmt.Println(err)
	}
	// Revisaremos las interfaces más adelante. Mientras tanto...
	aprendeConcurrencia()
}

// c es un canal, un objeto de comunicación concurrente seguro.
func inc(i int, c chan int) {
	c &lt;- i &#43; 1 // &lt;- es el operador &#34;enviar&#34; cuando aparece un canal a la
				  // izquierda.
}

// Utilizaremos inc para incrementar algunos números concurrentemente.
func aprendeConcurrencia() {
	// Misma función make utilizada antes para crear un sector. Make asigna
	// e inicia sectores, mapas y canales.
	c := make(chan int)
	// Inicia tres rutinasgo concurrentes. Los números serán incrementados
	// concurrentemente, quizás en paralelo si la máquina es capaz y está
	// correctamente configurada. Las tres envían al mismo canal.
	go inc(0, c) // go es una instrucción que inicia una nueva rutinago.
	go inc(10, c)
	go inc(-805, c)
	// Lee los tres resultados del canal y los imprime.
	// ¡No se puede saber en que orden llegarán los resultados!
	fmt.Println(&lt;-c, &lt;-c, &lt;-c) // Canal a la derecha, &lt;- es el operador
										// &#34;recibe&#34;.

	cs := make(chan string)       // Otro canal, este gestiona cadenas.
	ccs := make(chan chan string) // Un canal de canales cadena.
	go func() { c &lt;- 84 }()       // Inicia una nueva rutinago solo para
											// enviar un valor.
	go func() { cs &lt;- &#34;verboso&#34; }() // Otra vez, para cs en esta ocasión.
	// Select tiene una sintáxis parecida a la instrucción switch pero cada
	// caso involucra una operacion con un canal. Selecciona un caso de
	// forma aleatoria de los casos que están listos para comunicarse.
	select {
	case i := &lt;-c: // El valor recibido se puede asignar a una variable,
		fmt.Printf(&#34;es un %T&#34;, i)
	case &lt;-cs:     // o el valor se puede descartar.
		fmt.Println(&#34;es una cadena&#34;)
	case &lt;-ccs:    // Canal vacío, no está listo para la comunicación.
		fmt.Println(&#34;no sucedió.&#34;)
	}

	// En este punto un valor fue devuelto de c o cs. Una de las dos
	// rutinasgo que se iniciaron se ha completado, la otrá permancerá
	// bloqueada.
}
