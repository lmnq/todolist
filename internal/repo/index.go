package repo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// создание индекса для уникальных полей.
func CreateUniqueIndexes(client *mongo.Client) error {
	collection := client.Database(mongoDBName).Collection(taskCollection)

	// создание модели индекса
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"title", 1}, {"activeAt", 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	return nil
}
