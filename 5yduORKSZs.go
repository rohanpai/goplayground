package main

import (
	"log"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

func main() {

	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %s", err)
	}

	db := session.DB("TestDatabase")

	maxSize := int(4 * 1024 * 1024 * 1024) // 4294967296 bytes

	collectionName := "TestCollection"
	collection := db.C(collectionName)

	collectionInfo := &mgo.CollectionInfo{
		Capped:   true,
		MaxBytes: maxSize,
	}

	err = collection.Create(collectionInfo)
	if err != nil {
		log.Fatalf("Failed to create the collection: %s", err)
	}

	result := &bson.D{}
	err = db.Run(&bson.D{bson.DocElem{"collstats", collectionName}}, result)
	if err != nil {
		log.Fatalf("Failed to get collection stats: %s", err)
	}

	storageSize := result.Map()["storageSize"]

	log.Printf("Maximum Size was: %d", maxSize)
	log.Printf("Storage Size was: %d", storageSize)

}
