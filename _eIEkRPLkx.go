package main

import (
  "encoding/json"
  "fmt"
  "github.com/hoisie/web"
  "labix.org/v2/mgo"
  "os"
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

  web.Get("/machines/?", index)
  web.Post("/machines/?", create)
  web.Run("0.0.0.0:3000")
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
    if globalConfiguration.DatabaseUsername != "" && globalConfiguration.DatabasePassword != "" {
      credentials = globalConfiguration.DatabaseUsername + ":" + globalConfiguration.DatabasePassword + "@"
    }

    url := "mongodb://" + credentials + globalConfiguration.DatabaseServer + ":" + globalConfiguration.DatabasePort
    fmt.Println("Connecting to: " + url)

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
  collection := getCollection(session, "test", "readings")

  // insert the reading
  err := collection.Insert(readings)
  if err != nil {
    fmt.Println("error insertReadings:", err)
    panic(err)
  }
}

func getReadings() []Reading {
  // Setup session
  session := getSession()
  defer session.Close()

  // Setup collection
  collection := getCollection(session, "test", "readings")

  readings := make([]Reading, 1)
  err := collection.Find(nil).All(&readings)

  if err != nil {
    fmt.Println("error getReadings:", err)
    panic(err)
  }

  return readings
}

func prepareReadings() []Reading {
  var readings []Reading
  for i := 1; i <= 1; i++ {
    readings = append(readings, Reading{Name: "Thing"})
  }

  return readings
}

func loadConfiguration() {
  configFileName := "configs/mongodb.conf"
  if len(os.Args) > 1 {
    configFileName = os.Args[1]
  }

  configFile, err := os.Open(configFileName)
  if err != nil {
    fmt.Println("[ERROR] " + err.Error())
    fmt.Println("For your happiness an example config file is provided in the 'conf' directory in the repository.")
    os.Exit(1)
  }

  configDecoder := json.NewDecoder(configFile)
  err = configDecoder.Decode(globalConfiguration)
  if err != nil {
    fmt.Println("[CONFIG FILE FORMAT ERROR] " + err.Error())
    fmt.Println("Please ensure that your config file is in valid JSON format.")
    os.Exit(1)
  }
}
