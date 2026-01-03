package repository

import (
	"context"

	"github.com/ulumfr/ulumfr-be/pkg/domain"
	"github.com/ulumfr/ulumfr-be/prisma/db"
)

type categoryRepository struct {
	client *db.PrismaClient
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(client *db.PrismaClient) CategoryRepository {
	return &categoryRepository{client: client}
}

func (r *categoryRepository) FindAll(ctx context.Context) ([]domain.Category, error) {
	categories, err := r.client.Category.FindMany().
		OrderBy(db.Category.Name.Order(db.ASC)).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	result := make([]domain.Category, len(categories))
	for i, c := range categories {
		result[i] = *mapCategoryToDomain(&c)
	}

	return result, nil
}

func (r *categoryRepository) FindByID(ctx context.Context, id string) (*domain.Category, error) {
	category, err := r.client.Category.FindUnique(
		db.Category.ID.Equals(id),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return mapCategoryToDomain(category), nil
}

func (r *categoryRepository) Create(ctx context.Context, input domain.CreateCategoryInput) (*domain.Category, error) {
	category, err := r.client.Category.CreateOne(
		db.Category.Name.Set(input.Name),
		db.Category.Slug.Set(input.Slug),
		db.Category.Description.SetIfPresent(input.Description),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return mapCategoryToDomain(category), nil
}

func (r *categoryRepository) Update(ctx context.Context, id string, input domain.UpdateCategoryInput) (*domain.Category, error) {
	updates := []db.CategorySetParam{}

	if input.Name != nil {
		updates = append(updates, db.Category.Name.Set(*input.Name))
	}
	if input.Slug != nil {
		updates = append(updates, db.Category.Slug.Set(*input.Slug))
	}
	if input.Description != nil {
		updates = append(updates, db.Category.Description.Set(*input.Description))
	}

	category, err := r.client.Category.FindUnique(
		db.Category.ID.Equals(id),
	).Update(updates...).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return mapCategoryToDomain(category), nil
}

func (r *categoryRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.Category.FindUnique(
		db.Category.ID.Equals(id),
	).Delete().Exec(ctx)
	return err
}

func mapCategoryToDomain(c *db.CategoryModel) *domain.Category {
	cat := &domain.Category{
		ID:        c.ID,
		Name:      c.Name,
		Slug:      c.Slug,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}

	if desc, ok := c.Description(); ok {
		cat.Description = &desc
	}

	return cat
}
