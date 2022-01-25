package services

import "github.com/alexkarpovich/lst-api/src/internal/app"

type TrainingService interface {
	Build(app.Training) (*app.Training, error)
}
