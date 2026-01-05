package repository

import (
	"context"

	"github.com/ulumfr/ulumfr-be/pkg/domain"
	"github.com/ulumfr/ulumfr-be/prisma/db"
)

type tagRepository struct {
	client *db.PrismaClient
}

// NewTagRepository creates a new tag repository
func NewTagRepository(client *db.PrismaClient) TagRepository {
	return &tagRepository{client: client}
}

func (r *tagRepository) FindAll(ctx context.Context) ([]domain.Tag, error) {
	tags, err := r.client.Tag.FindMany().
		OrderBy(db.Tag.Name.Order(db.ASC)).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	result := make([]domain.Tag, len(tags))
	for i, t := range tags {
		tag := domain.Tag{
			ID:        t.ID,
			Name:      t.Name,
			Slug:      t.Slug,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
		}
		if iconUrl, ok := t.IconURL(); ok {
			tag.IconUrl = &iconUrl
		}
		result[i] = tag
	}

	return result, nil
}

func (r *tagRepository) FindByID(ctx context.Context, id string) (*domain.Tag, error) {
	tag, err := r.client.Tag.FindUnique(
		db.Tag.ID.Equals(id),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	result := &domain.Tag{
		ID:        tag.ID,
		Name:      tag.Name,
		Slug:      tag.Slug,
		CreatedAt: tag.CreatedAt,
		UpdatedAt: tag.UpdatedAt,
	}
	if iconUrl, ok := tag.IconURL(); ok {
		result.IconUrl = &iconUrl
	}

	return result, nil
}

func (r *tagRepository) Create(ctx context.Context, input domain.CreateTagInput) (*domain.Tag, error) {
	tag, err := r.client.Tag.CreateOne(
		db.Tag.Name.Set(input.Name),
		db.Tag.Slug.Set(input.Slug),
		db.Tag.IconURL.SetIfPresent(input.IconUrl),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	result := &domain.Tag{
		ID:        tag.ID,
		Name:      tag.Name,
		Slug:      tag.Slug,
		CreatedAt: tag.CreatedAt,
		UpdatedAt: tag.UpdatedAt,
	}
	if iconUrl, ok := tag.IconURL(); ok {
		result.IconUrl = &iconUrl
	}

	return result, nil
}

func (r *tagRepository) Update(ctx context.Context, id string, input domain.UpdateTagInput) (*domain.Tag, error) {
	updates := []db.TagSetParam{}

	if input.Name != nil {
		updates = append(updates, db.Tag.Name.Set(*input.Name))
	}
	if input.Slug != nil {
		updates = append(updates, db.Tag.Slug.Set(*input.Slug))
	}
	if input.IconUrl != nil {
		updates = append(updates, db.Tag.IconURL.Set(*input.IconUrl))
	}

	tag, err := r.client.Tag.FindUnique(
		db.Tag.ID.Equals(id),
	).Update(updates...).Exec(ctx)

	if err != nil {
		return nil, err
	}

	result := &domain.Tag{
		ID:        tag.ID,
		Name:      tag.Name,
		Slug:      tag.Slug,
		CreatedAt: tag.CreatedAt,
		UpdatedAt: tag.UpdatedAt,
	}
	if iconUrl, ok := tag.IconURL(); ok {
		result.IconUrl = &iconUrl
	}

	return result, nil
}

func (r *tagRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.Tag.FindUnique(
		db.Tag.ID.Equals(id),
	).Delete().Exec(ctx)
	return err
}
