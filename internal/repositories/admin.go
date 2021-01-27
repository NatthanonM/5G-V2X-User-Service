package repositories

import (
	"5g-v2x-user-service/internal/config"
	"5g-v2x-user-service/internal/infrastructures/database"
)

type AdminRepository struct {
	MONGO  *database.MongoDatabase
	config *config.Config
}

func NewAdminRepository(m *database.MongoDatabase, c *config.Config) *AdminRepository {
	return &AdminRepository{
		MONGO:  m,
		config: c,
	}
}
