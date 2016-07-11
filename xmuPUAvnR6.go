package main 


import (
    &#34;bufio&#34;
    &#34;bytes&#34;
    &#34;os&#34;
    &#34;fmt&#34;
    &#34;io/ioutil&#34;
    &#34;log&#34;
    &#34;net/http&#34;
    &#34;time&#34;
    &#34;encoding/json&#34;
    &#34;encoding/xml&#34;
    &#34;github.com/streadway/amqp&#34;
    // &#34;github.com/sethgrid/pester&#34;
)

var exit = make(chan bool, 1)
var channels = make([]string, 0)
var client = &amp;http.Client{}

type ACResponse struct {
    Name string `json:&#34;name&#34;`
    Parent string `json:&#34;parent&#34;`
}

type Channel struct {
    XMLName xml.Name `xml:&#34;Channel&#34;`
    ChannelId string `xml:&#34;channelId&#34;`
}

type Subscription struct {
    XMLName xml.Name `xml:&#34;Subscription&#34;`
    SubscriptionId string `xml:&#34;subscriptionId&#34;`
}

type Event struct {
    XMLName xml.Name `xml:&#34;xsi:Event&#34;`
    EventId string `xml:&#34;xsi:eventID&#34;`
}


func req(url string, data []byte) (c *http.Response){
    req, err := http.NewRequest(&#34;POST&#34;, url, bytes.NewBuffer(data))
    req.Header.Set(&#34;Authorization&#34;, fmt.Sprintf(&#34;Basic %s&#34;, os.Getenv(&#34;BSFT_TOKEN&#34;)))
    req.Header.Add(&#34;Content-Type&#34;, &#34;text/plain;charset=UTF-8&#34;)
    req.Header.Add(&#34;Accept-Encoding&#34;, &#34;identity&#34;)
    req.Header.Del(&#34;User-Agent&#34;)
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    return resp
}

func post(url string, data []byte) (c []byte){
    req, err := http.NewRequest(&#34;POST&#34;, url, bytes.NewBuffer(data))
    req.Header.Set(&#34;Authorization&#34;, fmt.Sprintf(&#34;Basic %s&#34;, os.Getenv(&#34;BSFT_TOKEN&#34;)))
    req.Header.Add(&#34;Content-Type&#34;, &#34;text/plain;charset=UTF-8&#34;)
    req.Header.Add(&#34;Accept-Encoding&#34;, &#34;identity&#34;)
    req.Header.Del(&#34;User-Agent&#34;)
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
    req, _ := http.NewRequest(&#34;PUT&#34;, fmt.Sprintf(&#34;http://URL/channel/%s/heartbeat&#34;, channelId), nil)
    req.Header.Add(&#34;Authorization&#34;, fmt.Sprintf(&#34;Basic %s&#34;, os.Getenv(&#34;BSFT_TOKEN&#34;)))
    req.Header.Add(&#34;Content-Type&#34;, &#34;text/plain;charset=UTF-8&#34;)
    req.Header.Add(&#34;Accept-Encoding&#34;, &#34;identity&#34;)
    req.Header.Del(&#34;User-Agent&#34;)
    resp, _ := client.Do(req)
    body, _ := ioutil.ReadAll(resp.Body)
    if (len(body) != 0) {
        log.Println(&#34;error on heartbeat for&#34;, channelId, string(body))
    }
    resp.Body.Close()
}

func channelTicker(channelId string) {
    ticker := time.NewTicker(time.Second * 15)
    fmt.Printf(&#34;%s Starting Channel heartbeat for %s\n&#34;, time.Now(), channelId)
    for t := range ticker.C {
        fmt.Printf(&#34;Heartbeat at %s for %s\n&#34;, t, channelId)
        channelHeartbeat(channelId)
    }
}


func subscribe(ch *amqp.Channel, q amqp.Queue, group ACResponse) (scanner *bufio.Scanner, err error){
    // Stream to read chunk by chunk.  
    res := req(&#34;http://URL/channel&#34;, channelTemplate(group.Name))
    reader := bufio.NewScanner(res.Body)
    // Make subscriptions
    fmt.Println(time.Now(), &#34;Subscribed to&#34;, group.Name)
    for reader.Scan() {
        text := reader.Text()
        c := Channel{ChannelId: &#34;none&#34;}
        merr := xml.Unmarshal([]byte(text), &amp;c)
        err = ch.Publish(
            &#34;events&#34;,     // exchange
            &#34;events&#34;, // routing key
            false,  // mandatory
            false,  // immediate
            amqp.Publishing{
                ContentType: &#34;text/plain&#34;,
                Body:        []byte(text),
            })
        // log.Printf(&#34; [x] Sent %s&#34;, text)
        if merr == nil {
            fmt.Println(&#34;ChannelId&#34;, c.ChannelId)
            go channelTicker(c.ChannelId)
        }
        failOnError(err, &#34;Failed to publish a message&#34;)
    }
    if err := reader.Err(); err != nil {
        panic(err)
        exit &lt;- true
    }
    return reader, err
}


func failOnError(err error, msg string) {
  if err != nil {
    log.Fatalf(&#34;%s: %s&#34;, msg, err)
    panic(fmt.Sprintf(&#34;%s: %s&#34;, msg, err))
  }
}

func main() {
    request, _ := http.NewRequest(&#34;GET&#34;, &#34;URL/subscriptions/?format=json&#34;, nil)
    request.Header.Add(&#34;Content-Type&#34;, &#34;application/json&#34;)
    resp, _ := client.Do(request)
    var ac []ACResponse
    json.NewDecoder(resp.Body).Decode(&amp;ac)
    resp.Body.Close()


    conn, err := amqp.Dial(&#34;amqp://URL:5672//&#34;)
    failOnError(err, &#34;Failed to connect to RabbitMQ&#34;)
    defer conn.Close()

    ch, err := conn.Channel()
    failOnError(err, &#34;Failed to open a channel&#34;)
    defer ch.Close()

    q, err := ch.QueueDeclare(
        &#34;events&#34;, // name
        false,   // durable
        false,   // delete when unused
        false,   // exclusive
        false,   // no-wait
        nil,     // arguments
    )
    failOnError(err, &#34;Failed to declare a queue&#34;)
    for _, item := range ac[:1] {
        go subscribe(ch, q, item)
        // time.Sleep(time.Second * 7)
    }

    &lt;-exit

}