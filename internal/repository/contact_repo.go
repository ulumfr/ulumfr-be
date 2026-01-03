package repository

import (
	"context"

	"github.com/ulumfr/ulumfr-be/internal/domain"
	"github.com/ulumfr/ulumfr-be/prisma/db"
)

type contactRepository struct {
	client *db.PrismaClient
}

// NewContactRepository creates a new contact repository
func NewContactRepository(client *db.PrismaClient) ContactRepository {
	return &contactRepository{client: client}
}

func (r *contactRepository) FindAll(ctx context.Context, params domain.ContactListParams) ([]domain.Contact, int64, error) {
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
	filters := []db.ContactWhereParam{}

	if params.IsRead != nil {
		filters = append(filters, db.Contact.IsRead.Equals(*params.IsRead))
	}

	// Count total
	count, err := r.client.Contact.FindMany(filters...).Exec(ctx)
	if err != nil {
		return nil, 0, err
	}
	total := int64(len(count))

	// Fetch contacts
	contacts, err := r.client.Contact.FindMany(filters...).
		OrderBy(db.Contact.CreatedAt.Order(db.DESC)).
		Skip(offset).
		Take(params.Limit).
		Exec(ctx)

	if err != nil {
		return nil, 0, err
	}

	result := make([]domain.Contact, len(contacts))
	for i, c := range contacts {
		result[i] = *mapContactToDomain(&c)
	}

	return result, total, nil
}

func (r *contactRepository) FindByID(ctx context.Context, id string) (*domain.Contact, error) {
	contact, err := r.client.Contact.FindUnique(
		db.Contact.ID.Equals(id),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return mapContactToDomain(contact), nil
}

func (r *contactRepository) Create(ctx context.Context, input domain.CreateContactInput) (*domain.Contact, error) {
	contact, err := r.client.Contact.CreateOne(
		db.Contact.Name.Set(input.Name),
		db.Contact.Email.Set(input.Email),
		db.Contact.Message.Set(input.Message),
		db.Contact.Subject.SetIfPresent(input.Subject),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return mapContactToDomain(contact), nil
}

func (r *contactRepository) MarkAsRead(ctx context.Context, id string) error {
	_, err := r.client.Contact.FindUnique(
		db.Contact.ID.Equals(id),
	).Update(
		db.Contact.IsRead.Set(true),
	).Exec(ctx)

	return err
}

func (r *contactRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.Contact.FindUnique(
		db.Contact.ID.Equals(id),
	).Delete().Exec(ctx)
	return err
}

func mapContactToDomain(c *db.ContactModel) *domain.Contact {
	contact := &domain.Contact{
		ID:        c.ID,
		Name:      c.Name,
		Email:     c.Email,
		Message:   c.Message,
		IsRead:    c.IsRead,
		CreatedAt: c.CreatedAt,
	}

	if subject, ok := c.Subject(); ok {
		contact.Subject = &subject
	}

	return contact
}
