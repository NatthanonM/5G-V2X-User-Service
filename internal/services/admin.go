package services

import (
	"5g-v2x-user-service/internal/config"
	"5g-v2x-user-service/internal/models"
	"5g-v2x-user-service/internal/repositories"
	"5g-v2x-user-service/internal/utils"
	"strconv"
	"time"

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

// VerifyAccessToken ...
func (as *AdminService) VerifyAccessToken(StringAccessToken string) (*models.Admin, error) {
	// extract username, expires and signature
	accessToken, err := utils.ExtractAccessToken(StringAccessToken)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Cannot extract access token.")
	}

	// check expire time
	timestamp, err := strconv.ParseInt(accessToken.Expires, 10, 64)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Cannot convert expires.")
	}
	wasExpired := utils.WasExpired(time.Unix(timestamp, 0))
	if wasExpired == true {
		return nil, status.Error(codes.Unauthenticated, "Access token was expired.")
	}

	// check username exist
	filter := make(map[string]interface{})
	filter["username"] = accessToken.Username

	admin, err := as.AdminRepository.FindOne(filter)
	if err != nil {
		return nil, status.Error(codes.NotFound, "This username does not exist.")
	}

	if err := utils.VerifyAccessToken(accessToken, admin.HashedPassword); err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid access token.")
	}

	return admin, nil
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

// CheckEmailPassword ...
func (as *AdminService) CheckEmailPassword(username, password string) (*models.Admin, error) {
	// find user
	filter := make(map[string]interface{})
	filter["username"] = username

	admin, err := as.AdminRepository.FindOne(filter)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Username or password is incorrect.")
	}

	// verify password
	if !admin.VerifyPassword(password) {
		return nil, status.Error(codes.Unauthenticated, "Username or password is incorrect.")
	}

	return admin, nil
}

// GetAccessTokens ...
func (as *AdminService) GetAccessTokens(username, hashedPassword string) (string, error) {
	// generate HashBasedToken
	token, err := utils.GenerateAccessToken(username, hashedPassword)
	if err != nil {
		return "", status.Error(codes.Internal, "Generate access token failed.")
	}

	return token, nil
}
