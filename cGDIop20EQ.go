package main

import &#34;log&#34;
import zmq &#34;github.com/pebbe/zmq3&#34;
import &#34;fmt&#34;
import &#34;time&#34;
import &#34;sync&#34;

//create router socket and just print received messages
func rrecv(socket *zmq.Socket) {
    log.Println(&#34;hooking up router...&#34;)
    socket.SetIdentity(&#34;router&#34;) 
    socket.Bind(&#34;tcp://127.0.0.1:9999&#34;)
    
    for {
        header, _ := socket.Recv(0)
        body, _ := socket.Recv(0)
        
        log.Println(&#34;router received from &#34;, header, &#34; body: &#34;, body)
    }
}

    
//messages received from router
//router doesnt send anything this will not output
func drecv(socket *zmq.Socket, m *sync.Mutex) {
    poller := zmq.NewPoller()
    poller.Add(socket, zmq.POLLIN)
    for {
        sockets, _ := poller.Poll(-1)
        for _, socket := range sockets {
            m.Lock()
            msg, _ := socket.Socket.Recv(0)
            log.Println(&#34;drecv &#34;, msg)
            m.Unlock()
        }
    }
}   

//just pumps messages via channel to the hub
//to send messages to router
func dealer(id string, ch chan string) {
    for {
        ch &lt;- fmt.Sprintf(&#34;%s:*************************&#34;, id)
    }
}   

func hub(ch chan string, m *sync.Mutex) {
    dealer, _ := zmq.NewSocket(zmq.DEALER)
    dealer.SetIdentity(&#34;dealer&#34;)
    dealer.Connect(&#34;tcp://127.0.0.1:9999&#34;)
    go drecv(dealer, m)
    
    for { 
        msg := &lt;- ch
        log.Println(&#34;hub: &#34;, msg)
        m.Lock()
        dealer.SendMessageDontwait(msg)
        m.Unlock()
    }
}


func main() {
    m := &amp;sync.Mutex{}
    //create router
    router, _ := zmq.NewSocket(zmq.ROUTER)
    //deal, _ := zmq.NewSocket(zmq.DEALER)
    ch := make(chan string)

    go hub(ch, m)
    go rrecv(router)
    go dealer(&#34;dealerA&#34;, ch)
    go dealer(&#34;dealerB&#34;, ch)
    go dealer(&#34;dealerC&#34;, ch)

    time.Sleep(3000 * time.Second)
}