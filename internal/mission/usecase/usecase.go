package usecase

import (
	"context"
	"errors"
	"go-test-assesment/internal/mission/domain"
)

type MissionUsecase struct {
	missionRepo domain.Repository
}

func NewMissionUsecase(mr domain.Repository) *MissionUsecase {
	return &MissionUsecase{missionRepo: mr}
}

func (uc *MissionUsecase) CreateMission(ctx context.Context, m *domain.Mission) error {
	return uc.missionRepo.CreateMission(ctx, m)
}

func (uc *MissionUsecase) GetMissionByID(ctx context.Context, id int64) (*domain.Mission, error) {
	return uc.missionRepo.GetMissionByID(ctx, id)
}

func (uc *MissionUsecase) ListMissions(ctx context.Context) ([]*domain.Mission, error) {
	return uc.missionRepo.ListMissions(ctx)
}

func (uc *MissionUsecase) UpdateMission(ctx context.Context, m *domain.Mission) error {
	existing, err := uc.missionRepo.GetMissionByID(ctx, m.ID)
	if err != nil {
		return err
	}
	if existing.Completed {
		return errors.New("cannot update a completed mission")
	}
	return uc.missionRepo.UpdateMission(ctx, m)
}

func (uc *MissionUsecase) DeleteMission(ctx context.Context, id int64) error {
	m, err := uc.missionRepo.GetMissionByID(ctx, id)
	if err != nil {
		return err
	}
	if m.CatID != nil && *m.CatID != 0 {
		return errors.New("mission cannot be deleted because it is assigned to a cat")
	}
	return uc.missionRepo.DeleteMission(ctx, id)
}

func (uc *MissionUsecase) AddTargets(ctx context.Context, missionID int64, targets []domain.Target) error {
	mission, err := uc.missionRepo.GetMissionByID(ctx, missionID)
	if err != nil {
		return err
	}
	if mission.Completed {
		return errors.New("cannot add targets to a completed mission")
	}
	return uc.missionRepo.AddTargets(ctx, targets)
}

func (uc *MissionUsecase) UpdateTarget(ctx context.Context, t *domain.Target) error {
	existingTarget, err := uc.missionRepo.GetTargetByID(ctx, t.ID)
	if err != nil {
		return err
	}
	mission, err := uc.missionRepo.GetMissionByID(ctx, t.MissionID)
	if err != nil {
		return err
	}
	if existingTarget.Completed || mission.Completed {
		if t.Notes != existingTarget.Notes {
			return errors.New("cannot update notes because target or mission is completed")
		}
	}
	return uc.missionRepo.UpdateTarget(ctx, t)
}

func (uc *MissionUsecase) DeleteTarget(ctx context.Context, id int64) error {
	t, err := uc.missionRepo.GetTargetByID(ctx, id)
	if err != nil {
		return err
	}
	if t.Completed {
		return errors.New("cannot delete a completed target")
	}
	return uc.missionRepo.DeleteTarget(ctx, id)
}

func (uc *MissionUsecase) AssignCatToMission(ctx context.Context, missionID, catID int64) error {
	mission, err := uc.missionRepo.GetMissionByID(ctx, missionID)
	if err != nil {
		return err
	}
	if mission.CatID != nil && *mission.CatID != 0 {
		return errors.New("mission already assigned to a cat")
	}
	mission.CatID = &catID
	return uc.missionRepo.UpdateMission(ctx, mission)
}
