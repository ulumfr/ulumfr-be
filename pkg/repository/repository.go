package repository

import (
	"context"
	"time"

	"github.com/ulumfr/ulumfr-be/pkg/domain"
	"github.com/ulumfr/ulumfr-be/prisma/db"
)

// Repositories holds all repository instances
type Repositories struct {
	User      UserRepository
	Session   SessionRepository
	Project   ProjectRepository
	Category  CategoryRepository
	Tag       TagRepository
	Career    CareerRepository
	Education EducationRepository
	Resume    ResumeRepository
	Contact   ContactRepository
}

// NewRepositories creates a new Repositories instance
func NewRepositories(client *db.PrismaClient) *Repositories {
	return &Repositories{
		User:      NewUserRepository(client),
		Session:   NewSessionRepository(client),
		Project:   NewProjectRepository(client),
		Category:  NewCategoryRepository(client),
		Tag:       NewTagRepository(client),
		Career:    NewCareerRepository(client),
		Education: NewEducationRepository(client),
		Resume:    NewResumeRepository(client),
		Contact:   NewContactRepository(client),
	}
}

// UserRepository defines user data access methods
type UserRepository interface {
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Create(ctx context.Context, name, email, password string) (*domain.User, error)
}

// SessionRepository defines session data access methods
type SessionRepository interface {
	FindByToken(ctx context.Context, token string) (*domain.Session, error)
	Create(ctx context.Context, userID, token string, expires time.Time) (*domain.Session, error)
	Delete(ctx context.Context, token string) error
	DeleteByUserID(ctx context.Context, userID string) error
}

// ProjectRepository defines project data access methods
type ProjectRepository interface {
	FindAll(ctx context.Context, params domain.ProjectListParams, publishedOnly bool) ([]domain.Project, int64, error)
	FindByID(ctx context.Context, id string) (*domain.Project, error)
	FindBySlug(ctx context.Context, slug string) (*domain.Project, error)
	Create(ctx context.Context, input domain.CreateProjectInput) (*domain.Project, error)
	Update(ctx context.Context, id string, input domain.UpdateProjectInput) (*domain.Project, error)
	Delete(ctx context.Context, id string) error
}

// CategoryRepository defines category data access methods
type CategoryRepository interface {
	FindAll(ctx context.Context) ([]domain.Category, error)
	FindByID(ctx context.Context, id string) (*domain.Category, error)
	Create(ctx context.Context, input domain.CreateCategoryInput) (*domain.Category, error)
	Update(ctx context.Context, id string, input domain.UpdateCategoryInput) (*domain.Category, error)
	Delete(ctx context.Context, id string) error
}

// TagRepository defines tag data access methods
type TagRepository interface {
	FindAll(ctx context.Context) ([]domain.Tag, error)
	FindByID(ctx context.Context, id string) (*domain.Tag, error)
	Create(ctx context.Context, input domain.CreateTagInput) (*domain.Tag, error)
	Update(ctx context.Context, id string, input domain.UpdateTagInput) (*domain.Tag, error)
	Delete(ctx context.Context, id string) error
}

// CareerRepository defines career data access methods
type CareerRepository interface {
	FindAll(ctx context.Context) ([]domain.Career, error)
	FindByID(ctx context.Context, id string) (*domain.Career, error)
	Create(ctx context.Context, input domain.CreateCareerInput) (*domain.Career, error)
	Update(ctx context.Context, id string, input domain.UpdateCareerInput) (*domain.Career, error)
	Delete(ctx context.Context, id string) error
}

// EducationRepository defines education data access methods
type EducationRepository interface {
	FindAll(ctx context.Context) ([]domain.Education, error)
	FindByID(ctx context.Context, id string) (*domain.Education, error)
	Create(ctx context.Context, input domain.CreateEducationInput) (*domain.Education, error)
	Update(ctx context.Context, id string, input domain.UpdateEducationInput) (*domain.Education, error)
	Delete(ctx context.Context, id string) error
}

// ResumeRepository defines resume data access methods
type ResumeRepository interface {
	FindAll(ctx context.Context) ([]domain.Resume, error)
	FindByID(ctx context.Context, id string) (*domain.Resume, error)
	FindActive(ctx context.Context) (*domain.Resume, error)
	Create(ctx context.Context, input domain.CreateResumeInput) (*domain.Resume, error)
	Update(ctx context.Context, id string, input domain.UpdateResumeInput) (*domain.Resume, error)
	Delete(ctx context.Context, id string) error
	SetActive(ctx context.Context, id string) error
}

// ContactRepository defines contact data access methods
type ContactRepository interface {
	FindAll(ctx context.Context, params domain.ContactListParams) ([]domain.Contact, int64, error)
	FindByID(ctx context.Context, id string) (*domain.Contact, error)
	Create(ctx context.Context, input domain.CreateContactInput) (*domain.Contact, error)
	MarkAsRead(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}
