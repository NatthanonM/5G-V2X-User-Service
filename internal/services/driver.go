package services

import (
	"5g-v2x-user-service/internal/config"
	"5g-v2x-user-service/internal/models"
	"5g-v2x-user-service/internal/repositories"
	"5g-v2x-user-service/internal/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DriverService struct {
	*repositories.DriverRepository
	*config.Config
}

func NewDriverService(DriverRepository *repositories.DriverRepository, Config *config.Config) *DriverService {
	return &DriverService{
		DriverRepository: DriverRepository,
		Config:           Config,
	}
}

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
