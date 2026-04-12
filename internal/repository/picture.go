package repository

import (
	"context"
	"database/sql"

	"vincadrn.com/santuy/internal/model"
	"vincadrn.com/santuy/internal/oci"
)

type PictureRepository interface {
	CreateObjectUploadURI(ctx context.Context, itineraryId string, format string) (*model.Picture, error)
	ListPictureByID(ctx context.Context, itineraryId string) (*model.Pictures, error)
}

type pictureRepository struct {
	db *sql.DB
}

func NewPictureRepository(db *sql.DB) PictureRepository {
	return &pictureRepository{db: db}
}

func (r *pictureRepository) CreateObjectUploadURI(ctx context.Context, itineraryId string, format string) (*model.Picture, error) {
	return oci.CreateUploadObjectURI("itineraryId", itineraryId, format)
}

func (r *pictureRepository) ListPictureByID(ctx context.Context, itineraryId string) (*model.Pictures, error) {
	return oci.ListObjectsByID("itineraryId", itineraryId)
}
