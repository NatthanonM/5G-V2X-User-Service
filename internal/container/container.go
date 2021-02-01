package container

import (
	"5g-v2x-user-service/internal/config"
	controller "5g-v2x-user-service/internal/controllers"
	"5g-v2x-user-service/internal/infrastructures/database"
	"5g-v2x-user-service/internal/infrastructures/http"
	"5g-v2x-user-service/internal/repositories"
	"5g-v2x-user-service/internal/services"

	"go.uber.org/dig"
)

type Container struct {
	container *dig.Container
	Error     error
}

func NewContainer() *Container {
	c := new(Container)
	c.Configure()
	return c
}

func (cn *Container) Configure() {
	cn.container = dig.New()

	// infrastructures
	if err := cn.container.Provide(http.NewGRPCServer); err != nil {
		cn.Error = err
	}

	if err := cn.container.Provide(database.NewMongoDatabase); err != nil {
		cn.Error = err
	}

	// config
	if err := cn.container.Provide(config.NewConfig); err != nil {
		cn.Error = err
	}

	// controllers
	if err := cn.container.Provide(controller.NewController); err != nil {
		cn.Error = err
	}

	if err := cn.container.Provide(controller.NewAdminController); err != nil {
		cn.Error = err
	}

	if err := cn.container.Provide(controller.NewDriverController); err != nil {
		cn.Error = err
	}

	// services
	if err := cn.container.Provide(services.NewAdminService); err != nil {
		cn.Error = err
	}

	if err := cn.container.Provide(services.NewDriverService); err != nil {
		cn.Error = err
	}

	// repositories
	if err := cn.container.Provide(repositories.NewAdminRepository); err != nil {
		cn.Error = err
	}

	if err := cn.container.Provide(repositories.NewDriverRepository); err != nil {
		cn.Error = err
	}

}

func (cn *Container) Run() *Container {
	if err := cn.container.Invoke(func(g *http.GRPCServer) {
		if err := g.Start(); err != nil {
			panic(err)
		}
	}); err != nil {
		panic(err)
	}
	return cn
}
