package service

import (
	"github.com/ulumfr/ulumfr-be/pkg/config"
	"github.com/ulumfr/ulumfr-be/pkg/repository"
	"github.com/ulumfr/ulumfr-be/pkg/storage"
)

// Services holds all service instances
type Services struct {
	Auth        *AuthService
	Project     *ProjectService
	Category    *CategoryService
	Tag         *TagService
	Career      *CareerService
	Education   *EducationService
	Resume      *ResumeService
	Contact     *ContactService
	Upload      *UploadService
	About       *AboutService
	Blog        *BlogService
	Certificate *CertificateService
}

// NewServices creates a new Services instance
func NewServices(repos *repository.Repositories, r2Client *storage.R2Client, cfg *config.Config) *Services {
	return &Services{
		Auth:        NewAuthService(repos.User, repos.Session, r2Client, cfg),
		Project:     NewProjectService(repos.Project, r2Client),
		Category:    NewCategoryService(repos.Category),
		Tag:         NewTagService(repos.Tag),
		Career:      NewCareerService(repos.Career, r2Client),
		Education:   NewEducationService(repos.Education, r2Client),
		Resume:      NewResumeService(repos.Resume, r2Client),
		Contact:     NewContactService(repos.Contact),
		Upload:      NewUploadService(r2Client),
		About:       NewAboutService(repos.About, r2Client),
		Blog:        NewBlogService(repos.Blog, r2Client),
		Certificate: NewCertificateService(repos.Certificate, r2Client),
	}
}
