package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/innoglobe/xmgo/internal/entity"
)

type CompanyRepositoryInterface interface {
	Create(ctx context.Context, company *entity.Company) (*entity.Company, error)
	Update(ctx context.Context, id uuid.UUID, company *entity.Company) (*entity.Company, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Get(ctx context.Context, id uuid.UUID) (*entity.Company, error)
}
