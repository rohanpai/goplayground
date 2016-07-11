package main

import (
	&#34;log&#34;

	&#34;labix.org/v2/mgo&#34;
	&#34;labix.org/v2/mgo/bson&#34;
)

func main() {

	session, err := mgo.Dial(&#34;localhost&#34;)
	if err != nil {
		log.Fatalf(&#34;Error connecting to MongoDB: %s&#34;, err)
	}

	db := session.DB(&#34;TestDatabase&#34;)

	maxSize := int(4 * 1024 * 1024 * 1024) // 4294967296 bytes

	collectionName := &#34;TestCollection&#34;
	collection := db.C(collectionName)

	collectionInfo := &amp;mgo.CollectionInfo{
		Capped:   true,
		MaxBytes: maxSize,
	}

	err = collection.Create(collectionInfo)
	if err != nil {
		log.Fatalf(&#34;Failed to create the collection: %s&#34;, err)
	}

	result := &amp;bson.D{}
	err = db.Run(&amp;bson.D{bson.DocElem{&#34;collstats&#34;, collectionName}}, result)
	if err != nil {
		log.Fatalf(&#34;Failed to get collection stats: %s&#34;, err)
	}

	storageSize := result.Map()[&#34;storageSize&#34;]

	log.Printf(&#34;Maximum Size was: %d&#34;, maxSize)
	log.Printf(&#34;Storage Size was: %d&#34;, storageSize)

}
