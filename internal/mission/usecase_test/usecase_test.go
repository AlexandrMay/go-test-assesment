package usecase_test

import (
	"context"
	"fmt"
	"go-test-assesment/internal/mission/domain"
	"go-test-assesment/internal/mission/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateMission(ctx context.Context, mission *domain.Mission) error {
	args := m.Called(ctx, mission)
	return args.Error(0)
}

func (m *MockRepository) GetMissionByID(ctx context.Context, id int64) (*domain.Mission, error) {
	args := m.Called(ctx, id)
	if obj := args.Get(0); obj != nil {
		return obj.(*domain.Mission), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) ListMissions(ctx context.Context) ([]*domain.Mission, error) {
	args := m.Called(ctx)
	if obj := args.Get(0); obj != nil {
		return obj.([]*domain.Mission), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) UpdateMission(ctx context.Context, mission *domain.Mission) error {
	args := m.Called(ctx, mission)
	return args.Error(0)
}

func (m *MockRepository) DeleteMission(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) AddTargets(ctx context.Context, targets []domain.Target) error {
	args := m.Called(ctx, targets)
	return args.Error(0)
}

func (m *MockRepository) UpdateTarget(ctx context.Context, target *domain.Target) error {
	args := m.Called(ctx, target)
	return args.Error(0)
}

func (m *MockRepository) GetTargetByID(ctx context.Context, id int64) (*domain.Target, error) {
	args := m.Called(ctx, id)
	if obj := args.Get(0); obj != nil {
		return obj.(*domain.Target), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) DeleteTarget(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestMissionUsecase_CreateMission(t *testing.T) {
	mockRepo := new(MockRepository)
	mockRepo.On("CreateMission", mock.Anything, mock.AnythingOfType("*domain.Mission")).Return(nil)

	uc := usecase.NewMissionUsecase(mockRepo)

	err := uc.CreateMission(context.Background(), &domain.Mission{})
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMissionUsecase_GetMissionByID(t *testing.T) {
	mockRepo := new(MockRepository)
	expectedMission := &domain.Mission{ID: 42}
	mockRepo.On("GetMissionByID", mock.Anything, int64(42)).Return(expectedMission, nil)

	uc := usecase.NewMissionUsecase(mockRepo)

	mission, err := uc.GetMissionByID(context.Background(), 42)
	assert.NoError(t, err)
	assert.Equal(t, expectedMission, mission)
	mockRepo.AssertExpectations(t)
}

func TestMissionUsecase_ListMissions(t *testing.T) {
	mockRepo := new(MockRepository)
	missions := []*domain.Mission{{ID: 1}, {ID: 2}}
	mockRepo.On("ListMissions", mock.Anything).Return(missions, nil)

	uc := usecase.NewMissionUsecase(mockRepo)

	result, err := uc.ListMissions(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, missions, result)
	mockRepo.AssertExpectations(t)
}

func TestMissionUsecase_UpdateMission(t *testing.T) {
	mockRepo := new(MockRepository)

	mockRepo.On("GetMissionByID", mock.Anything, int64(1)).Return(&domain.Mission{ID: 1, Completed: false}, nil)
	mockRepo.On("UpdateMission", mock.Anything, mock.AnythingOfType("*domain.Mission")).Return(nil)

	uc := usecase.NewMissionUsecase(mockRepo)
	err := uc.UpdateMission(context.Background(), &domain.Mission{ID: 1})
	assert.NoError(t, err)

	mockRepo.ExpectedCalls = nil
	mockRepo.On("GetMissionByID", mock.Anything, int64(1)).Return(&domain.Mission{ID: 1, Completed: true}, nil)

	err = uc.UpdateMission(context.Background(), &domain.Mission{ID: 1})
	assert.EqualError(t, err, "cannot update a completed mission")

	mockRepo.AssertExpectations(t)
}

func TestMissionUsecase_DeleteMission(t *testing.T) {
	mockRepo := new(MockRepository)

	mockRepo.On("GetMissionByID", mock.Anything, int64(1)).Return(&domain.Mission{ID: 1, CatID: nil}, nil)
	mockRepo.On("DeleteMission", mock.Anything, int64(1)).Return(nil)

	uc := usecase.NewMissionUsecase(mockRepo)
	err := uc.DeleteMission(context.Background(), 1)
	assert.NoError(t, err)

	mockRepo.ExpectedCalls = nil
	catID := int64(5)
	mockRepo.On("GetMissionByID", mock.Anything, int64(1)).Return(&domain.Mission{ID: 1, CatID: &catID}, nil)

	err = uc.DeleteMission(context.Background(), 1)
	assert.EqualError(t, err, "mission cannot be deleted because it is assigned to a cat")

	mockRepo.AssertExpectations(t)
}

func TestMissionUsecase_AddTargets(t *testing.T) {
	mockRepo := new(MockRepository)

	mockRepo.On("GetMissionByID", mock.Anything, int64(1)).Return(&domain.Mission{ID: 1, Completed: false}, nil)
	mockRepo.On("AddTargets", mock.Anything, mock.AnythingOfType("[]domain.Target")).Return(nil)

	uc := usecase.NewMissionUsecase(mockRepo)
	err := uc.AddTargets(context.Background(), 1, []domain.Target{{Name: "Target1"}})
	assert.NoError(t, err)

	mockRepo.ExpectedCalls = nil
	mockRepo.On("GetMissionByID", mock.Anything, int64(1)).Return(&domain.Mission{ID: 1, Completed: true}, nil)

	err = uc.AddTargets(context.Background(), 1, []domain.Target{{Name: "Target1"}})
	assert.EqualError(t, err, "cannot add targets to a completed mission")

	mockRepo.AssertExpectations(t)
}

func TestMissionUsecase_UpdateTarget(t *testing.T) {
	mockRepo := new(MockRepository)

	existingTarget := &domain.Target{
		ID:        1,
		MissionID: 10,
		Notes:     "old notes",
		Completed: false,
	}

	mission := &domain.Mission{
		ID:        10,
		Completed: false,
	}

	mockRepo.On("GetTargetByID", mock.Anything, int64(1)).Return(existingTarget, nil)
	mockRepo.On("GetMissionByID", mock.Anything, int64(10)).Return(mission, nil)

	mockRepo.On("UpdateTarget", mock.Anything, mock.MatchedBy(func(t *domain.Target) bool {
		fmt.Printf("UpdateTarget called with: %+v\n", t)
		return true
	})).Return(nil)

	uc := usecase.NewMissionUsecase(mockRepo)

	err := uc.UpdateTarget(context.Background(), &domain.Target{
		ID:        1,
		MissionID: 10,
		Notes:     "new notes",
		Completed: false,
	})
	assert.NoError(t, err)
}

func TestMissionUsecase_DeleteTarget(t *testing.T) {
	mockRepo := new(MockRepository)

	mockRepo.On("GetTargetByID", mock.Anything, int64(1)).Return(&domain.Target{ID: 1, Completed: false}, nil)
	mockRepo.On("DeleteTarget", mock.Anything, int64(1)).Return(nil)

	uc := usecase.NewMissionUsecase(mockRepo)
	err := uc.DeleteTarget(context.Background(), 1)
	assert.NoError(t, err)

	mockRepo.ExpectedCalls = nil
	mockRepo.On("GetTargetByID", mock.Anything, int64(1)).Return(&domain.Target{ID: 1, Completed: true}, nil)

	err = uc.DeleteTarget(context.Background(), 1)
	assert.EqualError(t, err, "cannot delete a completed target")

	mockRepo.AssertExpectations(t)
}

func TestMissionUsecase_AssignCatToMission(t *testing.T) {
	mockRepo := new(MockRepository)

	mockRepo.On("GetMissionByID", mock.Anything, int64(1)).Return(&domain.Mission{ID: 1, CatID: nil}, nil)
	mockRepo.On("UpdateMission", mock.Anything, mock.AnythingOfType("*domain.Mission")).Return(nil)

	uc := usecase.NewMissionUsecase(mockRepo)
	err := uc.AssignCatToMission(context.Background(), 1, 42)
	assert.NoError(t, err)

	mockRepo.ExpectedCalls = nil
	catID := int64(5)
	mockRepo.On("GetMissionByID", mock.Anything, int64(1)).Return(&domain.Mission{ID: 1, CatID: &catID}, nil)

	err = uc.AssignCatToMission(context.Background(), 1, 42)
	assert.EqualError(t, err, "mission already assigned to a cat")

	mockRepo.AssertExpectations(t)
}
