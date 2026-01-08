package repository

import (
	"context"

	"github.com/ulumfr/ulumfr-be/pkg/models"
	"github.com/ulumfr/ulumfr-be/prisma/db"
)

type userRepository struct {
	client *db.PrismaClient
}

// NewUserRepository creates a new user repository
func NewUserRepository(client *db.PrismaClient) UserRepository {
	return &userRepository{client: client}
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	user, err := r.client.User.FindUnique(
		db.User.ID.Equals(id),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return mapUserToDomain(user), nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := r.client.User.FindUnique(
		db.User.Email.Equals(email),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return mapUserToDomain(user), nil
}

func (r *userRepository) FindAll(ctx context.Context) ([]models.User, error) {
	users, err := r.client.User.FindMany().
		OrderBy(db.User.CreatedAt.Order(db.DESC)).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	result := make([]models.User, len(users))
	for i, u := range users {
		result[i] = *mapUserToDomain(&u)
	}

	return result, nil
}

func (r *userRepository) Create(ctx context.Context, name, email, password string) (*models.User, error) {
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

func (r *userRepository) Update(ctx context.Context, id string, input models.UpdateProfileInput) (*models.User, error) {
	updates := []db.UserSetParam{}

	if input.Name != nil {
		updates = append(updates, db.User.Name.Set(*input.Name))
	}
	if input.Email != nil {
		updates = append(updates, db.User.Email.Set(*input.Email))
	}
	if input.NewPassword != nil {
		// Password should be hashed before calling this function
		updates = append(updates, db.User.Password.Set(*input.NewPassword))
	}
	if input.Image != nil {
		updates = append(updates, db.User.Image.Set(*input.Image))
	}

	if len(updates) == 0 {
		// No updates to make, return the current user
		return r.FindByID(ctx, id)
	}

	user, err := r.client.User.FindUnique(
		db.User.ID.Equals(id),
	).Update(updates...).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return mapUserToDomain(user), nil
}

func mapUserToDomain(u *db.UserModel) *models.User {
	user := &models.User{
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
