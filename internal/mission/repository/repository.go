package repository

import (
	"context"
	"errors"
	"go-test-assesment/internal/mission/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MissionPostgres struct {
	pool *pgxpool.Pool
}

func NewMissionPostgres(pool *pgxpool.Pool) *MissionPostgres {
	return &MissionPostgres{pool: pool}
}

func (r *MissionPostgres) CreateMission(ctx context.Context, m *domain.Mission) error {
	query := `
		INSERT INTO missions (cat_id, completed, created_at, updated_at)
		VALUES ($1, $2, now(), now())
		RETURNING id, created_at, updated_at`
	return r.pool.QueryRow(ctx, query, m.CatID, m.Completed).
		Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
}

func (r *MissionPostgres) GetMissionByID(ctx context.Context, id int64) (*domain.Mission, error) {
	m := &domain.Mission{}
	query := `SELECT id, cat_id, completed, created_at, updated_at FROM missions WHERE id = $1`
	err := r.pool.QueryRow(ctx, query, id).
		Scan(&m.ID, &m.CatID, &m.Completed, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}

	targets, err := r.listTargetsByMissionID(ctx, m.ID)
	if err != nil {
		return nil, err
	}
	m.Targets = targets

	return m, nil
}

func (r *MissionPostgres) listTargetsByMissionID(ctx context.Context, missionID int64) ([]domain.Target, error) {
	query := `
		SELECT id, mission_id, name, country, notes, completed, created_at, updated_at
		FROM targets WHERE mission_id = $1`
	rows, err := r.pool.Query(ctx, query, missionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var targets []domain.Target
	for rows.Next() {
		var t domain.Target
		if err := rows.Scan(
			&t.ID, &t.MissionID, &t.Name, &t.Country,
			&t.Notes, &t.Completed, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		targets = append(targets, t)
	}
	return targets, nil
}

func (r *MissionPostgres) ListMissions(ctx context.Context) ([]*domain.Mission, error) {
	query := `SELECT id, cat_id, completed, created_at, updated_at FROM missions`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var missions []*domain.Mission
	for rows.Next() {
		m := &domain.Mission{}
		if err := rows.Scan(&m.ID, &m.CatID, &m.Completed, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		targets, err := r.listTargetsByMissionID(ctx, m.ID)
		if err != nil {
			return nil, err
		}
		m.Targets = targets
		missions = append(missions, m)
	}
	return missions, nil
}

func (r *MissionPostgres) UpdateMission(ctx context.Context, m *domain.Mission) error {
	query := `
		UPDATE missions
		SET cat_id = $1, completed = $2, updated_at = now()
		WHERE id = $3`
	res, err := r.pool.Exec(ctx, query, m.CatID, m.Completed, m.ID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("mission not found")
	}
	return nil
}

func (r *MissionPostgres) DeleteMission(ctx context.Context, id int64) error {
	var catID int64
	err := r.pool.QueryRow(ctx, `SELECT cat_id FROM missions WHERE id = $1`, id).Scan(&catID)
	if err != nil {
		return err
	}
	if catID != 0 {
		return errors.New("cannot delete mission assigned to a cat")
	}

	res, err := r.pool.Exec(ctx, `DELETE FROM missions WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("mission not found")
	}
	return nil
}

func (r *MissionPostgres) AddTargets(ctx context.Context, targets []domain.Target) error {
	for _, t := range targets {
		if t.MissionID == 0 {
			return errors.New("target must have mission_id")
		}
		_, err := r.pool.Exec(ctx,
			`INSERT INTO targets (mission_id, name, country, notes, completed, created_at, updated_at)
			 VALUES ($1, $2, $3, $4, $5, now(), now())`,
			t.MissionID, t.Name, t.Country, t.Notes, t.Completed,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *MissionPostgres) UpdateTarget(ctx context.Context, t *domain.Target) error {
	query := `
		UPDATE targets
		SET notes = $1, completed = $2, updated_at = now()
		WHERE id = $3`
	res, err := r.pool.Exec(ctx, query, t.Notes, t.Completed, t.ID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("target not found")
	}
	return nil
}

func (r *MissionPostgres) GetTargetByID(ctx context.Context, id int64) (*domain.Target, error) {
	var target domain.Target
	query := `SELECT id, mission_id, name, country, notes, completed, created_at, updated_at FROM targets WHERE id = $1`
	err := r.pool.QueryRow(ctx, query, id).
		Scan(&target.ID, &target.MissionID, &target.Name, &target.Country, &target.Notes,
			&target.Completed, &target.CreatedAt, &target.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &target, nil
}

func (r *MissionPostgres) DeleteTarget(ctx context.Context, id int64) error {
	var completed bool
	err := r.pool.QueryRow(ctx, `SELECT completed FROM targets WHERE id = $1`, id).Scan(&completed)
	if err != nil {
		return err
	}
	if completed {
		return errors.New("cannot delete a completed target")
	}

	res, err := r.pool.Exec(ctx, `DELETE FROM targets WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("target not found")
	}
	return nil
}
