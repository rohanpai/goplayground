package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func getredis() redis.Conn {
	c, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		panic("")
	}
	/*if _, err := c.Do("AUTH", "sedsed"); err != nil {
		c.Close()
		panic("")
	}*/
	return c
}

func Bench(b *testing.B) {
	b.StartTimer()
	for x := 1; x < b.N; x++ {
		createone()
	}
}

func BenchLua(b *testing.B) {
	b.StartTimer()
	for x := 1; x < b.N; x++ {
		createonelua()
	}
}

var LuaInsertScript string = `local id = redis.call('INCR', 'posthighestid')
local tags = {}
local tagstring = ARGV[1]
local hashstring = ARGV[2]
local extstring = ARGV[3]
tagstring:gsub('([^,]*),', function(c) table.insert(tags, c) end)

for nothing, t in pairs(tags) do 
	redis.call('ZADD', 'tag:' .. t, id, 'post:' .. id)
	redis.call('INCR', 'tagcount:' .. t)
end

redis.call('HMSET', 'post:'..id, 'id', id, 'hash', hashstring, 'ext', extstring, 'tags', tagstring)
redis.call('ZADD', 'posts', id, 'post:'..id)

return id
`

var tags []string
var r redis.Conn

func main() {

	r = getredis()

	tagfile, err := os.Open("tags.txt")
	defer tagfile.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	rand.Seed(time.Now().UnixNano())

	rdr := bufio.NewReader(tagfile)
	line, _, _ := rdr.ReadLine()
	tags = strings.Split(string(line), ",")

	fmt.Println(testing.Benchmark(Bench))
	fmt.Println(testing.Benchmark(BenchLua))

	return

	topause := 0
	for x := 1; x < 1000000; x++ {
		topause++
		createonelua()
		if topause%1000 == 0 {
			fmt.Println(topause)

			//time.Sleep(time.Second * 0)
		}
	}
}

func createonelua() {
	ptags := ""

	for x := 1; x < (rand.Intn(25) + 5); x++ {
		randtag := tags[rand.Intn(len(tags))]
		ptags = ptags + randtag + ","
	}

	hash := md5.New()
	hash.Write([]byte{byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255))})
	hash.Write([]byte{byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255))})

	_, err := r.Do("EVAL", LuaInsertScript, 0, ptags, fmt.Sprintf("%x", hash.Sum(nil)), ".jpg")

	if err != nil {
		fmt.Println("err: ", err.Error())
	}
}

func createone() {

	id, _ := redis.Int(r.Do("INCR", "posthighestid"))

	r.Send("MULTI")

	ptags := ""

	for x := 1; x < (rand.Intn(25) + 5); x++ {
		randtag := tags[rand.Intn(len(tags))]
		ptags = ptags + randtag + ","

		r.Send("ZADD", "tag:"+randtag, id, "post:"+strconv.Itoa(id))
		r.Send("INCR", "tagcount:"+randtag)
	}

	hash := md5.New()
	hash.Write([]byte{byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255))})
	hash.Write([]byte{byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255))})

	r.Send("HMSET", "post:"+strconv.Itoa(id), "id", id, "hash", fmt.Sprintf("%x", hash.Sum(nil)), "ext", ".jpg", "tags", ptags)
	r.Send("ZADD", "posts", id, "post:"+strconv.Itoa(id))
	r.Send("EXEC")
	r.Flush()
}
