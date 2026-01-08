package repository

import (
	"context"
	"time"

	"github.com/ulumfr/ulumfr-be/pkg/models"
	"github.com/ulumfr/ulumfr-be/prisma/db"
)

type sessionRepository struct {
	client *db.PrismaClient
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(client *db.PrismaClient) SessionRepository {
	return &sessionRepository{client: client}
}

func (r *sessionRepository) FindByToken(ctx context.Context, token string) (*models.Session, error) {
	session, err := r.client.Session.FindUnique(
		db.Session.SessionToken.Equals(token),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &models.Session{
		ID:           session.ID,
		SessionToken: session.SessionToken,
		UserID:       session.UserID,
		Expires:      session.Expires,
	}, nil
}

func (r *sessionRepository) Create(ctx context.Context, userID, token string, expires time.Time) (*models.Session, error) {
	session, err := r.client.Session.CreateOne(
		db.Session.SessionToken.Set(token),
		db.Session.Expires.Set(expires),
		db.Session.User.Link(db.User.ID.Equals(userID)),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &models.Session{
		ID:           session.ID,
		SessionToken: session.SessionToken,
		UserID:       session.UserID,
		Expires:      session.Expires,
	}, nil
}

func (r *sessionRepository) Delete(ctx context.Context, token string) error {
	_, err := r.client.Session.FindUnique(
		db.Session.SessionToken.Equals(token),
	).Delete().Exec(ctx)

	return err
}

func (r *sessionRepository) DeleteByUserID(ctx context.Context, userID string) error {
	_, err := r.client.Session.FindMany(
		db.Session.UserID.Equals(userID),
	).Delete().Exec(ctx)

	return err
}
