package repositories

import (
	"5g-v2x-user-service/internal/config"
	"5g-v2x-user-service/internal/infrastructures/database"
	"5g-v2x-user-service/internal/models"
	"context"
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

func (dr *DriverRepository) FindOne(driverID, username *string) (*models.Driver, error) {
	collection := dr.MONGO.Client.Database(dr.config.DatabaseName).Collection("driver")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result *models.Driver

	filterDeleted := bson.M{
		"deleted_at": bson.M{
			"$eq": nil,
		},
	}

	inputFilter := bson.M{}

	if driverID != nil {
		inputFilter = bson.M{
			"_id": *driverID,
		}
	} else if username != nil {
		inputFilter = bson.M{
			"username": *username,
		}
	}

	filter := bson.M{
		"$and": []bson.M{
			filterDeleted,
			inputFilter,
		},
	}

	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (dr *DriverRepository) Find() ([]*models.Driver, error) {
	collection := dr.MONGO.Client.Database(dr.config.DatabaseName).Collection("driver")

	var results []*models.Driver

	filter := bson.M{
		"deleted_at": bson.M{
			"$eq": nil,
		},
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

func (dr *DriverRepository) Delete(driverID string) error {
	collection := dr.MONGO.Client.Database(dr.config.DatabaseName).Collection("driver")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	bsonFilter := bson.M{"_id": driverID}
	bsonUpdate := bson.D{
		{"$set", bson.D{{"deleted_at", time.Now().UTC()}}},
	}

	_, err := collection.UpdateOne(ctx, bsonFilter, bsonUpdate)

	if err != nil {
		return err
	}
	return nil
}
