package service

import (
	"github.com/ulumfr/ulumfr-be/pkg/config"
	"github.com/ulumfr/ulumfr-be/pkg/repository"
	"github.com/ulumfr/ulumfr-be/pkg/storage"
)

// Services holds all service instances
type Services struct {
	Auth      *AuthService
	Project   *ProjectService
	Category  *CategoryService
	Tag       *TagService
	Career    *CareerService
	Education *EducationService
	Resume    *ResumeService
	Contact   *ContactService
	Upload    *UploadService
}

// NewServices creates a new Services instance
func NewServices(repos *repository.Repositories, r2Client *storage.R2Client, cfg *config.Config) *Services {
	return &Services{
		Auth:      NewAuthService(repos.User, repos.Session, cfg),
		Project:   NewProjectService(repos.Project),
		Category:  NewCategoryService(repos.Category),
		Tag:       NewTagService(repos.Tag),
		Career:    NewCareerService(repos.Career),
		Education: NewEducationService(repos.Education),
		Resume:    NewResumeService(repos.Resume),
		Contact:   NewContactService(repos.Contact),
		Upload:    NewUploadService(r2Client),
	}
}
