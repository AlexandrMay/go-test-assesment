package domain

import "context"

type Cat struct {
	ID                int64   `json:"id"`
	Name              string  `json:"name"`
	YearsOfExperience int     `json:"years_of_experience"`
	Breed             string  `json:"breed"`
	Salary            float64 `json:"salary"`
}

type Repository interface {
	Store(ctx context.Context, c *Cat) error
	GetByID(ctx context.Context, id int64) (*Cat, error)
	UpdateSalary(ctx context.Context, id int64, salary float64) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*Cat, error)
}
