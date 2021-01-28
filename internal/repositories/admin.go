package repositories

import (
	"5g-v2x-user-service/internal/config"
	"5g-v2x-user-service/internal/infrastructures/database"
	"5g-v2x-user-service/internal/models"
	"context"
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type AdminRepository struct {
	MONGO  *database.MongoDatabase
	config *config.Config
}

func NewAdminRepository(m *database.MongoDatabase, c *config.Config) *AdminRepository {
	return &AdminRepository{
		MONGO:  m,
		config: c,
	}
}

func (ar *AdminRepository) Create(m interface{}) error {
	collection := ar.MONGO.Client.Database(ar.config.DatabaseName).Collection("admin")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, m)
	if err != nil {
		return err
	}
	return nil
}

func (ar *AdminRepository) FindOne(filter map[string]interface{}) (*models.Admin, error) {
	collection := ar.MONGO.Client.Database(ar.config.DatabaseName).Collection("admin")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result *models.Admin

	jsonString, err := json.Marshal(filter)
	if err != nil {
		panic(err)
	}

	var bsonFilter interface{}
	err = bson.UnmarshalExtJSON([]byte(jsonString), true, &bsonFilter)
	if err != nil {
		panic(err)
	}

	err = collection.FindOne(ctx, bsonFilter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
