package controllers

import (
	"5g-v2x-user-service/internal/config"
	"5g-v2x-user-service/internal/models"
	"5g-v2x-user-service/internal/services"
	"5g-v2x-user-service/internal/utils"
	proto "5g-v2x-user-service/pkg/api"
	"context"
	"fmt"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (ds *DriverController) GetAllDriver(ctx context.Context, req *empty.Empty) (*proto.GetAllDriverResponse, error) {
	drivers, err := ds.DriverService.GetAllDriver()
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal error.")
	}
	var resDrivers []*proto.Driver
	for _, driver := range drivers {
		resDrivers = append(resDrivers, &proto.Driver{
			DriverId:    driver.DriverID,
			Firstname:   driver.Firstname,
			Lastname:    driver.Lastname,
			DateOfBirth: utils.WrapperTime(&driver.DateOfBirth),
			Gender:      driver.Gender,
			Username:    driver.Username,
		})
	}
	return &proto.GetAllDriverResponse{
		Drivers: resDrivers,
	}, nil
}

func (ds *DriverController) GetDriver(ctx context.Context, req *proto.GetDriverRequest) (*proto.Driver, error) {
	driver, err := ds.DriverService.GetDriver(req.DriverId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &proto.Driver{
		DriverId:    driver.DriverID,
		Firstname:   driver.Firstname,
		Lastname:    driver.Lastname,
		DateOfBirth: utils.WrapperTime(&driver.DateOfBirth),
		Gender:      driver.Gender,
		Username:    driver.Username,
	}, nil
}

func (ds *DriverController) GetDriverByUsername(ctx context.Context, req *proto.GetDriverByUsernameRequest) (*proto.Driver, error) {
	driver, err := ds.DriverService.GetDriverByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	fmt.Println(driver)
	return &proto.Driver{
		DriverId:    driver.DriverID,
		Firstname:   driver.Firstname,
		Lastname:    driver.Lastname,
		DateOfBirth: utils.WrapperTime(&driver.DateOfBirth),
		Gender:      driver.Gender,
		Username:    driver.Username,
	}, nil
}

func (ds *DriverController) LoginDriver(ctx context.Context, req *proto.LoginDriverRequest) (*proto.LoginDriverResponse, error) {
	driver, err := ds.DriverService.CheckEmailPassword(req.Username, req.Password)

	if err != nil {
		return nil, err
	}

	return &proto.LoginDriverResponse{
		DriverId: driver.DriverID,
	}, nil
}

func (ds *DriverController) UpdateDriver(ctx context.Context, req *proto.UpdateDriverRequest) (*proto.UpdateDriverResponse, error) {
	err := ds.DriverService.UpdateDriver(&models.Driver{
		DriverID:    req.DriverId,
		Firstname:   req.Firstname,
		Lastname:    req.Lastname,
		DateOfBirth: req.DateOfBirth.AsTime(),
	})

	if err != nil {
		return nil, err
	}

	return &proto.UpdateDriverResponse{}, nil
}
