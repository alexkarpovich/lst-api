package services

import (
	"github.com/alexkarpovich/lst-api/src/internal/app"
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
)

type TrainingService interface {
	Create() (*app.Training, error)
	NextItem() (*app.TrainingItem, error)
	ItemAnswers(*valueobject.ID) ([]*app.TrainingAnswer, error)
}
