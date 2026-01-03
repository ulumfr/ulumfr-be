package repository

import (
	"context"

	"github.com/ulumfr/ulumfr-be/pkg/domain"
	"github.com/ulumfr/ulumfr-be/prisma/db"
)

type careerRepository struct {
	client *db.PrismaClient
}

// NewCareerRepository creates a new career repository
func NewCareerRepository(client *db.PrismaClient) CareerRepository {
	return &careerRepository{client: client}
}

func (r *careerRepository) FindAll(ctx context.Context) ([]domain.Career, error) {
	careers, err := r.client.Career.FindMany().
		OrderBy(db.Career.SortOrder.Order(db.ASC), db.Career.StartDate.Order(db.DESC)).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	result := make([]domain.Career, len(careers))
	for i, c := range careers {
		result[i] = *mapCareerToDomain(&c)
	}

	return result, nil
}

func (r *careerRepository) FindByID(ctx context.Context, id string) (*domain.Career, error) {
	career, err := r.client.Career.FindUnique(
		db.Career.ID.Equals(id),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return mapCareerToDomain(career), nil
}

func (r *careerRepository) Create(ctx context.Context, input domain.CreateCareerInput) (*domain.Career, error) {
	career, err := r.client.Career.CreateOne(
		db.Career.Company.Set(input.Company),
		db.Career.Position.Set(input.Position),
		db.Career.StartDate.Set(input.StartDate),
		db.Career.IsCurrent.Set(input.IsCurrent),
		db.Career.SortOrder.Set(input.SortOrder),
		db.Career.Location.SetIfPresent(input.Location),
		db.Career.Description.SetIfPresent(input.Description),
		db.Career.EndDate.SetIfPresent(input.EndDate),
		db.Career.LogoURL.SetIfPresent(input.LogoURL),
		db.Career.CompanyURL.SetIfPresent(input.CompanyURL),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return mapCareerToDomain(career), nil
}

func (r *careerRepository) Update(ctx context.Context, id string, input domain.UpdateCareerInput) (*domain.Career, error) {
	updates := []db.CareerSetParam{}

	if input.Company != nil {
		updates = append(updates, db.Career.Company.Set(*input.Company))
	}
	if input.Position != nil {
		updates = append(updates, db.Career.Position.Set(*input.Position))
	}
	if input.Location != nil {
		updates = append(updates, db.Career.Location.Set(*input.Location))
	}
	if input.Description != nil {
		updates = append(updates, db.Career.Description.Set(*input.Description))
	}
	if input.StartDate != nil {
		updates = append(updates, db.Career.StartDate.Set(*input.StartDate))
	}
	if input.EndDate != nil {
		updates = append(updates, db.Career.EndDate.Set(*input.EndDate))
	}
	if input.IsCurrent != nil {
		updates = append(updates, db.Career.IsCurrent.Set(*input.IsCurrent))
	}
	if input.LogoURL != nil {
		updates = append(updates, db.Career.LogoURL.Set(*input.LogoURL))
	}
	if input.CompanyURL != nil {
		updates = append(updates, db.Career.CompanyURL.Set(*input.CompanyURL))
	}
	if input.SortOrder != nil {
		updates = append(updates, db.Career.SortOrder.Set(*input.SortOrder))
	}

	career, err := r.client.Career.FindUnique(
		db.Career.ID.Equals(id),
	).Update(updates...).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return mapCareerToDomain(career), nil
}

func (r *careerRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.Career.FindUnique(
		db.Career.ID.Equals(id),
	).Delete().Exec(ctx)
	return err
}

func mapCareerToDomain(c *db.CareerModel) *domain.Career {
	career := &domain.Career{
		ID:        c.ID,
		Company:   c.Company,
		Position:  c.Position,
		StartDate: c.StartDate,
		IsCurrent: c.IsCurrent,
		SortOrder: c.SortOrder,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}

	if loc, ok := c.Location(); ok {
		career.Location = &loc
	}
	if desc, ok := c.Description(); ok {
		career.Description = &desc
	}
	if endDate, ok := c.EndDate(); ok {
		career.EndDate = &endDate
	}
	if logo, ok := c.LogoURL(); ok {
		career.LogoURL = &logo
	}
	if url, ok := c.CompanyURL(); ok {
		career.CompanyURL = &url
	}

	return career
}
