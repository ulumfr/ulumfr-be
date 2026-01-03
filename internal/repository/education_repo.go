package repository

import (
	"context"

	"github.com/ulumfr/ulumfr-be/internal/domain"
	"github.com/ulumfr/ulumfr-be/prisma/db"
)

type educationRepository struct {
	client *db.PrismaClient
}

// NewEducationRepository creates a new education repository
func NewEducationRepository(client *db.PrismaClient) EducationRepository {
	return &educationRepository{client: client}
}

func (r *educationRepository) FindAll(ctx context.Context) ([]domain.Education, error) {
	educations, err := r.client.Education.FindMany().
		OrderBy(db.Education.SortOrder.Order(db.ASC), db.Education.StartDate.Order(db.DESC)).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	result := make([]domain.Education, len(educations))
	for i, e := range educations {
		result[i] = *mapEducationToDomain(&e)
	}

	return result, nil
}

func (r *educationRepository) FindByID(ctx context.Context, id string) (*domain.Education, error) {
	education, err := r.client.Education.FindUnique(
		db.Education.ID.Equals(id),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return mapEducationToDomain(education), nil
}

func (r *educationRepository) Create(ctx context.Context, input domain.CreateEducationInput) (*domain.Education, error) {
	education, err := r.client.Education.CreateOne(
		db.Education.School.Set(input.School),
		db.Education.Degree.Set(input.Degree),
		db.Education.StartDate.Set(input.StartDate),
		db.Education.SortOrder.Set(input.SortOrder),
		db.Education.Field.SetIfPresent(input.Field),
		db.Education.Location.SetIfPresent(input.Location),
		db.Education.Description.SetIfPresent(input.Description),
		db.Education.EndDate.SetIfPresent(input.EndDate),
		db.Education.Gpa.SetIfPresent(input.GPA),
		db.Education.LogoURL.SetIfPresent(input.LogoURL),
		db.Education.SchoolURL.SetIfPresent(input.SchoolURL),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return mapEducationToDomain(education), nil
}

func (r *educationRepository) Update(ctx context.Context, id string, input domain.UpdateEducationInput) (*domain.Education, error) {
	updates := []db.EducationSetParam{}

	if input.School != nil {
		updates = append(updates, db.Education.School.Set(*input.School))
	}
	if input.Degree != nil {
		updates = append(updates, db.Education.Degree.Set(*input.Degree))
	}
	if input.Field != nil {
		updates = append(updates, db.Education.Field.Set(*input.Field))
	}
	if input.Location != nil {
		updates = append(updates, db.Education.Location.Set(*input.Location))
	}
	if input.Description != nil {
		updates = append(updates, db.Education.Description.Set(*input.Description))
	}
	if input.StartDate != nil {
		updates = append(updates, db.Education.StartDate.Set(*input.StartDate))
	}
	if input.EndDate != nil {
		updates = append(updates, db.Education.EndDate.Set(*input.EndDate))
	}
	if input.GPA != nil {
		updates = append(updates, db.Education.Gpa.Set(*input.GPA))
	}
	if input.LogoURL != nil {
		updates = append(updates, db.Education.LogoURL.Set(*input.LogoURL))
	}
	if input.SchoolURL != nil {
		updates = append(updates, db.Education.SchoolURL.Set(*input.SchoolURL))
	}
	if input.SortOrder != nil {
		updates = append(updates, db.Education.SortOrder.Set(*input.SortOrder))
	}

	education, err := r.client.Education.FindUnique(
		db.Education.ID.Equals(id),
	).Update(updates...).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return mapEducationToDomain(education), nil
}

func (r *educationRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.Education.FindUnique(
		db.Education.ID.Equals(id),
	).Delete().Exec(ctx)
	return err
}

func mapEducationToDomain(e *db.EducationModel) *domain.Education {
	edu := &domain.Education{
		ID:        e.ID,
		School:    e.School,
		Degree:    e.Degree,
		StartDate: e.StartDate,
		SortOrder: e.SortOrder,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}

	if field, ok := e.Field(); ok {
		edu.Field = &field
	}
	if loc, ok := e.Location(); ok {
		edu.Location = &loc
	}
	if desc, ok := e.Description(); ok {
		edu.Description = &desc
	}
	if endDate, ok := e.EndDate(); ok {
		edu.EndDate = &endDate
	}
	if gpa, ok := e.Gpa(); ok {
		edu.GPA = &gpa
	}
	if logo, ok := e.LogoURL(); ok {
		edu.LogoURL = &logo
	}
	if url, ok := e.SchoolURL(); ok {
		edu.SchoolURL = &url
	}

	return edu
}
