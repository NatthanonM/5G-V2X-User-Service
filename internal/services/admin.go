package services

import (
	"5g-v2x-user-service/internal/config"
	"5g-v2x-user-service/internal/models"
	"5g-v2x-user-service/internal/repositories"
	"5g-v2x-user-service/internal/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AdminService struct {
	*repositories.AdminRepository
	*config.Config
}

func NewAdminService(AdminRepository *repositories.AdminRepository, Config *config.Config) *AdminService {
	return &AdminService{
		AdminRepository: AdminRepository,
		Config:          Config,
	}
}

func (as *AdminService) Register(username, password string) error {
	hashed, err := utils.HashAndSalt([]byte(password))
	if err != nil {
		return err
	}

	admin := models.Admin{
		Username:       username,
		HashedPassword: hashed,
	}

	filter := make(map[string]interface{})
	filter["username"] = username

	if _, err := as.AdminRepository.FindOne(filter); err == nil {
		return status.Error(codes.AlreadyExists, "Username is already existed")
	}

	if err := as.AdminRepository.Create(&admin); err != nil {
		return err
	}
	return nil
}
