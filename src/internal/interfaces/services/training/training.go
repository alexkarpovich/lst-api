package training

import (
	"github.com/alexkarpovich/lst-api/src/internal/app"
)

type TrainingService struct {
	NodeRepo app.NodeRepo
}

func (s *TrainingService) Build(trn app.Training) (*app.Training, error) {
	if trn.Type == app.TrainingDirect {
		directService := &trainingDirectService{s}

		return directService.Build(trn)
	}
	if trn.Type == app.TrainingCycles {
		throughService := &trainingCyclesService{s}

		return throughService.Build(trn)
	}

	return nil, nil
}
