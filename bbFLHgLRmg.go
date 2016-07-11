package main

import (
	&#34;bufio&#34;
	&#34;crypto/md5&#34;
	&#34;fmt&#34;
	&#34;github.com/garyburd/redigo/redis&#34;
	&#34;math/rand&#34;
	&#34;os&#34;
	&#34;strconv&#34;
	&#34;strings&#34;
	&#34;testing&#34;
	&#34;time&#34;
)

func getredis() redis.Conn {
	c, err := redis.Dial(&#34;tcp&#34;, &#34;localhost:6379&#34;)
	if err != nil {
		panic(&#34;&#34;)
	}
	/*if _, err := c.Do(&#34;AUTH&#34;, &#34;sedsed&#34;); err != nil {
		c.Close()
		panic(&#34;&#34;)
	}*/
	return c
}

func Bench(b *testing.B) {
	b.StartTimer()
	for x := 1; x &lt; b.N; x&#43;&#43; {
		createone()
	}
}

func BenchLua(b *testing.B) {
	b.StartTimer()
	for x := 1; x &lt; b.N; x&#43;&#43; {
		createonelua()
	}
}

var LuaInsertScript string = `local id = redis.call(&#39;INCR&#39;, &#39;posthighestid&#39;)
local tags = {}
local tagstring = ARGV[1]
local hashstring = ARGV[2]
local extstring = ARGV[3]
tagstring:gsub(&#39;([^,]*),&#39;, function(c) table.insert(tags, c) end)

for nothing, t in pairs(tags) do 
	redis.call(&#39;ZADD&#39;, &#39;tag:&#39; .. t, id, &#39;post:&#39; .. id)
	redis.call(&#39;INCR&#39;, &#39;tagcount:&#39; .. t)
end

redis.call(&#39;HMSET&#39;, &#39;post:&#39;..id, &#39;id&#39;, id, &#39;hash&#39;, hashstring, &#39;ext&#39;, extstring, &#39;tags&#39;, tagstring)
redis.call(&#39;ZADD&#39;, &#39;posts&#39;, id, &#39;post:&#39;..id)

return id
`

var tags []string
var r redis.Conn

func main() {

	r = getredis()

	tagfile, err := os.Open(&#34;tags.txt&#34;)
	defer tagfile.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	rand.Seed(time.Now().UnixNano())

	rdr := bufio.NewReader(tagfile)
	line, _, _ := rdr.ReadLine()
	tags = strings.Split(string(line), &#34;,&#34;)

	fmt.Println(testing.Benchmark(Bench))
	fmt.Println(testing.Benchmark(BenchLua))

	return

	topause := 0
	for x := 1; x &lt; 1000000; x&#43;&#43; {
		topause&#43;&#43;
		createonelua()
		if topause%1000 == 0 {
			fmt.Println(topause)

			//time.Sleep(time.Second * 0)
		}
	}
}

func createonelua() {
	ptags := &#34;&#34;

	for x := 1; x &lt; (rand.Intn(25) &#43; 5); x&#43;&#43; {
		randtag := tags[rand.Intn(len(tags))]
		ptags = ptags &#43; randtag &#43; &#34;,&#34;
	}

	hash := md5.New()
	hash.Write([]byte{byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255))})
	hash.Write([]byte{byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255))})

	_, err := r.Do(&#34;EVAL&#34;, LuaInsertScript, 0, ptags, fmt.Sprintf(&#34;%x&#34;, hash.Sum(nil)), &#34;.jpg&#34;)

	if err != nil {
		fmt.Println(&#34;err: &#34;, err.Error())
	}
}

func createone() {

	id, _ := redis.Int(r.Do(&#34;INCR&#34;, &#34;posthighestid&#34;))

	r.Send(&#34;MULTI&#34;)

	ptags := &#34;&#34;

	for x := 1; x &lt; (rand.Intn(25) &#43; 5); x&#43;&#43; {
		randtag := tags[rand.Intn(len(tags))]
		ptags = ptags &#43; randtag &#43; &#34;,&#34;

		r.Send(&#34;ZADD&#34;, &#34;tag:&#34;&#43;randtag, id, &#34;post:&#34;&#43;strconv.Itoa(id))
		r.Send(&#34;INCR&#34;, &#34;tagcount:&#34;&#43;randtag)
	}

	hash := md5.New()
	hash.Write([]byte{byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255))})
	hash.Write([]byte{byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255))})

	r.Send(&#34;HMSET&#34;, &#34;post:&#34;&#43;strconv.Itoa(id), &#34;id&#34;, id, &#34;hash&#34;, fmt.Sprintf(&#34;%x&#34;, hash.Sum(nil)), &#34;ext&#34;, &#34;.jpg&#34;, &#34;tags&#34;, ptags)
	r.Send(&#34;ZADD&#34;, &#34;posts&#34;, id, &#34;post:&#34;&#43;strconv.Itoa(id))
	r.Send(&#34;EXEC&#34;)
	r.Flush()
}
