package repository

import (
	"context"

	"github.com/ulumfr/ulumfr-be/internal/domain"
	"github.com/ulumfr/ulumfr-be/prisma/db"
)

type projectRepository struct {
	client *db.PrismaClient
}

// NewProjectRepository creates a new project repository
func NewProjectRepository(client *db.PrismaClient) ProjectRepository {
	return &projectRepository{client: client}
}

func (r *projectRepository) FindAll(ctx context.Context, params domain.ProjectListParams, publishedOnly bool) ([]domain.Project, int64, error) {
	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	offset := (params.Page - 1) * params.Limit

	// Build query filters
	filters := []db.ProjectWhereParam{}

	if publishedOnly {
		filters = append(filters, db.Project.IsPublished.Equals(true))
	}

	if params.IsFeatured != nil {
		filters = append(filters, db.Project.IsFeatured.Equals(*params.IsFeatured))
	}

	if params.Search != "" {
		filters = append(filters, db.Project.Or(
			db.Project.Title.Contains(params.Search),
			db.Project.Description.Contains(params.Search),
		))
	}

	// Count total
	count, err := r.client.Project.FindMany(filters...).Exec(ctx)
	if err != nil {
		return nil, 0, err
	}
	total := int64(len(count))

	// Fetch projects
	projects, err := r.client.Project.FindMany(filters...).
		With(
			db.Project.Categories.Fetch().With(db.ProjectCategory.Category.Fetch()),
			db.Project.Tags.Fetch().With(db.ProjectTag.Tag.Fetch()),
			db.Project.Images.Fetch(),
		).
		OrderBy(db.Project.SortOrder.Order(db.ASC), db.Project.CreatedAt.Order(db.DESC)).
		Skip(offset).
		Take(params.Limit).
		Exec(ctx)

	if err != nil {
		return nil, 0, err
	}

	result := make([]domain.Project, len(projects))
	for i, p := range projects {
		result[i] = *mapProjectToDomain(&p)
	}

	return result, total, nil
}

