package repository

import (
	"context"
	"time"

	"github.com/ulumfr/ulumfr-be/pkg/models"
	"github.com/ulumfr/ulumfr-be/prisma/db"
)

// Repositories holds all repository instances
type Repositories struct {
	User        UserRepository
	Session     SessionRepository
	Project     ProjectRepository
	Category    CategoryRepository
	Tag         TagRepository
	Career      CareerRepository
	Education   EducationRepository
	Resume      ResumeRepository
	Contact     ContactRepository
	About       AboutRepository
	Blog        BlogRepository
	Certificate CertificateRepository
}

// NewRepositories creates a new Repositories instance
func NewRepositories(client *db.PrismaClient) *Repositories {
	return &Repositories{
		User:        NewUserRepository(client),
		Session:     NewSessionRepository(client),
		Project:     NewProjectRepository(client),
		Category:    NewCategoryRepository(client),
		Tag:         NewTagRepository(client),
		Career:      NewCareerRepository(client),
		Education:   NewEducationRepository(client),
		Resume:      NewResumeRepository(client),
		Contact:     NewContactRepository(client),
		About:       NewAboutRepository(client),
		Blog:        NewBlogRepository(client),
		Certificate: NewCertificateRepository(client),
	}
}

// UserRepository defines user data access methods
type UserRepository interface {
	FindByID(ctx context.Context, id string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindAll(ctx context.Context) ([]models.User, error)
	Create(ctx context.Context, name, email, password string) (*models.User, error)
	Update(ctx context.Context, id string, input models.UpdateProfileInput) (*models.User, error)
}

// SessionRepository defines session data access methods
type SessionRepository interface {
	FindByToken(ctx context.Context, token string) (*models.Session, error)
	Create(ctx context.Context, userID, token string, expires time.Time) (*models.Session, error)
	Delete(ctx context.Context, token string) error
	DeleteByUserID(ctx context.Context, userID string) error
}

// ProjectRepository defines project data access methods
type ProjectRepository interface {
	FindAll(ctx context.Context, params models.ProjectListParams, publishedOnly bool) ([]models.Project, int64, error)
	FindByID(ctx context.Context, id string) (*models.Project, error)
	FindBySlug(ctx context.Context, slug string) (*models.Project, error)
	Create(ctx context.Context, input models.CreateProjectInput) (*models.Project, error)
	Update(ctx context.Context, id string, input models.UpdateProjectInput) (*models.Project, error)
	Delete(ctx context.Context, id string) error
}

// CategoryRepository defines category data access methods
type CategoryRepository interface {
	FindAll(ctx context.Context) ([]models.Category, error)
	FindByID(ctx context.Context, id string) (*models.Category, error)
	Create(ctx context.Context, input models.CreateCategoryInput) (*models.Category, error)
	Update(ctx context.Context, id string, input models.UpdateCategoryInput) (*models.Category, error)
	Delete(ctx context.Context, id string) error
}

// TagRepository defines tag data access methods
type TagRepository interface {
	FindAll(ctx context.Context) ([]models.Tag, error)
	FindByID(ctx context.Context, id string) (*models.Tag, error)
	Create(ctx context.Context, input models.CreateTagInput) (*models.Tag, error)
	Update(ctx context.Context, id string, input models.UpdateTagInput) (*models.Tag, error)
	Delete(ctx context.Context, id string) error
}

// CareerRepository defines career data access methods
type CareerRepository interface {
	FindAll(ctx context.Context) ([]models.Career, error)
	FindByID(ctx context.Context, id string) (*models.Career, error)
	Create(ctx context.Context, input models.CreateCareerInput) (*models.Career, error)
	Update(ctx context.Context, id string, input models.UpdateCareerInput) (*models.Career, error)
	Delete(ctx context.Context, id string) error
}

// EducationRepository defines education data access methods
type EducationRepository interface {
	FindAll(ctx context.Context) ([]models.Education, error)
	FindByID(ctx context.Context, id string) (*models.Education, error)
	Create(ctx context.Context, input models.CreateEducationInput) (*models.Education, error)
	Update(ctx context.Context, id string, input models.UpdateEducationInput) (*models.Education, error)
	Delete(ctx context.Context, id string) error
}

// ResumeRepository defines resume data access methods
type ResumeRepository interface {
	FindAll(ctx context.Context) ([]models.Resume, error)
	FindByID(ctx context.Context, id string) (*models.Resume, error)
	FindActive(ctx context.Context) (*models.Resume, error)
	Create(ctx context.Context, input models.CreateResumeInput) (*models.Resume, error)
	Update(ctx context.Context, id string, input models.UpdateResumeInput) (*models.Resume, error)
	Delete(ctx context.Context, id string) error
	SetActive(ctx context.Context, id string) error
}

// ContactRepository defines contact data access methods
type ContactRepository interface {
	FindAll(ctx context.Context, params models.ContactListParams) ([]models.Contact, int64, error)
	FindByID(ctx context.Context, id string) (*models.Contact, error)
	Create(ctx context.Context, input models.CreateContactInput) (*models.Contact, error)
	MarkAsRead(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}

// AboutRepository defines about data access methods
type AboutRepository interface {
	FindAll(ctx context.Context) ([]models.About, error)
	FindByID(ctx context.Context, id string) (*models.About, error)
	FindActive(ctx context.Context) (*models.About, error)
	Create(ctx context.Context, input models.CreateAboutInput) (*models.About, error)
	Update(ctx context.Context, id string, input models.UpdateAboutInput) (*models.About, error)
	Delete(ctx context.Context, id string) error
}

// BlogRepository defines blog data access methods
type BlogRepository interface {
	FindAll(ctx context.Context, params models.BlogListParams, publishedOnly bool) ([]models.Blog, int64, error)
	FindByID(ctx context.Context, id string) (*models.Blog, error)
	FindBySlug(ctx context.Context, slug string) (*models.Blog, error)
	Create(ctx context.Context, input models.CreateBlogInput) (*models.Blog, error)
	Update(ctx context.Context, id string, input models.UpdateBlogInput) (*models.Blog, error)
	Delete(ctx context.Context, id string) error
}

// CertificateRepository defines certificate data access methods
type CertificateRepository interface {
	FindAll(ctx context.Context) ([]models.Certificate, error)
	FindByID(ctx context.Context, id string) (*models.Certificate, error)
	Create(ctx context.Context, input models.CreateCertificateInput) (*models.Certificate, error)
	Update(ctx context.Context, id string, input models.UpdateCertificateInput) (*models.Certificate, error)
	Delete(ctx context.Context, id string) error
}
