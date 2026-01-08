package repository

import (
	"context"

	"github.com/ulumfr/ulumfr-be/pkg/models"
	"github.com/ulumfr/ulumfr-be/prisma/db"
)

type aboutRepository struct {
	client *db.PrismaClient
}

// NewAboutRepository creates a new about repository
func NewAboutRepository(client *db.PrismaClient) AboutRepository {
	return &aboutRepository{client: client}
}

func (r *aboutRepository) FindAll(ctx context.Context) ([]models.About, error) {
	abouts, err := r.client.About.FindMany().
		OrderBy(db.About.CreatedAt.Order(db.DESC)).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	result := make([]models.About, len(abouts))
	for i, a := range abouts {
		result[i] = *mapAboutToDomain(&a)
	}

	return result, nil
}

func (r *aboutRepository) FindByID(ctx context.Context, id string) (*models.About, error) {
	about, err := r.client.About.FindUnique(
		db.About.ID.Equals(id),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return mapAboutToDomain(about), nil
}

func (r *aboutRepository) FindActive(ctx context.Context) (*models.About, error) {
	about, err := r.client.About.FindFirst(
		db.About.IsActive.Equals(true),
	).OrderBy(db.About.UpdatedAt.Order(db.DESC)).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return mapAboutToDomain(about), nil
}

func (r *aboutRepository) Create(ctx context.Context, input models.CreateAboutInput) (*models.About, error) {
	about, err := r.client.About.CreateOne(
		db.About.FullName.Set(input.FullName),
		db.About.Role.Set(input.Role),
		db.About.IsActive.Set(input.IsActive),
		db.About.Nickname.SetIfPresent(input.Nickname),
		db.About.Bio.SetIfPresent(input.Bio),
		db.About.AvatarURL.SetIfPresent(input.AvatarURL),
		db.About.CoverURL.SetIfPresent(input.CoverURL),
		db.About.Location.SetIfPresent(input.Location),
		db.About.Email.SetIfPresent(input.Email),
		db.About.Phone.SetIfPresent(input.Phone),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return mapAboutToDomain(about), nil
}

func (r *aboutRepository) Update(ctx context.Context, id string, input models.UpdateAboutInput) (*models.About, error) {
	updates := []db.AboutSetParam{}

	if input.FullName != nil {
		updates = append(updates, db.About.FullName.Set(*input.FullName))
	}
	if input.Nickname != nil {
		updates = append(updates, db.About.Nickname.Set(*input.Nickname))
	}
	if input.Role != nil {
		updates = append(updates, db.About.Role.Set(*input.Role))
	}
	if input.Bio != nil {
		updates = append(updates, db.About.Bio.Set(*input.Bio))
	}
	if input.AvatarURL != nil {
		updates = append(updates, db.About.AvatarURL.Set(*input.AvatarURL))
	}
	if input.CoverURL != nil {
		updates = append(updates, db.About.CoverURL.Set(*input.CoverURL))
	}
	if input.Location != nil {
		updates = append(updates, db.About.Location.Set(*input.Location))
	}
	if input.Email != nil {
		updates = append(updates, db.About.Email.Set(*input.Email))
	}
	if input.Phone != nil {
		updates = append(updates, db.About.Phone.Set(*input.Phone))
	}
	if input.IsActive != nil {
		updates = append(updates, db.About.IsActive.Set(*input.IsActive))
	}

	about, err := r.client.About.FindUnique(
		db.About.ID.Equals(id),
	).Update(updates...).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return mapAboutToDomain(about), nil
}

func (r *aboutRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.About.FindUnique(
		db.About.ID.Equals(id),
	).Delete().Exec(ctx)
	return err
}

func mapAboutToDomain(a *db.AboutModel) *models.About {
	about := &models.About{
		ID:        a.ID,
		FullName:  a.FullName,
		Role:      a.Role,
		IsActive:  a.IsActive,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}

	if nickname, ok := a.Nickname(); ok {
		about.Nickname = &nickname
	}
	if bio, ok := a.Bio(); ok {
		about.Bio = &bio
	}
	if avatar, ok := a.AvatarURL(); ok {
		about.AvatarURL = &avatar
	}
	if cover, ok := a.CoverURL(); ok {
		about.CoverURL = &cover
	}
	if loc, ok := a.Location(); ok {
		about.Location = &loc
	}
	if email, ok := a.Email(); ok {
		about.Email = &email
	}
	if phone, ok := a.Phone(); ok {
		about.Phone = &phone
	}

	return about
}
