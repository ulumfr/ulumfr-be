package repository

import (
	"context"

	"github.com/ulumfr/ulumfr-be/pkg/domain"
	"github.com/ulumfr/ulumfr-be/prisma/db"
)

type userRepository struct {
	client *db.PrismaClient
}

// NewUserRepository creates a new user repository
func NewUserRepository(client *db.PrismaClient) UserRepository {
	return &userRepository{client: client}
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := r.client.User.FindUnique(
		db.User.ID.Equals(id),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return mapUserToDomain(user), nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := r.client.User.FindUnique(
		db.User.Email.Equals(email),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return mapUserToDomain(user), nil
}

func (r *userRepository) Create(ctx context.Context, name, email, password string) (*domain.User, error) {
	user, err := r.client.User.CreateOne(
		db.User.Email.Set(email),
		db.User.Name.Set(name),
		db.User.Password.Set(password),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return mapUserToDomain(user), nil
}

func mapUserToDomain(u *db.UserModel) *domain.User {
	user := &domain.User{
		ID:        u.ID,
		Email:     u.Email,
		Role:      string(u.Role),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}

	if name, ok := u.Name(); ok {
		user.Name = &name
	}
	if image, ok := u.Image(); ok {
		user.Image = &image
	}
	if emailVerified, ok := u.EmailVerified(); ok {
		user.EmailVerified = &emailVerified
	}
	if password, ok := u.Password(); ok {
		user.Password = &password
	}

	return user
}
