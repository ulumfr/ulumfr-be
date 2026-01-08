package repository

import (
	"context"
	"strings"

	"github.com/ulumfr/ulumfr-be/pkg/models"
	"github.com/ulumfr/ulumfr-be/prisma/db"
)

type blogRepository struct {
	client *db.PrismaClient
}

// NewBlogRepository creates a new blog repository
func NewBlogRepository(client *db.PrismaClient) BlogRepository {
	return &blogRepository{client: client}
}

func (r *blogRepository) FindAll(ctx context.Context, params models.BlogListParams, publishedOnly bool) ([]models.Blog, int64, error) {
	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 10
	}

	// Build filters
	filters := []db.BlogWhereParam{}

	if publishedOnly {
		filters = append(filters, db.Blog.IsPublished.Equals(true))
	}

	if params.TagID != "" {
		filters = append(filters, db.Blog.Tags.Some(
			db.BlogTag.TagID.Equals(params.TagID),
		))
	}

	if params.IsFeatured != nil {
		filters = append(filters, db.Blog.IsFeatured.Equals(*params.IsFeatured))
	}

	if params.Search != "" {
		search := strings.ToLower(params.Search)
		filters = append(filters, db.Blog.Or(
			db.Blog.Title.Contains(search),
			db.Blog.Excerpt.Contains(search),
		))
	}

	// Get total count
	count, err := r.client.Blog.FindMany(filters...).Exec(ctx)
	if err != nil {
		return nil, 0, err
	}
	total := int64(len(count))

	// Get paginated results
	offset := (params.Page - 1) * params.Limit
	blogs, err := r.client.Blog.FindMany(filters...).
		With(db.Blog.Tags.Fetch().With(db.BlogTag.Tag.Fetch())).
		OrderBy(db.Blog.SortOrder.Order(db.ASC), db.Blog.CreatedAt.Order(db.DESC)).
		Skip(offset).
		Take(params.Limit).
		Exec(ctx)

	if err != nil {
		return nil, 0, err
	}

	result := make([]models.Blog, len(blogs))
	for i, b := range blogs {
		result[i] = *mapBlogToDomain(&b)
	}

	return result, total, nil
}

func (r *blogRepository) FindByID(ctx context.Context, id string) (*models.Blog, error) {
	blog, err := r.client.Blog.FindUnique(
		db.Blog.ID.Equals(id),
	).With(db.Blog.Tags.Fetch().With(db.BlogTag.Tag.Fetch())).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return mapBlogToDomain(blog), nil
}

func (r *blogRepository) FindBySlug(ctx context.Context, slug string) (*models.Blog, error) {
	blog, err := r.client.Blog.FindUnique(
		db.Blog.Slug.Equals(slug),
	).With(db.Blog.Tags.Fetch().With(db.BlogTag.Tag.Fetch())).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return mapBlogToDomain(blog), nil
}

func (r *blogRepository) Create(ctx context.Context, input models.CreateBlogInput) (*models.Blog, error) {
	blog, err := r.client.Blog.CreateOne(
		db.Blog.Title.Set(input.Title),
		db.Blog.Slug.Set(input.Slug),
		db.Blog.IsPublished.Set(input.IsPublished),
		db.Blog.IsFeatured.Set(input.IsFeatured),
		db.Blog.SortOrder.Set(input.SortOrder),
		db.Blog.Excerpt.SetIfPresent(input.Excerpt),
		db.Blog.Content.SetIfPresent(input.Content),
		db.Blog.CoverImage.SetIfPresent(input.CoverImage),
		db.Blog.PublishedAt.SetIfPresent(input.PublishedAt),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	// Link tags
	if len(input.TagIDs) > 0 {
		for _, tagID := range input.TagIDs {
			_, err := r.client.BlogTag.CreateOne(
				db.BlogTag.Blog.Link(db.Blog.ID.Equals(blog.ID)),
				db.BlogTag.Tag.Link(db.Tag.ID.Equals(tagID)),
			).Exec(ctx)
			if err != nil {
				// Log but continue - tag might not exist
				continue
			}
		}
	}

	return r.FindByID(ctx, blog.ID)
}

func (r *blogRepository) Update(ctx context.Context, id string, input models.UpdateBlogInput) (*models.Blog, error) {
	updates := []db.BlogSetParam{}

	if input.Title != nil {
		updates = append(updates, db.Blog.Title.Set(*input.Title))
	}
	if input.Slug != nil {
		updates = append(updates, db.Blog.Slug.Set(*input.Slug))
	}
	if input.Excerpt != nil {
		updates = append(updates, db.Blog.Excerpt.Set(*input.Excerpt))
	}
	if input.Content != nil {
		updates = append(updates, db.Blog.Content.Set(*input.Content))
	}
	if input.CoverImage != nil {
		updates = append(updates, db.Blog.CoverImage.Set(*input.CoverImage))
	}
	if input.IsPublished != nil {
		updates = append(updates, db.Blog.IsPublished.Set(*input.IsPublished))
	}
	if input.IsFeatured != nil {
		updates = append(updates, db.Blog.IsFeatured.Set(*input.IsFeatured))
	}
	if input.PublishedAt != nil {
		updates = append(updates, db.Blog.PublishedAt.Set(*input.PublishedAt))
	}
	if input.SortOrder != nil {
		updates = append(updates, db.Blog.SortOrder.Set(*input.SortOrder))
	}

	_, err := r.client.Blog.FindUnique(
		db.Blog.ID.Equals(id),
	).Update(updates...).Exec(ctx)

	if err != nil {
		return nil, err
	}

	// Update tags if provided
	if input.TagIDs != nil {
		// Remove existing tags
		r.client.BlogTag.FindMany(
			db.BlogTag.BlogID.Equals(id),
		).Delete().Exec(ctx)

		// Add new tags
		for _, tagID := range input.TagIDs {
			r.client.BlogTag.CreateOne(
				db.BlogTag.Blog.Link(db.Blog.ID.Equals(id)),
				db.BlogTag.Tag.Link(db.Tag.ID.Equals(tagID)),
			).Exec(ctx)
		}
	}

	return r.FindByID(ctx, id)
}

func (r *blogRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.Blog.FindUnique(
		db.Blog.ID.Equals(id),
	).Delete().Exec(ctx)
	return err
}

func mapBlogToDomain(b *db.BlogModel) *models.Blog {
	blog := &models.Blog{
		ID:          b.ID,
		Title:       b.Title,
		Slug:        b.Slug,
		IsPublished: b.IsPublished,
		IsFeatured:  b.IsFeatured,
		SortOrder:   b.SortOrder,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}

	if excerpt, ok := b.Excerpt(); ok {
		blog.Excerpt = &excerpt
	}
	if content, ok := b.Content(); ok {
		blog.Content = &content
	}
	if cover, ok := b.CoverImage(); ok {
		blog.CoverImage = &cover
	}
	if publishedAt, ok := b.PublishedAt(); ok {
		blog.PublishedAt = &publishedAt
	}

	// Map tags
	if tags := b.Tags(); tags != nil {
		blog.Tags = make([]models.Tag, 0, len(tags))
		for _, bt := range tags {
			if tag := bt.Tag(); tag != nil {
				blog.Tags = append(blog.Tags, models.Tag{
					ID:        tag.ID,
					Name:      tag.Name,
					Slug:      tag.Slug,
					CreatedAt: tag.CreatedAt,
					UpdatedAt: tag.UpdatedAt,
				})
			}
		}
	}

	return blog
}
