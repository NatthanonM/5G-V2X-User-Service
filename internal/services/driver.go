package services

import (
	"5g-v2x-user-service/internal/config"
	"5g-v2x-user-service/internal/models"
	"5g-v2x-user-service/internal/repositories"
	"5g-v2x-user-service/internal/utils"

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

	if _, err := ds.DriverRepository.FindOne(nil, &driver.Username); err == nil {
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
	drivers, err := ds.DriverRepository.Find()
	if err != nil {
		return nil, err
	}
	return drivers, nil
}

// GetDriver ...
func (ds *DriverService) GetDriver(driverID string) (*models.Driver, error) {
	driver, err := ds.DriverRepository.FindOne(&driverID, nil)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Driver not found.")
	}
	return driver, nil
}

func (ds *DriverService) GetDriverByUsername(username string) (*models.Driver, error) {
	driver, err := ds.DriverRepository.FindOne(nil, &username)
	if err != nil {
		return nil, err
	}

	return driver, nil
}

// CheckEmailPassword ...
func (ds *DriverService) CheckEmailPassword(username, password string) (*models.Driver, error) {
	// find user
	driver, err := ds.DriverRepository.FindOne(nil, &username)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Username or password is incorrect.")
	}

	// verify password
	if !driver.VerifyPassword(password) {
		return nil, status.Error(codes.Unauthenticated, "Username or password is incorrect.")
	}

	return driver, nil
}

func (ds *DriverService) UpdateDriver(updateDriver *models.Driver) error {
	_, err := ds.GetDriver(updateDriver.DriverID)

	if err != nil {
		return status.Error(codes.NotFound, "Driver not found")
	}

	err = ds.DriverRepository.Update(updateDriver)
	if err != nil {
		return status.Error(codes.Internal, "Update driver failed")
	}

	return err
}

func (ds *DriverService) DeleteDriver(driverID string) error {
	_, err := ds.GetDriver(driverID)

	if err != nil {
		return status.Error(codes.NotFound, "Driver not found")
	}

	err = ds.DriverRepository.Delete(driverID)
	if err != nil {
		return status.Error(codes.Internal, "Delete driver failed")
	}

	return err
}
