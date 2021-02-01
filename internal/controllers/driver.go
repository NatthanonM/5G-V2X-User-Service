package controllers

import (
	"5g-v2x-user-service/internal/config"
	"5g-v2x-user-service/internal/models"
	"5g-v2x-user-service/internal/services"
	proto "5g-v2x-user-service/pkg/api"
	"context"
)

type DriverController struct {
	*services.DriverService
	*config.Config
}

func NewDriverController(DriverService *services.DriverService, Config *config.Config) *DriverController {
	return &DriverController{
		DriverService: DriverService,
		Config:        Config,
	}
}

func (ds *DriverController) AddNewDriver(ctx context.Context, req *proto.AddNewDriverRequest) (*proto.AddNewDriverReponse, error) {
	driver := models.Driver{
		Firstname:   req.Firstname,
		Lastname:    req.Lastname,
		Username:    req.Username,
		Password:    req.Password,
		DateOfBirth: req.DateOfBirth.AsTime(),
		Gender:      req.Gender.String(),
	}
	driverID, err := ds.DriverService.AddNewDriver(&driver)
	if err != nil {
		return nil, err
	}
	return &proto.AddNewDriverReponse{
		DriverId: *driverID,
	}, nil
}
