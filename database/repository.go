package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/feastM/HatParty/config"
	"github.com/feastM/HatParty/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(config.Cfg.MongoConnectionString))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
	return client
}

var DB *mongo.Client = ConnectDB()

func InsertHats(client *mongo.Client, collectionName string, hats []models.Hat) {
	var docs []interface{}

	for _, i := range hats {
		docs = append(docs, i)
	}

	result, err := client.Database("HatParty").Collection(collectionName).InsertMany(context.TODO(), docs)

	fmt.Println(result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
		}
		log.Fatal(err)
	}
}

func AddParty(client *mongo.Client, party models.Party) (errs error) {
	hatCollection := client.Database("HatParty").Collection("Hats")
	partyCollection := client.Database("HatParty").Collection("Parties")

	err := client.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}

		filter := bson.D{
			{"$or",
				bson.A{
					bson.D{{"cleaning", bson.D{{"$lt", time.Now().Local().Add(time.Minute * time.Duration(config.Cfg.CleaningTimeInHours))}}}},
					bson.D{{"usedby", nil}},
				}},
			{"$or",
				bson.A{
					bson.D{{"cleaning", nil}},
					bson.D{{"usedby", nil}},
				}}}
		ops := options.FindOneAndUpdate().SetSort(bson.D{{"priority", -1}, {"cleaning", 1}})

		var hat models.Hat

		for i := 0; i < party.HatsRequested; i++ {
			err := hatCollection.FindOneAndUpdate(sessionContext, filter, bson.D{{"$set", bson.D{{"usedby", party.Id}}}}, ops).Decode(&hat)
			if err != nil {
				return err
			}
			hat.UsedBy = &party.Id
			party.Hats = append(party.Hats, hat)
		}

		_, err = partyCollection.InsertOne(sessionContext, party)

		if err != nil {
			sessionContext.AbortTransaction(sessionContext)
			return err
		}

		if err = sessionContext.CommitTransaction(sessionContext); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return
}

func StopParty(client *mongo.Client, partyId string) (errs error) {
	hatCollection := client.Database("HatParty").Collection("Hats")
	partyCollection := client.Database("HatParty").Collection("Parties")

	err := client.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}

		partyCollection.FindOneAndUpdate(sessionContext, bson.D{{"id", partyId}}, bson.D{{"$set", bson.D{{"status", 0}}}})
		if err != nil {
			return err
		}

		_, err = hatCollection.UpdateMany(sessionContext, bson.D{{"usedby", partyId}}, bson.D{{"$set", bson.D{{"usedby", nil}, {"cleaning", time.Now()}}}})

		if err != nil {
			sessionContext.AbortTransaction(sessionContext)
			return err
		}

		if err = sessionContext.CommitTransaction(sessionContext); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return
}
