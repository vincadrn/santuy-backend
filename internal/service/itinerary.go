package service

import (
	"context"

	"vincadrn.com/santuy/internal/model"
	"vincadrn.com/santuy/internal/repository"
)

type ItineraryService struct {
	repo repository.ItineraryRepository
}

func NewItineraryService(repo repository.ItineraryRepository) *ItineraryService {
	return &ItineraryService{repo: repo}
}

func (svc *ItineraryService) ListItineraries(ctx context.Context, groupId string, vacationId string) (*model.Itineraries, error) {
	return svc.repo.ListItineraries(ctx, groupId, vacationId)
}

func (svc *ItineraryService) ListItineraryDetails(ctx context.Context, groupId string, itineraryId string) (*model.ItineraryDetails, error) {
	return svc.repo.ListItineraryDetails(ctx, groupId, itineraryId)
}
