package main

import "log"
import zmq "github.com/pebbe/zmq3"
import "fmt"
import "time"
import "sync"

//create router socket and just print received messages
func rrecv(socket *zmq.Socket) {
    log.Println("hooking up router...")
    socket.SetIdentity("router") 
    socket.Bind("tcp://127.0.0.1:9999")
    
    for {
        header, _ := socket.Recv(0)
        body, _ := socket.Recv(0)
        
        log.Println("router received from ", header, " body: ", body)
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
            log.Println("drecv ", msg)
            m.Unlock()
        }
    }
}   

//just pumps messages via channel to the hub
//to send messages to router
func dealer(id string, ch chan string) {
    for {
        ch <- fmt.Sprintf("%s:*************************", id)
    }
}   

func hub(ch chan string, m *sync.Mutex) {
    dealer, _ := zmq.NewSocket(zmq.DEALER)
    dealer.SetIdentity("dealer")
    dealer.Connect("tcp://127.0.0.1:9999")
    go drecv(dealer, m)
    
    for { 
        msg := <- ch
        log.Println("hub: ", msg)
        m.Lock()
        dealer.SendMessageDontwait(msg)
        m.Unlock()
    }
}


func main() {
    m := &sync.Mutex{}
    //create router
    router, _ := zmq.NewSocket(zmq.ROUTER)
    //deal, _ := zmq.NewSocket(zmq.DEALER)
    ch := make(chan string)

    go hub(ch, m)
    go rrecv(router)
    go dealer("dealerA", ch)
    go dealer("dealerB", ch)
    go dealer("dealerC", ch)

    time.Sleep(3000 * time.Second)
}