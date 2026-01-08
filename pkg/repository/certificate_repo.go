package repository

import (
	"context"

	"github.com/ulumfr/ulumfr-be/pkg/models"
	"github.com/ulumfr/ulumfr-be/prisma/db"
)

type certificateRepository struct {
	client *db.PrismaClient
}

// NewCertificateRepository creates a new certificate repository
func NewCertificateRepository(client *db.PrismaClient) CertificateRepository {
	return &certificateRepository{client: client}
}

func (r *certificateRepository) FindAll(ctx context.Context) ([]models.Certificate, error) {
	certificates, err := r.client.Certificate.FindMany().
		OrderBy(db.Certificate.SortOrder.Order(db.ASC), db.Certificate.IssueDate.Order(db.DESC)).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	result := make([]models.Certificate, len(certificates))
	for i, c := range certificates {
		result[i] = *mapCertificateToDomain(&c)
	}

	return result, nil
}

func (r *certificateRepository) FindByID(ctx context.Context, id string) (*models.Certificate, error) {
	cert, err := r.client.Certificate.FindUnique(
		db.Certificate.ID.Equals(id),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return mapCertificateToDomain(cert), nil
}

func (r *certificateRepository) Create(ctx context.Context, input models.CreateCertificateInput) (*models.Certificate, error) {
	cert, err := r.client.Certificate.CreateOne(
		db.Certificate.Name.Set(input.Name),
		db.Certificate.Issuer.Set(input.Issuer),
		db.Certificate.IssueDate.Set(input.IssueDate),
		db.Certificate.SortOrder.Set(input.SortOrder),
		db.Certificate.ExpiryDate.SetIfPresent(input.ExpiryDate),
		db.Certificate.CredentialID.SetIfPresent(input.CredentialID),
		db.Certificate.CredentialURL.SetIfPresent(input.CredentialURL),
		db.Certificate.ImageURL.SetIfPresent(input.ImageURL),
		db.Certificate.Description.SetIfPresent(input.Description),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return mapCertificateToDomain(cert), nil
}

func (r *certificateRepository) Update(ctx context.Context, id string, input models.UpdateCertificateInput) (*models.Certificate, error) {
	updates := []db.CertificateSetParam{}

	if input.Name != nil {
		updates = append(updates, db.Certificate.Name.Set(*input.Name))
	}
	if input.Issuer != nil {
		updates = append(updates, db.Certificate.Issuer.Set(*input.Issuer))
	}
	if input.IssueDate != nil {
		updates = append(updates, db.Certificate.IssueDate.Set(*input.IssueDate))
	}
	if input.ExpiryDate != nil {
		updates = append(updates, db.Certificate.ExpiryDate.Set(*input.ExpiryDate))
	}
	if input.CredentialID != nil {
		updates = append(updates, db.Certificate.CredentialID.Set(*input.CredentialID))
	}
	if input.CredentialURL != nil {
		updates = append(updates, db.Certificate.CredentialURL.Set(*input.CredentialURL))
	}
	if input.ImageURL != nil {
		updates = append(updates, db.Certificate.ImageURL.Set(*input.ImageURL))
	}
	if input.Description != nil {
		updates = append(updates, db.Certificate.Description.Set(*input.Description))
	}
	if input.SortOrder != nil {
		updates = append(updates, db.Certificate.SortOrder.Set(*input.SortOrder))
	}

	cert, err := r.client.Certificate.FindUnique(
		db.Certificate.ID.Equals(id),
	).Update(updates...).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return mapCertificateToDomain(cert), nil
}

func (r *certificateRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.Certificate.FindUnique(
		db.Certificate.ID.Equals(id),
	).Delete().Exec(ctx)
	return err
}

func mapCertificateToDomain(c *db.CertificateModel) *models.Certificate {
	cert := &models.Certificate{
		ID:        c.ID,
		Name:      c.Name,
		Issuer:    c.Issuer,
		IssueDate: c.IssueDate,
		SortOrder: c.SortOrder,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}

	if expiry, ok := c.ExpiryDate(); ok {
		cert.ExpiryDate = &expiry
	}
	if credID, ok := c.CredentialID(); ok {
		cert.CredentialID = &credID
	}
	if credURL, ok := c.CredentialURL(); ok {
		cert.CredentialURL = &credURL
	}
	if img, ok := c.ImageURL(); ok {
		cert.ImageURL = &img
	}
	if desc, ok := c.Description(); ok {
		cert.Description = &desc
	}

	return cert
}
