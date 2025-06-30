package repository

import (
	"context"
	"go-test-assesment/internal/cat/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresCatRepository struct {
	db *pgxpool.Pool
}

func NewPostgresCatRepository(db *pgxpool.Pool) domain.Repository {
	return &postgresCatRepository{db: db}
}

func (r *postgresCatRepository) Store(ctx context.Context, c *domain.Cat) error {
	query := `INSERT INTO cats (name, years_of_experience, breed, salary) VALUES ($1, $2, $3, $4) RETURNING id`
	return r.db.QueryRow(ctx, query, c.Name, c.YearsOfExperience, c.Breed, c.Salary).Scan(&c.ID)
}

func (r *postgresCatRepository) GetByID(ctx context.Context, id int64) (*domain.Cat, error) {
	query := `SELECT id, name, years_of_experience, breed, salary FROM cats WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var c domain.Cat
	err := row.Scan(&c.ID, &c.Name, &c.YearsOfExperience, &c.Breed, &c.Salary)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *postgresCatRepository) UpdateSalary(ctx context.Context, id int64, salary float64) error {
	query := `UPDATE cats SET salary = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, salary, id)
	return err
}

func (r *postgresCatRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM cats WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *postgresCatRepository) List(ctx context.Context) ([]*domain.Cat, error) {
	query := `SELECT id, name, years_of_experience, breed, salary FROM cats`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []*domain.Cat
	for rows.Next() {
		var c domain.Cat
		err := rows.Scan(&c.ID, &c.Name, &c.YearsOfExperience, &c.Breed, &c.Salary)
		if err != nil {
			return nil, err
		}
		cats = append(cats, &c)
	}

	return cats, nil
}
