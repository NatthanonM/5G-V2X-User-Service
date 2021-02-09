package services

import (
	"5g-v2x-user-service/internal/config"
	"5g-v2x-user-service/internal/models"
	"5g-v2x-user-service/internal/repositories"
	"5g-v2x-user-service/internal/utils"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DriverService ...
type DriverService struct {
	*repositories.DriverRepository
	*config.Config
}

// NewDriverService ...
func NewDriverService(DriverRepository *repositories.DriverRepository, Config *config.Config) *DriverService {
	return &DriverService{
		DriverRepository: DriverRepository,
		Config:           Config,
	}
}

// AddNewDriver ...
func (ds *DriverService) AddNewDriver(driver *models.Driver) (*string, error) {
	hashed, err := utils.HashAndSalt([]byte(driver.Password))
	if err != nil {
		return nil, err
	}
	driver.HashedPassword = hashed

	filter := make(map[string]interface{})
	filter["username"] = driver.Username

	if _, err := ds.DriverRepository.FindOne(filter); err == nil {
		return nil, status.Error(codes.AlreadyExists, "Username is already existed")
	}

	driverID, err := ds.DriverRepository.Create(driver)
	if err != nil {
		return nil, err
	}

	return &driverID, nil
}

// GetAllDriver ...
func (ds *DriverService) GetAllDriver() ([]*models.Driver, error) {
	filter := make(map[string]interface{})
	drivers, err := ds.DriverRepository.Find(filter)
	if err != nil {
		return nil, err
	}
	return drivers, nil
}

// GetDriver ...
func (ds *DriverService) GetDriver(driverID string) (*models.Driver, error) {
	filter := make(map[string]interface{})
	filter["_id"] = driverID
	driver, err := ds.DriverRepository.FindOne(filter)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Driver not found.")
	}
	return driver, nil
}

func (ds *DriverService) GetDriverByUsername(username string) (*models.Driver, error) {
	filter := make(map[string]interface{})
	filter["username"] = username

	driver, err := ds.DriverRepository.FindOne(filter)
	if err != nil {
		return nil, err
	}

	return driver, nil
}
