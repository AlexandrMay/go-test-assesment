package usecase

import (
	"context"
	"errors"
	"fmt"
	cat "go-test-assesment/internal/cat/domain"
)

type CatUsecase struct {
	repo           cat.Repository
	breedValidator cat.BreedValidator
}

func NewCatUsecase(repo cat.Repository, bv cat.BreedValidator) *CatUsecase {
	return &CatUsecase{
		repo:           repo,
		breedValidator: bv,
	}
}

func (uc *CatUsecase) Create(ctx context.Context, c *cat.Cat) error {
	if c.Name == "" {
		return errors.New("cat name cannot be empty")
	}

	valid, err := uc.breedValidator.ValidateBreed(ctx, c.Breed)
	if err != nil {
		return fmt.Errorf("error validating breed: %w", err)
	}
	if !valid {
		return errors.New("invalid breed")
	}

	return uc.repo.Store(ctx, c)
}

func (uc *CatUsecase) GetByID(ctx context.Context, id int64) (*cat.Cat, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *CatUsecase) UpdateSalary(ctx context.Context, id int64, salary float64) error {
	if salary < 0 {
		return errors.New("salary cannot be negative")
	}
	return uc.repo.UpdateSalary(ctx, id, salary)
}

func (uc *CatUsecase) Delete(ctx context.Context, id int64) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *CatUsecase) List(ctx context.Context) ([]*cat.Cat, error) {
	return uc.repo.List(ctx)
}
