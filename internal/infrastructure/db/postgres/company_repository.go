package postgresrepository

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/innoglobe/xmgo/internal/customerrors"
	"github.com/innoglobe/xmgo/internal/entity"
	"github.com/innoglobe/xmgo/internal/interface/repository"
	"gorm.io/gorm"
	"strings"
)

// interface assertion to make sure it implements all methods
var _ repository.CompanyRepositoryInterface = &PostgresRepository{}

// PostgresRepository struct
type PostgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository creates a new instance of PostgresRepository
func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// Create insert company into the database
func (r *PostgresRepository) Create(ctx context.Context, company *entity.Company) (*entity.Company, error) {
	if err := r.db.WithContext(ctx).Create(company).Error; err != nil {
		fmt.Printf("ERROR: %T\n", err)
		// We use the sqlstate search because gorm.ErrDuplicatedKey doesn't catch the unique constraint violation
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			return nil, &customerrors.CompanyExistsError{Name: company.Name}
		}
		if strings.Contains(err.Error(), "connection refused") {
			return nil, &customerrors.DBConnectionError{}
		}
		return nil, &customerrors.GenericTxError{Msg: err.Error()}
	}

	return company, nil
}

// Update updates company in the database
func (r *PostgresRepository) Update(ctx context.Context, id uuid.UUID, company *entity.Company) (*entity.Company, error) {
	var existingCompany entity.Company
	if err := r.db.WithContext(ctx).First(&existingCompany, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &customerrors.RecordNotFoundError{ID: id}
		}
		if strings.Contains(err.Error(), "connection refused") {
			return nil, &customerrors.DBConnectionError{}
		}
		return nil, &customerrors.GenericTxError{Msg: err.Error()}
	}

	if company.ID != uuid.Nil && existingCompany.ID != company.ID {
		return nil, &customerrors.IDUpdateError{ID: id}
	}

	if err := r.db.WithContext(ctx).Model(&existingCompany).Updates(company).Error; err != nil {
		return nil, &customerrors.GenericTxError{Msg: err.Error()}
	}

	return &existingCompany, nil
}

// Delete deletes a company from the database
func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	var company entity.Company
	if err := r.db.WithContext(ctx).First(&company, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &customerrors.RecordNotFoundError{ID: id}
		}
		if strings.Contains(err.Error(), "connection refused") {
			return &customerrors.DBConnectionError{}
		}
		return &customerrors.GenericTxError{Msg: err.Error()}
	}

	if err := r.db.WithContext(ctx).Delete(&company).Error; err != nil {
		return &customerrors.GenericTxError{Msg: err.Error()}
	}

	return nil
}

// Get gets a company from the database
func (r *PostgresRepository) Get(ctx context.Context, id uuid.UUID) (*entity.Company, error) {
	var company entity.Company
	if err := r.db.WithContext(ctx).First(&company, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &customerrors.RecordNotFoundError{ID: id}
		}
		if strings.Contains(err.Error(), "connection refused") {
			return nil, &customerrors.DBConnectionError{}
		}
		return nil, &customerrors.GenericTxError{Msg: err.Error()}
	}
	return &company, nil
}
