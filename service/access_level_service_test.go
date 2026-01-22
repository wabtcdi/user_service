package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/wabtcdi/user_service/dto"
	"github.com/wabtcdi/user_service/mocks"
	"github.com/wabtcdi/user_service/models"
)

func TestAccessLevelService_CreateAccessLevel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAccessLevelRepository(ctrl)
	service := NewAccessLevelService(mockRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		req := &dto.CreateAccessLevelRequest{
			Name:        "Admin",
			Description: "Administrator access level",
		}

		mockRepo.EXPECT().
			GetByName(ctx, req.Name).
			Return(nil, errors.New("not found"))

		mockRepo.EXPECT().
			Create(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, al *models.AccessLevel) error {
				al.ID = 1
				return nil
			})

		resp, err := service.CreateAccessLevel(ctx, req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resp.Name != req.Name {
			t.Errorf("Expected name %s, got %s", req.Name, resp.Name)
		}
		if resp.Description != req.Description {
			t.Errorf("Expected description %s, got %s", req.Description, resp.Description)
		}
		if resp.ID != 1 {
			t.Errorf("Expected ID 1, got %d", resp.ID)
		}
	})

	t.Run("SuccessWithoutDescription", func(t *testing.T) {
		req := &dto.CreateAccessLevelRequest{
			Name: "User",
		}

		mockRepo.EXPECT().
			GetByName(ctx, req.Name).
			Return(nil, errors.New("not found"))

		mockRepo.EXPECT().
			Create(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, al *models.AccessLevel) error {
				al.ID = 2
				return nil
			})

		resp, err := service.CreateAccessLevel(ctx, req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resp.Name != req.Name {
			t.Errorf("Expected name %s, got %s", req.Name, resp.Name)
		}
		if resp.Description != "" {
			t.Errorf("Expected empty description, got %s", resp.Description)
		}
	})

	t.Run("AlreadyExists", func(t *testing.T) {
		req := &dto.CreateAccessLevelRequest{
			Name:        "Admin",
			Description: "Administrator access level",
		}

		existingAccessLevel := &models.AccessLevel{
			ID:   1,
			Name: "Admin",
		}

		mockRepo.EXPECT().
			GetByName(ctx, req.Name).
			Return(existingAccessLevel, nil)

		_, err := service.CreateAccessLevel(ctx, req)
		if err == nil {
			t.Fatal("Expected error for duplicate access level, got nil")
		}
	})

	t.Run("CreateError", func(t *testing.T) {
		req := &dto.CreateAccessLevelRequest{
			Name:        "Editor",
			Description: "Editor access level",
		}

		mockRepo.EXPECT().
			GetByName(ctx, req.Name).
			Return(nil, errors.New("not found"))

		mockRepo.EXPECT().
			Create(ctx, gomock.Any()).
			Return(errors.New("database error"))

		_, err := service.CreateAccessLevel(ctx, req)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestAccessLevelService_GetAccessLevel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAccessLevelRepository(ctrl)
	service := NewAccessLevelService(mockRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		desc := "Admin access level"
		accessLevel := &models.AccessLevel{
			ID:          1,
			Name:        "Admin",
			Description: &desc,
		}

		mockRepo.EXPECT().
			GetByID(ctx, 1).
			Return(accessLevel, nil)

		resp, err := service.GetAccessLevel(ctx, 1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resp.ID != 1 {
			t.Errorf("Expected ID 1, got %d", resp.ID)
		}
		if resp.Name != "Admin" {
			t.Errorf("Expected name Admin, got %s", resp.Name)
		}
		if resp.Description != desc {
			t.Errorf("Expected description %s, got %s", desc, resp.Description)
		}
	})

	t.Run("SuccessWithoutDescription", func(t *testing.T) {
		accessLevel := &models.AccessLevel{
			ID:          2,
			Name:        "User",
			Description: nil,
		}

		mockRepo.EXPECT().
			GetByID(ctx, 2).
			Return(accessLevel, nil)

		resp, err := service.GetAccessLevel(ctx, 2)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resp.ID != 2 {
			t.Errorf("Expected ID 2, got %d", resp.ID)
		}
		if resp.Description != "" {
			t.Errorf("Expected empty description, got %s", resp.Description)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.EXPECT().
			GetByID(ctx, 999).
			Return(nil, errors.New("not found"))

		_, err := service.GetAccessLevel(ctx, 999)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestAccessLevelService_ListAccessLevels(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAccessLevelRepository(ctrl)
	service := NewAccessLevelService(mockRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		desc1 := "Admin access level"
		desc2 := "User access level"
		accessLevels := []*models.AccessLevel{
			{
				ID:          1,
				Name:        "Admin",
				Description: &desc1,
			},
			{
				ID:          2,
				Name:        "User",
				Description: &desc2,
			},
		}

		mockRepo.EXPECT().
			List(ctx).
			Return(accessLevels, nil)

		resp, err := service.ListAccessLevels(ctx)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(resp) != 2 {
			t.Fatalf("Expected 2 access levels, got %d", len(resp))
		}

		if resp[0].Name != "Admin" {
			t.Errorf("Expected first name Admin, got %s", resp[0].Name)
		}
		if resp[1].Name != "User" {
			t.Errorf("Expected second name User, got %s", resp[1].Name)
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		mockRepo.EXPECT().
			List(ctx).
			Return([]*models.AccessLevel{}, nil)

		resp, err := service.ListAccessLevels(ctx)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(resp) != 0 {
			t.Errorf("Expected 0 access levels, got %d", len(resp))
		}
	})

	t.Run("MixedDescriptions", func(t *testing.T) {
		desc := "Admin access level"
		accessLevels := []*models.AccessLevel{
			{
				ID:          1,
				Name:        "Admin",
				Description: &desc,
			},
			{
				ID:          2,
				Name:        "User",
				Description: nil,
			},
		}

		mockRepo.EXPECT().
			List(ctx).
			Return(accessLevels, nil)

		resp, err := service.ListAccessLevels(ctx)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resp[0].Description != desc {
			t.Errorf("Expected first description %s, got %s", desc, resp[0].Description)
		}
		if resp[1].Description != "" {
			t.Errorf("Expected second description empty, got %s", resp[1].Description)
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo.EXPECT().
			List(ctx).
			Return(nil, errors.New("database error"))

		_, err := service.ListAccessLevels(ctx)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}
