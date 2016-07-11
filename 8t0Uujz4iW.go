// En el ejemplo anterior vimos como manejar un estado
// simple usando operaciones atómicas. Para estados más
// complejos podemos usar un _[mutex](http://en.wikipedia.org/wiki/Mutual_exclusion)_
// (_mecanismo de exclusión mutua_) para acceder a los
// datos desde múltiples goroutines.

package main

import (
    &#34;fmt&#34;
    &#34;math/rand&#34;
    &#34;runtime&#34;
    &#34;sync&#34;
    &#34;sync/atomic&#34;
    &#34;time&#34;
)

func main() {

    // Para nuestro ejemplo el `estado` será un mapa.
    var estado = make(map[int]int)

    // Este `mutex` sincronizará el acceso al `estado`.
    var mutex = &amp;sync.Mutex{}

    // Para comparar el uso de mutex con otro
    // mecanismo que veremos después usaremos la variable
    // `ops` que contará el número de operaciones
    // hechas con el `estado`
    var ops int64 = 0

    // Aquí vamos a iniciar 100 gorutinas que leerán
    // contastantemente el `estado`
    for r := 0; r &lt; 100; r&#43;&#43; {
        go func() {
            total := 0
            for {

                // En cada lectura tomamos una llave
                // para acceder, aseguramos el acceso
                // exclusivo al `estado` llamando al
                // método `Lock()` del mutex, leemos
                // el valor de la llave elegida,
                // desbloqueamos el mutex llamando a `Unlock()`
                // e incrementamos el contador `ops`
                llave := rand.Intn(5)
                mutex.Lock()
                total &#43;= estado[llave]
                mutex.Unlock()
                atomic.AddInt64(&amp;ops, 1)

                // Para asegurar que esta gorutina
                // no asfixie al scheduler, vamos a ceder
                // explicitamente después de cada operación
                // llamando `runtime.Gosched()`. Este &#34;ceder&#34;
                // es manejado automáticamente con operaciones
                // en canales y al realizar llamadas
                // bloqueantes como `time.Sleep`, pero en
                // este caso lo tenemos que hacer manualmente.
                runtime.Gosched()
            }
        }()
    }

    // Vamos a iniciar 10 gorutinas para simular escrituras,
    // usando el mismo patrón que usamos en las lecturas.
    for w := 0; w &lt; 10; w&#43;&#43; {
        go func() {
            for {
                llave := rand.Intn(5)
                val := rand.Intn(100)
                mutex.Lock()
                estado[llave] = val
                mutex.Unlock()
                atomic.AddInt64(&amp;ops, 1)
                runtime.Gosched()
            }
        }()
    }

    // Dejamos que que las 10 gorutinas trabajen en el `estado`
    // y en el `mutex` por un segundo.
    time.Sleep(time.Second)

    // Tomamos y reportamos un conteo final de operaciones
    opsFinal := atomic.LoadInt64(&amp;ops)
    fmt.Println(&#34;ops:&#34;, opsFinal)

    // Hacemo un bloqueo final del `estado` y mostramos como
    // terminó
    mutex.Lock()
    fmt.Println(&#34;estado:&#34;, estado)
    mutex.Unlock()
}
