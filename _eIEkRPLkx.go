package main

import (
  &#34;encoding/json&#34;
  &#34;fmt&#34;
  &#34;github.com/hoisie/web&#34;
  &#34;labix.org/v2/mgo&#34;
  &#34;os&#34;
)

var (
  dbSession           *mgo.Session
  globalConfiguration *Configuration = new(Configuration)
)

// structs
type Reading struct {
  Id   string
  Name string
}

type Configuration struct {
  DatabaseServer string
  DatabasePort   string
  DatabaseName   string
  DatabaseUsername string
  DatabasePassword string
}

func main() {
  loadConfiguration()

  web.Get(&#34;/machines/?&#34;, index)
  web.Post(&#34;/machines/?&#34;, create)
  web.Run(&#34;0.0.0.0:3000&#34;)
}

func index(webContext *web.Context) {
  readings := getReadings()

  webContext.WriteString(readingsToString(readings))
}

func create(webContext *web.Context) {
  // Setup readings
  readings := prepareReadings()

  // Insert readings
  insertReadings(readings)

  // return readings
  webContext.WriteString(readingsToString(readings))
}

// private
func readingsToString(readings []Reading) string {
  data, err := json.Marshal(readings)
  if err != nil {
    panic(err)
  }

  return string(data)
}

func getSession() *mgo.Session {
  if dbSession == nil {
    var credentials string
    if globalConfiguration.DatabaseUsername != &#34;&#34; &amp;&amp; globalConfiguration.DatabasePassword != &#34;&#34; {
      credentials = globalConfiguration.DatabaseUsername &#43; &#34;:&#34; &#43; globalConfiguration.DatabasePassword &#43; &#34;@&#34;
    }

    url := &#34;mongodb://&#34; &#43; credentials &#43; globalConfiguration.DatabaseServer &#43; &#34;:&#34; &#43; globalConfiguration.DatabasePort
    fmt.Println(&#34;Connecting to: &#34; &#43; url)

    var err error
    dbSession, err = mgo.Dial(url)
    if err != nil {
      panic(err) // no, not really
    }
  }

  return dbSession.Clone()
}

func getCollection(session *mgo.Session, databaseName string, tableName string) *mgo.Collection {
  // Optional. Switch the session to a monotonic behavior.
  // database.Session.SetMode(mgo.Monotonic, true)

  collection := session.DB(globalConfiguration.DatabaseName).C(tableName)

  return collection
}

func insertReadings(readings []Reading) {
  // Setup session
  session := getSession()
  defer session.Close()

  // Setup collection
  collection := getCollection(session, &#34;test&#34;, &#34;readings&#34;)

  // insert the reading
  err := collection.Insert(readings)
  if err != nil {
    fmt.Println(&#34;error insertReadings:&#34;, err)
    panic(err)
  }
}

func getReadings() []Reading {
  // Setup session
  session := getSession()
  defer session.Close()

  // Setup collection
  collection := getCollection(session, &#34;test&#34;, &#34;readings&#34;)

  readings := make([]Reading, 1)
  err := collection.Find(nil).All(&amp;readings)

  if err != nil {
    fmt.Println(&#34;error getReadings:&#34;, err)
    panic(err)
  }

  return readings
}

func prepareReadings() []Reading {
  var readings []Reading
  for i := 1; i &lt;= 1; i&#43;&#43; {
    readings = append(readings, Reading{Name: &#34;Thing&#34;})
  }

  return readings
}

func loadConfiguration() {
  configFileName := &#34;configs/mongodb.conf&#34;
  if len(os.Args) &gt; 1 {
    configFileName = os.Args[1]
  }

  configFile, err := os.Open(configFileName)
  if err != nil {
    fmt.Println(&#34;[ERROR] &#34; &#43; err.Error())
    fmt.Println(&#34;For your happiness an example config file is provided in the &#39;conf&#39; directory in the repository.&#34;)
    os.Exit(1)
  }

  configDecoder := json.NewDecoder(configFile)
  err = configDecoder.Decode(globalConfiguration)
  if err != nil {
    fmt.Println(&#34;[CONFIG FILE FORMAT ERROR] &#34; &#43; err.Error())
    fmt.Println(&#34;Please ensure that your config file is in valid JSON format.&#34;)
    os.Exit(1)
  }
}
