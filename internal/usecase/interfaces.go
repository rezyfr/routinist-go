package usecase

import (
	"context"

	"routinist/internal/entity"
)

type (
	Auth interface {
		Register(ctx context.Context, e entity.RegisterRequestDTO) (entity.AuthResponseDTO, error)
		Login(ctx context.Context, e entity.LoginRequestDTO) (entity.AuthResponseDTO, error)
	}

	AuthRepo interface {
		Register(ctx context.Context, e entity.RegisterRequestDTO) (entity.AuthResponseDTO, error)
		Login(ctx context.Context, e entity.LoginRequestDTO) (entity.AuthResponseDTO, error)
	}
)
