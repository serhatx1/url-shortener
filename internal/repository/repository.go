package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type URLMapping struct {
	ShortURL    string    `bson:"short_url"`
	OriginalURL string    `bson:"original_url"`
	CreatedAt   time.Time `bson:"created_at"`
}

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		collection: db.Collection("url-shortener"),
	}
}

func (r *Repository) Save(originalURL string, shortURL string) error {
	urlMAP := URLMapping{
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}
	_, err := r.collection.InsertOne(context.Background(), urlMAP)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) IfLongUrlExist(originalURL string) (string, bool) {
	filter := bson.M{"original_url": originalURL}
	var result URLMapping
	err := r.collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return "", false
	}
	if err == mongo.ErrNoDocuments {
		return "", false
	}
	return result.ShortURL, true
}

func (r *Repository) FindOriginal(shortURL string) (string, bool) {
	filter := bson.M{"short_url": "http://localhost:8080" + shortURL}
	var result URLMapping
	err := r.collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return "", false
	}
	if err == mongo.ErrNoDocuments {
		return "", false
	}
	return result.OriginalURL, true
}
