package utils

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"strings"
)

type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

var Mg MongoInstance

var Secret string

const dbName = "test"

func MongoDBConnect() (error, func()) {
	mongoURI := getMongoURLandPopulateSecretString()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	deferFunc := func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}
	if err != nil {
		fmt.Println(err)
		return err, deferFunc
	}

	db := client.Database(dbName)

	if err != nil {
		return err, deferFunc
	}

	Mg = MongoInstance{
		Client: client,
		Db:     db,
	}
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"email", 1}},
		Options: options.Index().SetUnique(true),
	}

	// Create the index
	_, err = Mg.Db.Collection("users").Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatal("Error creating the index")
	}

	fmt.Println("Connected to MongoDB cloud")
	return nil, deferFunc
}

func getMongoURLandPopulateSecretString() string {
	if os.Getenv("ENV") == "DEV" {
		Secret = os.Getenv("JWT_S")
		return os.Getenv("M_URI")
	}
	name := "projects/1037996227658/secrets/quizapp_s/versions/4"
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatal("failed to create secretmanager client: %w", err)

	}
	defer client.Close()
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		log.Fatal("failed to access secret version: %w", err)
	}
	stringVal := string(result.Payload.Data)
	words := strings.Fields(stringVal)
	Secret = words[1]
	return words[0]

}
