package domain

import (
	"context"
	"time"
)

type Mission struct {
	ID        int64     `json:"id"`
	CatID     *int64    `json:"cat_id,omitempty"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Targets   []Target  `json:"targets,omitempty"`
}

type Target struct {
	ID        int64     `json:"id"`
	MissionID int64     `json:"mission_id"`
	Name      string    `json:"name"`
	Country   string    `json:"country"`
	Notes     string    `json:"notes"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Repository interface {
	CreateMission(ctx context.Context, mission *Mission) error
	GetMissionByID(ctx context.Context, id int64) (*Mission, error)
	ListMissions(ctx context.Context) ([]*Mission, error)
	UpdateMission(ctx context.Context, mission *Mission) error
	DeleteMission(ctx context.Context, id int64) error
	GetTargetByID(ctx context.Context, id int64) (*Target, error)
	AddTargets(ctx context.Context, targets []Target) error
	UpdateTarget(ctx context.Context, target *Target) error
	DeleteTarget(ctx context.Context, id int64) error
}

type Usecase interface {
	CreateMission(ctx context.Context, mission *Mission) error
	GetMissionByID(ctx context.Context, id int64) (*Mission, error)
	ListMissions(ctx context.Context) ([]*Mission, error)
	UpdateMission(ctx context.Context, mission *Mission) error
	DeleteMission(ctx context.Context, id int64) error

	AddTargets(ctx context.Context, missionID int64, targets []Target) error
	UpdateTarget(ctx context.Context, target *Target) error
	DeleteTarget(ctx context.Context, targetID int64) error

	AssignCatToMission(ctx context.Context, missionID, catID int64) error
}
