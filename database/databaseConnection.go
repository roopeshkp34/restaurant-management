package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBInstance() *mongo.Client {
	MongoDb := "mongodb://localhost:28017"

	log.Println(MongoDb)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(MongoDb))
	if err != nil {
		log.Fatal(err)
	}
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	// err = client.Connect(ctx)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	log.Println("Connected to mongoDb")

	return client
}

var Client *mongo.Client = DBInstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("restaurant").Collection(collectionName)

	return collection

}
