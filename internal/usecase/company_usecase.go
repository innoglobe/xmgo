package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/innoglobe/xmgo/internal/entity"
	"github.com/innoglobe/xmgo/internal/interface/repository"
	eventservice "github.com/innoglobe/xmgo/internal/service"
)

type CompanyUsecaseInterface interface {
	CreateCompany(ctx context.Context, company *entity.Company) (*entity.Company, error)
	UpdateCompany(ctx context.Context, id uuid.UUID, company *entity.Company) (*entity.Company, error)
	DeleteCompany(ctx context.Context, id uuid.UUID) error
	GetCompany(ctx context.Context, id uuid.UUID) (*entity.Company, error)
}

type companyUsecase struct {
	repo          repository.CompanyRepositoryInterface
	eventProducer eventservice.Producer
}

func NewCompanyUsecase(repo repository.CompanyRepositoryInterface, eventProducer eventservice.Producer) CompanyUsecaseInterface {
	return &companyUsecase{repo: repo, eventProducer: eventProducer}
}

func (u *companyUsecase) CreateCompany(ctx context.Context, company *entity.Company) (*entity.Company, error) {
	if company == nil {
		return nil, errors.New("company can't be nil")
	}

	if err := company.Type.IsValid(); err != nil {
		return nil, err
	}

	res, err := u.repo.Create(ctx, company)
	if err != nil {
		return nil, err
	}

	u.eventProducer.Produce(&eventservice.Event{
		Operation: "create",
		Entity:    "company",
		Data:      *res,
	})
	return res, nil
}

func (u *companyUsecase) UpdateCompany(ctx context.Context, id uuid.UUID, company *entity.Company) (*entity.Company, error) {
	if company == nil {
		return nil, errors.New("company can't be nil")
	}

	if err := company.Type.IsValid(); err != nil {
		return nil, err
	}

	res, err := u.repo.Update(ctx, id, company)
	if err != nil {
		return nil, err
	}

	u.eventProducer.Produce(&eventservice.Event{
		Operation: "update",
		Entity:    "company",
		Data:      *res,
	})

	return res, nil
}

func (u *companyUsecase) DeleteCompany(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid id")
	}

	err := u.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	u.eventProducer.Produce(&eventservice.Event{
		Operation: "delete",
		Entity:    "company",
		Data:      id,
	})

	return nil
}

func (u *companyUsecase) GetCompany(ctx context.Context, id uuid.UUID) (*entity.Company, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid id")
	}
	return u.repo.Get(ctx, id)
}
