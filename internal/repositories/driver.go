package repositories

import (
	"5g-v2x-user-service/internal/config"
	"5g-v2x-user-service/internal/infrastructures/database"
	"5g-v2x-user-service/internal/models"
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type DriverRepository struct {
	MONGO  *database.MongoDatabase
	config *config.Config
}

func NewDriverRepository(m *database.MongoDatabase, c *config.Config) *DriverRepository {
	return &DriverRepository{
		MONGO:  m,
		config: c,
	}
}

func (dr *DriverRepository) Create(driver *models.Driver) (string, error) {
	collection := dr.MONGO.Client.Database(dr.config.DatabaseName).Collection("driver")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id := uuid.New()
	driver.DriverID = id.String()
	_, err := collection.InsertOne(ctx, driver)
	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func (dr *DriverRepository) FindOne(filter map[string]interface{}) (*models.Driver, error) {
	collection := dr.MONGO.Client.Database(dr.config.DatabaseName).Collection("driver")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result *models.Driver

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

func (dr *DriverRepository) Find(filter map[string]interface{}) ([]*models.Driver, error) {
	collection := dr.MONGO.Client.Database(dr.config.DatabaseName).Collection("driver")

	var results []*models.Driver

	jsonString, err := json.Marshal(filter)
	if err != nil {
		panic(err)
	}

	var bsonFilter interface{}
	err = bson.UnmarshalExtJSON([]byte(jsonString), true, &bsonFilter)
	if err != nil {
		panic(err)
	}

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		var elem models.Driver
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(context.TODO())

	return results, nil
}

func (dr *DriverRepository) Update(updateDriver *models.Driver) error {
	collection := dr.MONGO.Client.Database(dr.config.DatabaseName).Collection("driver")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	bsonFilter := bson.M{"_id": updateDriver.DriverID}
	bsonUpdate := bson.D{
		{"$set", bson.D{{"firstname", updateDriver.Firstname}, {"lastname", updateDriver.Lastname}, {"date_of_birth", updateDriver.DateOfBirth}}},
	}

	_, err := collection.UpdateOne(ctx, bsonFilter, bsonUpdate)

	if err != nil {
		return err
	}
	return nil
}
