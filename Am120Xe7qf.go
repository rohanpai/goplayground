// Commento su riga singola
/* Commento
 su riga multipla */

// In cima a ogni file è necessario specificare il package.
// Main è un package speciale che identifica un eseguibile anziché una libreria.
package main

// Con import sono dichiarate tutte le librerie a cui si fa riferimento 
// all&#39;interno del file.
import (
	&#34;fmt&#34;       // Un package nella libreria standard di Go.
	&#34;io/ioutil&#34; // Implementa alcune funzioni di utility per l&#39;I/O.
	m &#34;math&#34;    // Libreria matematica, con alias locale m
	&#34;net/http&#34;  // Sì, un web server!
	&#34;strconv&#34;   // Package per la conversione di stringhe.
)

// Una definizione di funzione. Il main è speciale: è il punto di ingresso
// per il programma. Amalo o odialo, ma Go usa le parentesi graffe.
func main() {
    // Println stampa una riga a schermo.
    // Questa funzione è all&#39;interno del package fmt.
	fmt.Println(&#34;Ciao mondo!&#34;)

    // Chiama un&#39;altra funzione all&#39;interno di questo package.
	oltreIlCiaoMondo()
}

// Le funzioni ricevono i parametri all&#39;interno di parentesi tonde.
// Se la funzione non riceve parametri, vanno comunque messe le parentesi (vuote).
func oltreIlCiaoMondo() {
	var x int // Dichiarazione di una variabile. Ricordati di dichiarare sempre le variabili prima di usarle!
	x = 3     // Assegnazione di una variabile.
    // E&#39; possibile la dichiarazione &#34;rapida&#34; := per inferire il tipo, dichiarare e assegnare contemporaneamente.
	y := 4
    // Una funzione che ritorna due valori.
	somma, prod := imparaMoltepliciValoriDiRitorno(x, y)
	fmt.Println(&#34;somma:&#34;, somma, &#34;prodotto:&#34;, prod)    // Semplice output.
	imparaTipi()                                      // &lt; y minuti, devi imparare ancora!
}

/* &lt;- commento su righe multiple
Le funzioni possono avere parametri e ritornare (molteplici!) valori.
Qua, x e y sono gli argomenti, mentre somma e prod sono i valori ritornati.
Da notare il fatto che x e somma vengono dichiarati come interi.
*/
func imparaMoltepliciValoriDiRitorno(x, y int) (somma, prod int) {
	return x &#43; y, x * y // Ritorna due valori.
}

// Ecco alcuni tipi presenti in Go
func imparaTipi() {
    // La dichiarazione rapida di solito fa il suo lavoro.
	str := &#34;Impara il Go!&#34; // Tipo stringa.

	s2 := `Una stringa letterale
puo&#39; includere andata a capo.` // Sempre di tipo stringa.

    // Stringa letterale non ASCII. I sorgenti Go sono in UTF-8.
	g := &#39;Σ&#39; // Il tipo runa, alias per int32, è costituito da un code point unicode.

	f := 3.14195 // float64, un numero in virgola mobile a 64-bit (IEEE-754)

	c := 3 &#43; 4i  // complex128, rappresentato internamente con due float64.

    // Inizializzare le variabili con var.
	var u uint = 7 // Senza segno, ma la dimensione dipende dall&#39;implementazione (come l&#39;int)
	var pi float32 = 22. / 7 

    // Sintassi per la conversione.
	n := byte(&#39;\n&#39;) // Il tipo byte è un alias per uint8.

    // I vettori hanno dimensione fissa, stabilita durante la compilazione.
	var a4 [4]int           // Un vettore di 4 interi, tutti inizializzati a 0.
	a3 := [...]int{3, 1, 5} // Un vettore inizializzato con una dimensione fissa pari a 3, i cui elementi sono 3, 1 e 5.

    // Gli slice hanno dimensione variabile. Vettori e slice hanno pro e contro,
    // ma generalmente si tende a usare più spesso gli slice.
	s3 := []int{4, 5, 9}    // La differenza con a3 è che qua non ci sono i 3 punti all&#39;interno delle parentesi quadre.
	s4 := make([]int, 4)    // Alloca uno slice di 4 interi, tutti inizializzati a 0.
	var d2 [][]float64      // Semplice dichiarazione, non vengono fatte allocazioni.
	bs := []byte(&#34;uno slice&#34;) // Sintassi per la conversione.

    // Poiché gli slice sono dinamici, è possibile aggiungere elementi
    // quando è necessario. Per farlo, si usa la funzione append(). Il primo
    // argomento è lo slice a cui stiamo aggiungendo elementi. Di solito
    // lo slice viene aggiornato, senza fare una copia, come nell&#39;esempio:
	s := []int{1, 2, 3}		// Il risultato è uno slice di dimensione 3.
	s = append(s, 4, 5, 6)	// Aggiunge 3 elementi: lo slice ha dimensione 6.
	fmt.Println(s) // Lo slice aggiornato è [1 2 3 4 5 6]
    // Per aggiungere un altro slice, invece che elencare gli elementi uno ad
    // uno, è possibile passare alla funzione append un riferimento ad uno
    // slice, oppure uno slice letterale: in questo caso si usano i tre punti,
    // dopo lo slice, a significare &#34;prendi ciascun elemento dello slice&#34;:
	s = append(s, []int{7, 8, 9}...) // Il secondo argomento è uno slice letterale.
	fmt.Println(s)	// Lo slice aggiornato è [1 2 3 4 5 6 7 8 9]

	p, q := imparaLaMemoria() // Dichiara due puntatori a intero: p e q.
	fmt.Println(*p, *q)   // * dereferenzia un puntatore. Questo stampa due interi.

    // Una variabile di tipo map è un vettore associativo di dimensione variabile,
    // e funzionano come le tabelle di hash o i dizionari in altri linguaggi.
	m := map[string]int{&#34;tre&#34;: 3, &#34;quattro&#34;: 4}
	m[&#34;uno&#34;] = 1

    // Le variabili dichiarate e non usate sono un errore in Go.
    // L&#39;underscore permette di &#34;usare&#34; una variabile, scartandone il valore.
	_, _, _, _, _, _, _, _, _, _ = str, s2, g, f, u, pi, n, a3, s4, bs
	// Stampare a schermo ovviamente significa usare una variabile.
	fmt.Println(s, c, a4, s3, d2, m)

	imparaControlloDiFlusso() // Torniamo in carreggiata.
}

// In Go è possibile associare dei nomi ai valori di ritorno di una funzione.
// Assegnare un nome al tipo di dato ritornato permette di fare return in vari
// punti all&#39;interno del corpo della funzione, ma anche di usare return senza
// specificare in modo esplicito che cosa ritornare.
func imparaValoriDiRitornoConNome(x, y int) (z int) {
	z = x * y
	return // z è implicito, perchè compare nella definizione di funzione.
}

// Go è dotato di garbage collection. Ha i puntatori, ma non l&#39;aritmetica dei
// puntatori. Puoi fare errori coi puntatori a nil, ma non puoi direttamente
// incrementare un puntatore.
func imparaLaMemoria() (p, q *int) {
    // I valori di ritorno (con nome) p e q sono puntatori a int.
	p = new(int) // La funzione new si occupa di allocare memoria.
    // L&#39;int allocato viene inizializzato a 0, dunque p non è più nil.
	s := make([]int, 20) // Alloca 20 int come un singolo blocco di memoria.
	s[3] = 7             // Ne assegna uno.
	r := -2              // Dichiara un&#39;altra variabile locale
	return &amp;s[3], &amp;r     // &amp; &#34;prende&#34; l&#39;indirizzo di un oggetto.
}

func calcoloCostoso() float64 {
	return m.Exp(10)
}

func imparaControlloDiFlusso() {
    // L&#39;istruzione if richiede parentesi graffe per il corpo, mentre non ha
    // bisogno di parentesi tonde per la condizione.
	if true {
		fmt.Println(&#34;te l&#39;ho detto&#34;)
	}
    // Eseguendo &#34;go fmt&#34; da riga di comando, il codice viene formattato
    // in maniera standard.
	if false {
		// :(
	} else {
		// :D
	}
    // L&#39;istruzione switch serve ad evitare tanti if messi in cascata.
	x := 42.0
	switch x {
	case 0:
	case 1:
	case 42:
		// Quando è soddisfatta la condizione all&#39;interno di un case, il
        // programma esce dal switch senza che siano specificate istruzioni
        // di tipo &#34;break&#34;. In Go infatti di default non è presente il
        // cosiddetto &#34;fall through&#34; all&#39;interno dell&#39;istruzione switch.
        // Tuttavia, il linguaggio mette a disposizione la parola chiave
        // fallthrough per permettere, in casi particolari, questo comportamento.
	case 43:
		// Non si arriva qua.
	default:
		// Il caso di default è opzionale.
	}
    // Come l&#39;if, anche il for non usa parentesi tonde per la condizione.
    // Le variabili dichiarate all&#39;interno di if/for sono locali al loro scope.
	for x := 0; x &lt; 3; x&#43;&#43; { // &#43;&#43; è un&#39;istruzione!
		fmt.Println(&#34;ciclo numero&#34;, x)
	}
	// x == 42 qua.

    // Il for è l&#39;unica istruzione per ciclare in Go, ma ha varie forme.
	for { // Ciclo infinito.
		break    // Scherzavo.
		continue // Non si arriva qua.
	}

    // Puoi usare range per ciclare su un vettore, slice, stringa, mappa o canale.
    // range ritorna uno (per i canali) o due valori (vettore, slice, stringa, mappa).
	for chiave, valore := range map[string]int{&#34;uno&#34;: 1, &#34;due&#34;: 2, &#34;tre&#34;: 3} {
        // per ogni coppia dentro la mappa, stampa chiave e valore
		fmt.Printf(&#34;chiave=%s, valore=%d\n&#34;, chiave, valore)
	}

    // Come nel for, := dentro la condizione dell&#39;if è usato per dichiarare
    // e assegnare y, poi testare se y &gt; x.
	if y := calcoloCostoso(); y &gt; x {
		x = y
	}
	// Le funzioni letterali sono closure.
	xGrande := func() bool {
		return x &gt; 10000 // Si riferisce a x dichiarata sopra al switch (vedi sopra).
	}
	fmt.Println(&#34;xGrande:&#34;, xGrande()) // true (abbiamo assegnato e^10 a x).
	x = 1.3e3                          // Adesso x == 1300
	fmt.Println(&#34;xGrande:&#34;, xGrande()) // false ora.

    // Inoltre le funzioni letterali possono essere definite e chiamate
    // inline, col ruolo di parametri di funzione, a patto che:
    // a) la funzione letterale venga chiamata subito (),
    // b) il valore ritornato è in accordo con il tipo dell&#39;argomento.
	fmt.Println(&#34;Somma e raddoppia due numeri: &#34;,
		func(a, b int) int {
			return (a &#43; b) * 2
		}(10, 2)) // Chiamata con argomenti 10 e 2
	// =&gt; Somma e raddoppia due numeri: 24

	// Quando ti servirà, lo amerai.
	goto amore
amore:

	imparaFabbricaDiFunzioni() // Una funzione che ritorna un&#39;altra funzione è divertente!
	imparaDefer()              // Un tour veloce di una parola chiave importante.
	imparaInterfacce()         // Arriva la roba buona!
}

func imparaFabbricaDiFunzioni() {
    // Questi due blocchi di istruzioni sono equivalenti, ma il secondo è più semplice da capire.
	fmt.Println(fabbricaDiFrasi(&#34;estate&#34;)(&#34;Una bella giornata&#34;, &#34;giornata!&#34;))

	d := fabbricaDiFrasi(&#34;estate&#34;)
	fmt.Println(d(&#34;Una bella&#34;, &#34;giornata!&#34;))
	fmt.Println(d(&#34;Un pigro&#34;, &#34;pomeriggio!&#34;))
}

// I decoratori sono comuni in alcuni linguaggi. Si può fare lo stesso in Go
// con le funzioni letterali che accettano argomenti.
func fabbricaDiFrasi(miaStringa string) func(prima, dopo string) string {
	return func(prima, dopo string) string {
		return fmt.Sprintf(&#34;%s %s %s&#34;, prima, miaStringa, dopo) // Nuova stringa
	}
}

func imparaDefer() (ok bool) {
    // Le istruzioni dette &#34;deferred&#34; (rinviate) sono eseguite
    // appena prima che la funzione ritorni.
	defer fmt.Println(&#34;le istruzioni &#39;deferred&#39; sono eseguite in ordine inverso (LIFO).&#34;)
	defer fmt.Println(&#34;\nQuesta riga viene stampata per prima perché&#34;)
    // defer viene usato di solito per chiudere un file, così la funzione che
    // chiude il file viene messa vicino a quella che lo apre.
	return true
}

// Definisce Stringer come un&#39;interfaccia con un metodo, String.
type Stringer interface {
	String() string
}

// Definisce coppia come una struct con due campi interi, chiamati x e y.
type coppia struct {
	x, y int
}

// Definisce un metodo sul tipo coppia, che adesso implementa Stringer.
func (p coppia) String() string { // p viene definito &#34;ricevente&#34;
    // Sprintf è un&#39;altra funzione del package ftm.
    // La notazione con il punto serve per richiamare i campi di p.
	return fmt.Sprintf(&#34;(%d, %d)&#34;, p.x, p.y)
}

func imparaInterfacce() {
	// Brace syntax is a &#34;struct literal&#34;. It evaluates to an initialized
	// struct. The := syntax declares and initializes p to this struct.
    // Le parentesi graffe sono usate per le cosiddette &#34;struct letterali&#34;.
    // Con :=, p viene dichiarata e inizializzata a questa struct.
	p := coppia{3, 4}
	fmt.Println(p.String()) // Chiama il metodo String di p, che è di tipo coppia.
	var i Stringer          // Dichiara i come interfaccia Stringer.
	i = p                   // Valido perchè coppia implementa Stringer.
    // Chiama il metodo String di i, che è di tipo Stringer. Output uguale a sopra.
	fmt.Println(i.String())

	// Functions in the fmt package call the String method to ask an object
	// for a printable representation of itself.
    // Le funzioni dentro al package fmt chiamano il metodo String per
    // chiedere ad un oggetto una rappresentazione in stringhe di sé stesso.
	fmt.Println(p) // Output uguale a sopra. Println chiama il metodo String.
	fmt.Println(i) // Output uguale a sopra.

	imparaParametriVariadici(&#34;grande&#34;, &#34;imparando&#34;, &#34;qua!&#34;)
}

// Le funzioni possono avere parametri variadici (ovvero di lunghezza variabile).
func imparaParametriVariadici(mieStringhe ...interface{}) {
    // Cicla su ogni valore variadico.
    // L&#39;underscore serve a ignorare l&#39;indice del vettore.
	for _, param := range mieStringhe {
		fmt.Println(&#34;parametro:&#34;, param)
	}

    // Passa un valore variadico come parametro variadico.
	fmt.Println(&#34;parametri:&#34;, fmt.Sprintln(mieStringhe...))

	imparaGestioneErrori()
}

func imparaGestioneErrori() {
    // La sintassi &#34;, ok&#34; è usata per indicare se qualcosa ha funzionato o no.
	m := map[int]string{3: &#34;tre&#34;, 4: &#34;quattro&#34;}
	if x, ok := m[1]; !ok { // ok sarà false perchè 1 non è dentro la mappa.
		fmt.Println(&#34;qua non c&#39;è nessuno!&#34;)
	} else {
		fmt.Print(x) // x sarebbe il valore che corrisponde alla chiave 1, se fosse nella mappa.
	}
    // Un errore non riporta soltanto &#34;ok&#34; ma è più specifico riguardo al problema.
	if _, err := strconv.Atoi(&#34;non_intero&#34;); err != nil { // _ scarta il valore
		// stampa &#39;strconv.ParseInt: parsing &#34;non_intero&#34;: invalid syntax&#39;
		fmt.Println(err)
	}
    // Approfondiremo le interfacce un&#39;altra volta. Nel frattempo,
	imparaConcorrenza()
}

// c è un canale, un oggetto per comunicare in modo concorrente e sicuro.
func inc(i int, c chan int) {
	c &lt;- i &#43; 1 // &lt;- è l&#39;operatore di &#34;invio&#34; quando un canale sta a sinistra.
}

// Useremo inc per incrementare alcuni numeri in modo concorrente.
func imparaConcorrenza() {
    // Stessa funzione usata prima per creare uno slice. Make alloca e
    // inizializza slice, mappe e canali.
	c := make(chan int)
    // Lancia tre goroutine. I numeri saranno incrementati in modo concorrente,
    // forse in parallelo se la macchina lo supporta. Tutti e tre inviano dati
    // sullo stesso canale.
	go inc(0, c) // go è un&#39;istruzione che avvia una goroutine.
	go inc(10, c)
	go inc(-805, c)
    // Legge tre risultati dal canale e li stampa a schermo.
    // Non si conosce a priori l&#39;ordine in cui i risultati arriveranno!
	fmt.Println(&lt;-c, &lt;-c, &lt;-c) // &lt;- è l&#39;operatore di &#34;ricevuta&#34; quando
    // un canale sta a destra.

	cs := make(chan string)       // Un altro canale, gestisce le stringhe.
	ccs := make(chan chan string) // Un canale che gestisce canali di stringhe.
	go func() { c &lt;- 84 }()       // Lancia una goroutine, solo per inviare un valore.
	go func() { cs &lt;- &#34;parolina&#34; }() // Stessa cosa ma per cs.
    // select è simile a switch, ma ogni case riguarda un&#39;operazione su un
    // canale. Seleziona, in modo random, uno tra i canali che sono pronti
    // a comunicare.
	select {
	case i := &lt;-c: // Il valore ricevuto può essere assegnato a una variabile,
		fmt.Printf(&#34;E&#39; un %T&#34;, i)
	case &lt;-cs: // oppure il valore ricevuto può essere scartato.
		fmt.Println(&#34;E&#39; una stringa.&#34;)
	case &lt;-ccs: // Canale vuoto, non pronto per comunicare.
		fmt.Println(&#34;Non succede niente.&#34;)
	}
    // A questo punto un valore è stato preso da c o cs. Una delle tue goroutine
    // cominciate sopra ha completato l&#39;esecuzione, l&#39;altra rimarrà bloccata.
}