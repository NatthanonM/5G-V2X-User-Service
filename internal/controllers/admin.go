package controllers

import (
	"5g-v2x-user-service/internal/config"
	"5g-v2x-user-service/internal/services"
	proto "5g-v2x-user-service/pkg/api"
	"context"

	"github.com/golang/protobuf/ptypes/empty"
)

type AdminController struct {
	*services.AdminService
	*config.Config
}

func NewAdminController(AdminService *services.AdminService, Config *config.Config) *AdminController {
	return &AdminController{
		AdminService: AdminService,
		Config:       Config,
	}
}

func (ac *AdminController) VerifyAdminAccessToken(ctx context.Context, req *proto.VerifyAdminAccessTokenRequest) (*proto.VerifyAdminAccessTokenResponse, error) {
	return &proto.VerifyAdminAccessTokenResponse{
		AccessToken: "",
	}, nil
}

func (ac *AdminController) RegisterAdmin(ctx context.Context, req *proto.RegisterAdminRequest) (*empty.Empty, error) {
	if err := ac.AdminService.Register(req.Username, req.Password); err != nil {
		return nil, err
	}
	return new(empty.Empty), nil
}

func (ac *AdminController) LoginAdmin(ctx context.Context, req *proto.LoginAdminRequest) (*proto.LoginAdminResponse, error) {
	return &proto.LoginAdminResponse{
		AccessToken: "",
	}, nil
}
