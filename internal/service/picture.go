package service

import (
	"context"

	"vincadrn.com/santuy/internal/model"
	"vincadrn.com/santuy/internal/repository"
)

type PictureService struct {
	pictRepo repository.PictureRepository
}

func NewPictureService(pictRepo repository.PictureRepository, itineraryRepo repository.ItineraryRepository) *PictureService {
	return &PictureService{
		pictRepo: pictRepo,
	}
}

func (svc *PictureService) CreateObjectUploadURI(ctx context.Context, itineraryId string, format string) (*model.Picture, error) {
	return svc.pictRepo.CreateObjectUploadURI(ctx, itineraryId, format)
}

func (svc *PictureService) ListPicturesByItinerary(ctx context.Context, itineraryId string) (*model.Pictures, error) {
	return svc.pictRepo.ListPictureByID(ctx, itineraryId)
}
