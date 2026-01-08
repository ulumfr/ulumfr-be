package repository

import (
	"context"

	"github.com/ulumfr/ulumfr-be/pkg/models"
	"github.com/ulumfr/ulumfr-be/prisma/db"
)

type resumeRepository struct {
	client *db.PrismaClient
}

// NewResumeRepository creates a new resume repository
func NewResumeRepository(client *db.PrismaClient) ResumeRepository {
	return &resumeRepository{client: client}
}

func (r *resumeRepository) FindAll(ctx context.Context) ([]models.Resume, error) {
	resumes, err := r.client.Resume.FindMany().
		OrderBy(db.Resume.CreatedAt.Order(db.DESC)).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	result := make([]models.Resume, len(resumes))
	for i, res := range resumes {
		result[i] = *mapResumeToDomain(&res)
	}

	return result, nil
}

func (r *resumeRepository) FindByID(ctx context.Context, id string) (*models.Resume, error) {
	resume, err := r.client.Resume.FindUnique(
		db.Resume.ID.Equals(id),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return mapResumeToDomain(resume), nil
}

func (r *resumeRepository) FindActive(ctx context.Context) (*models.Resume, error) {
	resume, err := r.client.Resume.FindFirst(
		db.Resume.IsActive.Equals(true),
	).OrderBy(db.Resume.CreatedAt.Order(db.DESC)).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return mapResumeToDomain(resume), nil
}

func (r *resumeRepository) Create(ctx context.Context, input models.CreateResumeInput) (*models.Resume, error) {
	// If this resume is being set as active, deactivate all others
	if input.IsActive {
		_, err := r.client.Resume.FindMany(
			db.Resume.IsActive.Equals(true),
		).Update(
			db.Resume.IsActive.Set(false),
		).Exec(ctx)
		if err != nil {
			return nil, err
		}
	}

	resume, err := r.client.Resume.CreateOne(
		db.Resume.FileURL.Set(input.FileURL),
		db.Resume.FileName.Set(input.FileName),
		db.Resume.IsActive.Set(input.IsActive),
		db.Resume.FileSize.SetIfPresent(input.FileSize),
		db.Resume.Version.SetIfPresent(input.Version),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return mapResumeToDomain(resume), nil
}

func (r *resumeRepository) Update(ctx context.Context, id string, input models.UpdateResumeInput) (*models.Resume, error) {
	updates := []db.ResumeSetParam{}

	if input.FileURL != nil {
		updates = append(updates, db.Resume.FileURL.Set(*input.FileURL))
	}
	if input.FileName != nil {
		updates = append(updates, db.Resume.FileName.Set(*input.FileName))
	}
	if input.FileSize != nil {
		updates = append(updates, db.Resume.FileSize.Set(*input.FileSize))
	}
	if input.Version != nil {
		updates = append(updates, db.Resume.Version.Set(*input.Version))
	}
	if input.IsActive != nil {
		// If setting this as active, deactivate others
		if *input.IsActive {
			_, err := r.client.Resume.FindMany(
				db.Resume.IsActive.Equals(true),
				db.Resume.ID.Not(id),
			).Update(
				db.Resume.IsActive.Set(false),
			).Exec(ctx)
			if err != nil {
				return nil, err
			}
		}
		updates = append(updates, db.Resume.IsActive.Set(*input.IsActive))
	}

	resume, err := r.client.Resume.FindUnique(
		db.Resume.ID.Equals(id),
	).Update(updates...).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return mapResumeToDomain(resume), nil
}

func (r *resumeRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.Resume.FindUnique(
		db.Resume.ID.Equals(id),
	).Delete().Exec(ctx)
	return err
}

func (r *resumeRepository) SetActive(ctx context.Context, id string) error {
	// Deactivate all resumes
	_, err := r.client.Resume.FindMany(
		db.Resume.IsActive.Equals(true),
	).Update(
		db.Resume.IsActive.Set(false),
	).Exec(ctx)
	if err != nil {
		return err
	}

	// Activate the specified resume
	_, err = r.client.Resume.FindUnique(
		db.Resume.ID.Equals(id),
	).Update(
		db.Resume.IsActive.Set(true),
	).Exec(ctx)

	return err
}

func mapResumeToDomain(r *db.ResumeModel) *models.Resume {
	resume := &models.Resume{
		ID:        r.ID,
		FileURL:   r.FileURL,
		FileName:  r.FileName,
		IsActive:  r.IsActive,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}

	if size, ok := r.FileSize(); ok {
		resume.FileSize = &size
	}
	if version, ok := r.Version(); ok {
		resume.Version = &version
	}

	return resume
}
