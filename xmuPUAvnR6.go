package main 


import (
    "bufio"
    "bytes"
    "os"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "time"
    "encoding/json"
    "encoding/xml"
    "github.com/streadway/amqp"
    // "github.com/sethgrid/pester"
)

var exit = make(chan bool, 1)
var channels = make([]string, 0)
var client = &http.Client{}

type ACResponse struct {
    Name string `json:"name"`
    Parent string `json:"parent"`
}

type Channel struct {
    XMLName xml.Name `xml:"Channel"`
    ChannelId string `xml:"channelId"`
}

type Subscription struct {
    XMLName xml.Name `xml:"Subscription"`
    SubscriptionId string `xml:"subscriptionId"`
}

type Event struct {
    XMLName xml.Name `xml:"xsi:Event"`
    EventId string `xml:"xsi:eventID"`
}


func req(url string, data []byte) (c *http.Response){
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
    req.Header.Set("Authorization", fmt.Sprintf("Basic %s", os.Getenv("BSFT_TOKEN")))
    req.Header.Add("Content-Type", "text/plain;charset=UTF-8")
    req.Header.Add("Accept-Encoding", "identity")
    req.Header.Del("User-Agent")
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    return resp
}

func post(url string, data []byte) (c []byte){
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
    req.Header.Set("Authorization", fmt.Sprintf("Basic %s", os.Getenv("BSFT_TOKEN")))
    req.Header.Add("Content-Type", "text/plain;charset=UTF-8")
    req.Header.Add("Accept-Encoding", "identity")
    req.Header.Del("User-Agent")
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    return body
}

func subscriptionTemplate(event string, system_id string) (xml []byte) {
    return []byte(fmt.Sprintf(`XML`, event, system_id))
}

func channelTemplate(system_id string) (xml []byte) {
    return []byte(fmt.Sprintf(`XML`, system_id))
}

func refreshTemplate() (xml []byte) {
    return []byte(fmt.Sprintf(`XML`))
}

func refreshChannel() (xml []byte) {
    return []byte(fmt.Sprintf(`XML`))
}


func channelHeartbeat(channelId string) {
    req, _ := http.NewRequest("PUT", fmt.Sprintf("http://URL/channel/%s/heartbeat", channelId), nil)
    req.Header.Add("Authorization", fmt.Sprintf("Basic %s", os.Getenv("BSFT_TOKEN")))
    req.Header.Add("Content-Type", "text/plain;charset=UTF-8")
    req.Header.Add("Accept-Encoding", "identity")
    req.Header.Del("User-Agent")
    resp, _ := client.Do(req)
    body, _ := ioutil.ReadAll(resp.Body)
    if (len(body) != 0) {
        log.Println("error on heartbeat for", channelId, string(body))
    }
    resp.Body.Close()
}

func channelTicker(channelId string) {
    ticker := time.NewTicker(time.Second * 15)
    fmt.Printf("%s Starting Channel heartbeat for %s\n", time.Now(), channelId)
    for t := range ticker.C {
        fmt.Printf("Heartbeat at %s for %s\n", t, channelId)
        channelHeartbeat(channelId)
    }
}


func subscribe(ch *amqp.Channel, q amqp.Queue, group ACResponse) (scanner *bufio.Scanner, err error){
    // Stream to read chunk by chunk.  
    res := req("http://URL/channel", channelTemplate(group.Name))
    reader := bufio.NewScanner(res.Body)
    // Make subscriptions
    fmt.Println(time.Now(), "Subscribed to", group.Name)
    for reader.Scan() {
        text := reader.Text()
        c := Channel{ChannelId: "none"}
        merr := xml.Unmarshal([]byte(text), &c)
        err = ch.Publish(
            "events",     // exchange
            "events", // routing key
            false,  // mandatory
            false,  // immediate
            amqp.Publishing{
                ContentType: "text/plain",
                Body:        []byte(text),
            })
        // log.Printf(" [x] Sent %s", text)
        if merr == nil {
            fmt.Println("ChannelId", c.ChannelId)
            go channelTicker(c.ChannelId)
        }
        failOnError(err, "Failed to publish a message")
    }
    if err := reader.Err(); err != nil {
        panic(err)
        exit <- true
    }
    return reader, err
}


func failOnError(err error, msg string) {
  if err != nil {
    log.Fatalf("%s: %s", msg, err)
    panic(fmt.Sprintf("%s: %s", msg, err))
  }
}

func main() {
    request, _ := http.NewRequest("GET", "URL/subscriptions/?format=json", nil)
    request.Header.Add("Content-Type", "application/json")
    resp, _ := client.Do(request)
    var ac []ACResponse
    json.NewDecoder(resp.Body).Decode(&ac)
    resp.Body.Close()


    conn, err := amqp.Dial("amqp://URL:5672//")
    failOnError(err, "Failed to connect to RabbitMQ")
    defer conn.Close()

    ch, err := conn.Channel()
    failOnError(err, "Failed to open a channel")
    defer ch.Close()

    q, err := ch.QueueDeclare(
        "events", // name
        false,   // durable
        false,   // delete when unused
        false,   // exclusive
        false,   // no-wait
        nil,     // arguments
    )
    failOnError(err, "Failed to declare a queue")
    for _, item := range ac[:1] {
        go subscribe(ch, q, item)
        // time.Sleep(time.Second * 7)
    }

    <-exit

}