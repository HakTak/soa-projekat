package service

import (
	"context"
	"time"
	"tour/internal/model"
	"tour/internal/repository"

	"github.com/google/uuid"
)

type TourService struct {
	repo repository.TourRepository
}

func NewTourService(repo repository.TourRepository) *TourService {
	return &TourService{repo}
}

// Create nova tura
func (s *TourService) Create(ctx context.Context, authorID, name, description string, dificulty model.DificultyLevel, status model.Status, tags []model.Tag) (*model.Tour, error) {

	var newTags []string
	println(status)
	var newstatus = model.StatusDraft
	for _, tag := range tags {
		switch string(tag) {
		case "cultural":
			newTags = append(newTags, string(model.Cultural))
		case "food":
			newTags = append(newTags, string(model.Food))
		case "attraction":
			newTags = append(newTags, string(model.Attraction))
		case "beach":
			newTags = append(newTags, string(model.Beach))
		case "mountain":
			newTags = append(newTags, string(model.Mountain))
		case "nature":
			newTags = append(newTags, string(model.Nature))
		case "adventure":
			newTags = append(newTags, string(model.Adventure))
		case "relax":
			newTags = append(newTags, string(model.Relax))
		case "nightlife":
			newTags = append(newTags, string(model.Nightlife))
		case "shopping":
			newTags = append(newTags, string(model.Shopping))
		default:
			// ignoris nepoznate tagove
		}
	}

	tour := &model.Tour{
		ID:          uuid.NewString(),
		IDAuthor:    authorID,
		Name:        name,
		Description: description,
		Price:       0,
		Dificulty:   dificulty,
		Status:      newstatus,
		Tags:        newTags,
		PublishedAt: time.Now(), // ili time.Now(), zavisi šta koristiš
	}

	println("AuthorID: ", authorID)

	err := s.repo.Create(ctx, tour)
	if err != nil {
		return nil, err
	}

	return tour, nil
}

// GetAllByAuthorID vraća sve ture jednog autora
func (s *TourService) GetAllByAuthorID(ctx context.Context, authorID string) ([]*model.Tour, error) {
	return s.repo.GetByAuthor(ctx, authorID)
}