func (r *projectRepository) FindByID(ctx context.Context, id string) (*domain.Project, error) {
	project, err := r.client.Project.FindUnique(
		db.Project.ID.Equals(id),
	).With(
		db.Project.Categories.Fetch().With(db.ProjectCategory.Category.Fetch()),
		db.Project.Tags.Fetch().With(db.ProjectTag.Tag.Fetch()),
		db.Project.Images.Fetch(),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return mapProjectToDomain(project), nil
}

func (r *projectRepository) FindBySlug(ctx context.Context, slug string) (*domain.Project, error) {
	project, err := r.client.Project.FindUnique(
		db.Project.Slug.Equals(slug),
	).With(
		db.Project.Categories.Fetch().With(db.ProjectCategory.Category.Fetch()),
		db.Project.Tags.Fetch().With(db.ProjectTag.Tag.Fetch()),
		db.Project.Images.Fetch(),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return mapProjectToDomain(project), nil
}

func (r *projectRepository) Create(ctx context.Context, input domain.CreateProjectInput) (*domain.Project, error) {
	// Create project
	project, err := r.client.Project.CreateOne(
		db.Project.Title.Set(input.Title),
		db.Project.Slug.Set(input.Slug),
		db.Project.IsPublished.Set(input.IsPublished),
		db.Project.IsFeatured.Set(input.IsFeatured),
		db.Project.SortOrder.Set(input.SortOrder),
		db.Project.Description.SetIfPresent(input.Description),
		db.Project.Content.SetIfPresent(input.Content),
		db.Project.ThumbnailURL.SetIfPresent(input.ThumbnailURL),
		db.Project.DemoURL.SetIfPresent(input.DemoURL),
		db.Project.RepoURL.SetIfPresent(input.RepoURL),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	// Add categories
	for _, catID := range input.CategoryIDs {
		_, err = r.client.ProjectCategory.CreateOne(
			db.ProjectCategory.Project.Link(db.Project.ID.Equals(project.ID)),
			db.ProjectCategory.Category.Link(db.Category.ID.Equals(catID)),
		).Exec(ctx)
		if err != nil {
			return nil, err
		}
	}

	// Add tags
	for _, tagID := range input.TagIDs {
		_, err = r.client.ProjectTag.CreateOne(
			db.ProjectTag.Project.Link(db.Project.ID.Equals(project.ID)),
			db.ProjectTag.Tag.Link(db.Tag.ID.Equals(tagID)),
		).Exec(ctx)
		if err != nil {
			return nil, err
		}
	}

	return r.FindByID(ctx, project.ID)
}

func (r *projectRepository) Update(ctx context.Context, id string, input domain.UpdateProjectInput) (*domain.Project, error) {
	updates := []db.ProjectSetParam{}

	if input.Title != nil {
		updates = append(updates, db.Project.Title.Set(*input.Title))
	}
	if input.Slug != nil {
		updates = append(updates, db.Project.Slug.Set(*input.Slug))
	}
	if input.Description != nil {
		updates = append(updates, db.Project.Description.Set(*input.Description))
	}
	if input.Content != nil {
		updates = append(updates, db.Project.Content.Set(*input.Content))
	}
	if input.ThumbnailURL != nil {
		updates = append(updates, db.Project.ThumbnailURL.Set(*input.ThumbnailURL))
	}
	if input.DemoURL != nil {
		updates = append(updates, db.Project.DemoURL.Set(*input.DemoURL))
	}
	if input.RepoURL != nil {
		updates = append(updates, db.Project.RepoURL.Set(*input.RepoURL))
	}
	if input.IsPublished != nil {
		updates = append(updates, db.Project.IsPublished.Set(*input.IsPublished))
	}
	if input.IsFeatured != nil {
		updates = append(updates, db.Project.IsFeatured.Set(*input.IsFeatured))
	}
	if input.SortOrder != nil {
		updates = append(updates, db.Project.SortOrder.Set(*input.SortOrder))
	}

	if len(updates) > 0 {
		_, err := r.client.Project.FindUnique(
			db.Project.ID.Equals(id),
		).Update(updates...).Exec(ctx)
		if err != nil {
			return nil, err
		}
	}

	// Update categories if provided
	if input.CategoryIDs != nil {
		// Delete existing
		_, err := r.client.ProjectCategory.FindMany(
			db.ProjectCategory.ProjectID.Equals(id),
		).Delete().Exec(ctx)
		if err != nil {
			return nil, err
		}

		// Add new
		for _, catID := range input.CategoryIDs {
			_, err = r.client.ProjectCategory.CreateOne(
				db.ProjectCategory.Project.Link(db.Project.ID.Equals(id)),
				db.ProjectCategory.Category.Link(db.Category.ID.Equals(catID)),
			).Exec(ctx)
			if err != nil {
				return nil, err
			}
		}
	}

	// Update tags if provided
	if input.TagIDs != nil {
		// Delete existing
		_, err := r.client.ProjectTag.FindMany(
			db.ProjectTag.ProjectID.Equals(id),
		).Delete().Exec(ctx)
		if err != nil {
			return nil, err
		}

		// Add new
		for _, tagID := range input.TagIDs {
			_, err = r.client.ProjectTag.CreateOne(
				db.ProjectTag.Project.Link(db.Project.ID.Equals(id)),
				db.ProjectTag.Tag.Link(db.Tag.ID.Equals(tagID)),
			).Exec(ctx)
			if err != nil {
				return nil, err
			}
		}
	}

	return r.FindByID(ctx, id)
}

func (r *projectRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.Project.FindUnique(
		db.Project.ID.Equals(id),
	).Delete().Exec(ctx)
	return err
}

func mapProjectToDomain(p *db.ProjectModel) *domain.Project {
	project := &domain.Project{
		ID:          p.ID,
		Title:       p.Title,
		Slug:        p.Slug,
		IsPublished: p.IsPublished,
		IsFeatured:  p.IsFeatured,
		SortOrder:   p.SortOrder,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}

	if desc, ok := p.Description(); ok {
		project.Description = &desc
	}
	if content, ok := p.Content(); ok {
		project.Content = &content
	}
	if thumb, ok := p.ThumbnailURL(); ok {
		project.ThumbnailURL = &thumb
	}
	if demo, ok := p.DemoURL(); ok {
		project.DemoURL = &demo
	}
	if repo, ok := p.RepoURL(); ok {
		project.RepoURL = &repo
	}

	// Map categories
	if cats := p.Categories(); cats != nil {
		project.Categories = make([]domain.Category, len(cats))
		for i, pc := range cats {
			cat := pc.Category()
			project.Categories[i] = domain.Category{
				ID:        cat.ID,
				Name:      cat.Name,
				Slug:      cat.Slug,
				CreatedAt: cat.CreatedAt,
				UpdatedAt: cat.UpdatedAt,
			}
			if desc, ok := cat.Description(); ok {
				project.Categories[i].Description = &desc
			}
		}
	}

	// Map tags
	if tags := p.Tags(); tags != nil {
		project.Tags = make([]domain.Tag, len(tags))
		for i, pt := range tags {
			t := pt.Tag()
			project.Tags[i] = domain.Tag{
				ID:        t.ID,
				Name:      t.Name,
				Slug:      t.Slug,
				CreatedAt: t.CreatedAt,
				UpdatedAt: t.UpdatedAt,
			}
		}
	}

	// Map images
	if imgs := p.Images(); imgs != nil {
		project.Images = make([]domain.ProjectImage, len(imgs))
		for i, img := range imgs {
			project.Images[i] = domain.ProjectImage{
				ID:        img.ID,
				ProjectID: img.ProjectID,
				URL:       img.URL,
				SortOrder: img.SortOrder,
				CreatedAt: img.CreatedAt,
			}
			if alt, ok := img.Alt(); ok {
				project.Images[i].Alt = &alt
			}
		}
	}

	return project
}
