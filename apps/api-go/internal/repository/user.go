package repository

import (
	"context"

	"github.com/dogab/vitalstack/api/internal/models"
	"github.com/supabase-community/supabase-go"
)

type UserRepository interface {
	GetProfile(ctx context.Context, id string) (*models.Profile, error)
	CreateProfile(ctx context.Context, profile *models.Profile) error
	UpdateProfile(ctx context.Context, profile *models.Profile) error
}

type userRepository struct {
	client *supabase.Client
}

func NewUserRepository(client *supabase.Client) UserRepository {
	return &userRepository{client: client}
}

func (r *userRepository) GetProfile(ctx context.Context, id string) (*models.Profile, error) {
	var profile models.Profile
	_, err := r.client.From("public.profiles").Select("*", "exact", false).Eq("id", id).Single().ExecuteTo(&profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *userRepository) CreateProfile(ctx context.Context, profile *models.Profile) error {
	_, _, err := r.client.From("public.profiles").Insert(profile, false, "", "", "").Execute()
	return err
}

func (r *userRepository) UpdateProfile(ctx context.Context, profile *models.Profile) error {
	_, _, err := r.client.From("public.profiles").Update(profile, "", "exact").Eq("id", profile.ID).Execute()
	return err
}
