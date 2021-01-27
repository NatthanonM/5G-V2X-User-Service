package services

import (
	"5g-v2x-user-service/internal/config"
	"5g-v2x-user-service/internal/repositories"
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
