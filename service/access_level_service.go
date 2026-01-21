package service

import (
	"context"
	"fmt"

	"github.com/wabtcdi/user_service/dto"
	"github.com/wabtcdi/user_service/models"
	"github.com/wabtcdi/user_service/repository"
)

type AccessLevelService struct {
	repo repository.AccessLevelRepository
}

func NewAccessLevelService(repo repository.AccessLevelRepository) *AccessLevelService {
	return &AccessLevelService{repo: repo}
}

func (s *AccessLevelService) CreateAccessLevel(ctx context.Context, req *dto.CreateAccessLevelRequest) (*dto.AccessLevelResponse, error) {
	// Check if access level already exists
	existing, _ := s.repo.GetByName(ctx, req.Name)
	if existing != nil {
		return nil, fmt.Errorf("access level with name %s already exists", req.Name)
	}

	accessLevel := &models.AccessLevel{
		Name: req.Name,
	}
	if req.Description != "" {
		accessLevel.Description = &req.Description
	}

	err := s.repo.Create(ctx, accessLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to create access level: %w", err)
	}

	desc := ""
	if accessLevel.Description != nil {
		desc = *accessLevel.Description
	}

	return &dto.AccessLevelResponse{
		ID:          accessLevel.ID,
		Name:        accessLevel.Name,
		Description: desc,
	}, nil
}

func (s *AccessLevelService) GetAccessLevel(ctx context.Context, id int) (*dto.AccessLevelResponse, error) {
	accessLevel, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	desc := ""
	if accessLevel.Description != nil {
		desc = *accessLevel.Description
	}

	return &dto.AccessLevelResponse{
		ID:          accessLevel.ID,
		Name:        accessLevel.Name,
		Description: desc,
	}, nil
}

func (s *AccessLevelService) ListAccessLevels(ctx context.Context) ([]dto.AccessLevelResponse, error) {
	accessLevels, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.AccessLevelResponse, 0, len(accessLevels))
	for _, al := range accessLevels {
		desc := ""
		if al.Description != nil {
			desc = *al.Description
		}
		responses = append(responses, dto.AccessLevelResponse{
			ID:          al.ID,
			Name:        al.Name,
			Description: desc,
		})
	}

	return responses, nil
}
