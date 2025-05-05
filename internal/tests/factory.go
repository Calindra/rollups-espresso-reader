package tests

import (
	"context"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type PostgresContainer struct {
	container *postgres.PostgresContainer
	config    *PostgresContainerDBConfig
}

type PostgresContainerDBConfig struct {
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"-"`
}

const postgresImageVersion = "postgres:17.4-bookworm"

func NewPostgresContainer(name, user, password string) *PostgresContainer {
	return &PostgresContainer{
		config: &PostgresContainerDBConfig{
			Name:     name,
			User:     user,
			Password: password,
		},
	}
}

func (pc *PostgresContainer) RunPostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	// Create a new PostgreSQL container
	pgContainer, err := postgres.Run(ctx, postgresImageVersion)
	if err != nil {
		return nil, err
	}

	// Wait for the container to be ready
	if err := pgContainer.Start(ctx); err != nil {
		return nil, err
	}

	return &PostgresContainer{container: pgContainer, config: pc.config}, nil
}
