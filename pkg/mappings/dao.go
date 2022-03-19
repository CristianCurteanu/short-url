package mappings

import (
	"context"
	"encoding/base64"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewUrlMappingDao(db *mongo.Collection) URLMappingDAO {
	db.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "key", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	return &urlMappingDAO{db}
}

func CreateKey(url string) string {
	baseUrl := base64.StdEncoding.EncodeToString([]byte(url))
	len := len(baseUrl)
	return baseUrl[len-9 : len-1]
}

type URLMappingDAO interface {
	Add(ctx context.Context, mapping UrlMapping) error
	SearchById(ctx context.Context, key string) (UrlMapping, error)
	Delete(ctx context.Context, key string) error
	IncrementCounter(ctx context.Context, key string) error
}

type urlMappingDAO struct {
	db *mongo.Collection
}

func (ud *urlMappingDAO) Add(ctx context.Context, mapping UrlMapping) error {
	_, err := ud.db.InsertOne(ctx, mapping)

	return err
}

func (ud *urlMappingDAO) SearchById(ctx context.Context, key string) (UrlMapping, error) {
	var result UrlMapping
	err := ud.db.FindOne(ctx, bson.M{"key": key}).Decode(&result)
	if err != nil {
		return UrlMapping{}, err
	}

	return result, nil
}

func (ud *urlMappingDAO) Delete(ctx context.Context, key string) error {
	_, err := ud.db.DeleteOne(ctx, bson.M{"key": key})
	return err
}

func (ud *urlMappingDAO) IncrementCounter(ctx context.Context, key string) error {
	filter := bson.M{"key": key}
	update := bson.D{{"$inc", bson.D{{"counter", 1}}}}

	_, err := ud.db.UpdateOne(ctx, filter, update)

	return err
}
