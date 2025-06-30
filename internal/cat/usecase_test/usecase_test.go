package usecase_test

import (
	"context"
	"errors"
	"testing"

	cat "go-test-assesment/internal/cat/domain"
	"go-test-assesment/internal/cat/usecase"
)

type mockCatRepo struct {
	storeFn        func(ctx context.Context, c *cat.Cat) error
	getByIDFn      func(ctx context.Context, id int64) (*cat.Cat, error)
	updateSalaryFn func(ctx context.Context, id int64, salary float64) error
	deleteFn       func(ctx context.Context, id int64) error
	listFn         func(ctx context.Context) ([]*cat.Cat, error)
}

func (m *mockCatRepo) Store(ctx context.Context, c *cat.Cat) error {
	return m.storeFn(ctx, c)
}
func (m *mockCatRepo) GetByID(ctx context.Context, id int64) (*cat.Cat, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockCatRepo) UpdateSalary(ctx context.Context, id int64, salary float64) error {
	return m.updateSalaryFn(ctx, id, salary)
}
func (m *mockCatRepo) Delete(ctx context.Context, id int64) error {
	return m.deleteFn(ctx, id)
}
func (m *mockCatRepo) List(ctx context.Context) ([]*cat.Cat, error) {
	return m.listFn(ctx)
}

type mockBreedValidator struct {
	validateFn func(ctx context.Context, breed string) (bool, error)
}

func (m *mockBreedValidator) ValidateBreed(ctx context.Context, breed string) (bool, error) {
	return m.validateFn(ctx, breed)
}

func TestCatUsecase_Create(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		inputCat       *cat.Cat
		repoStoreErr   error
		validatorValid bool
		validatorErr   error
		wantErr        bool
	}{
		{
			name:           "success",
			inputCat:       &cat.Cat{Name: "Tom", Breed: "Siamese", Salary: 1000, YearsOfExperience: 2},
			repoStoreErr:   nil,
			validatorValid: true,
			validatorErr:   nil,
			wantErr:        false,
		},
		{
			name:     "empty name",
			inputCat: &cat.Cat{Name: "", Breed: "Siamese", Salary: 1000},
			wantErr:  true,
		},
		{
			name:           "invalid breed",
			inputCat:       &cat.Cat{Name: "Tom", Breed: "Unknown", Salary: 1000},
			validatorValid: false,
			wantErr:        true,
		},
		{
			name:         "breed validator error",
			inputCat:     &cat.Cat{Name: "Tom", Breed: "Siamese", Salary: 1000},
			validatorErr: errors.New("API error"),
			wantErr:      true,
		},
		{
			name:           "repository store error",
			inputCat:       &cat.Cat{Name: "Tom", Breed: "Siamese", Salary: 1000},
			repoStoreErr:   errors.New("db error"),
			validatorValid: true,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockCatRepo{
				storeFn: func(ctx context.Context, c *cat.Cat) error {
					return tt.repoStoreErr
				},
			}
			validator := &mockBreedValidator{
				validateFn: func(ctx context.Context, breed string) (bool, error) {
					return tt.validatorValid, tt.validatorErr
				},
			}

			uc := usecase.NewCatUsecase(repo, validator)
			err := uc.Create(ctx, tt.inputCat)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCatUsecase_UpdateSalary(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		id         int64
		salary     float64
		repoErr    error
		wantErr    bool
		errMessage string
	}{
		{
			name:    "success",
			id:      1,
			salary:  2000,
			repoErr: nil,
			wantErr: false,
		},
		{
			name:       "negative salary",
			id:         1,
			salary:     -10,
			wantErr:    true,
			errMessage: "salary cannot be negative",
		},
		{
			name:    "repository error",
			id:      1,
			salary:  1000,
			repoErr: errors.New("db update error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockCatRepo{
				updateSalaryFn: func(ctx context.Context, id int64, salary float64) error {
					return tt.repoErr
				},
			}

			validator := &mockBreedValidator{}

			uc := usecase.NewCatUsecase(repo, validator)

			err := uc.UpdateSalary(ctx, tt.id, tt.salary)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateSalary() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.errMessage != "" && err != nil && err.Error() != tt.errMessage {
				t.Errorf("UpdateSalary() error message = %v, want %v", err.Error(), tt.errMessage)
			}
		})
	}
}

func TestCatUsecase_Delete(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		id      int64
		repoErr error
		wantErr bool
	}{
		{
			name:    "success",
			id:      1,
			repoErr: nil,
			wantErr: false,
		},
		{
			name:    "repository error",
			id:      2,
			repoErr: errors.New("db delete error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockCatRepo{
				deleteFn: func(ctx context.Context, id int64) error {
					return tt.repoErr
				},
			}

			validator := &mockBreedValidator{}

			uc := usecase.NewCatUsecase(repo, validator)

			err := uc.Delete(ctx, tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCatUsecase_GetByID(t *testing.T) {
	ctx := context.Background()

	expectedCat := &cat.Cat{
		ID:                1,
		Name:              "Tom",
		Breed:             "Siamese",
		Salary:            1500,
		YearsOfExperience: 3,
	}

	tests := []struct {
		name    string
		id      int64
		cat     *cat.Cat
		repoErr error
		wantErr bool
	}{
		{
			name:    "success",
			id:      1,
			cat:     expectedCat,
			repoErr: nil,
			wantErr: false,
		},
		{
			name:    "repository error",
			id:      2,
			cat:     nil,
			repoErr: errors.New("not found"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockCatRepo{
				getByIDFn: func(ctx context.Context, id int64) (*cat.Cat, error) {
					return tt.cat, tt.repoErr
				},
			}

			validator := &mockBreedValidator{}

			uc := usecase.NewCatUsecase(repo, validator)

			c, err := uc.GetByID(ctx, tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
			}

			if c != tt.cat {
				t.Errorf("GetByID() cat = %v, want %v", c, tt.cat)
			}
		})
	}
}

func TestCatUsecase_List(t *testing.T) {
	ctx := context.Background()

	expectedCats := []*cat.Cat{
		{ID: 1, Name: "Tom", Breed: "Siamese", Salary: 1500},
		{ID: 2, Name: "Jerry", Breed: "Persian", Salary: 1800},
	}

	tests := []struct {
		name    string
		cats    []*cat.Cat
		repoErr error
		wantErr bool
	}{
		{
			name:    "success",
			cats:    expectedCats,
			repoErr: nil,
			wantErr: false,
		},
		{
			name:    "repository error",
			cats:    nil,
			repoErr: errors.New("db error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockCatRepo{
				listFn: func(ctx context.Context) ([]*cat.Cat, error) {
					return tt.cats, tt.repoErr
				},
			}

			validator := &mockBreedValidator{}

			uc := usecase.NewCatUsecase(repo, validator)

			cats, err := uc.List(ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(cats) != len(tt.cats) {
				t.Errorf("List() cats count = %d, want %d", len(cats), len(tt.cats))
			}
		})
	}
}
