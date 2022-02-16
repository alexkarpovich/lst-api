package training

import (
	"github.com/alexkarpovich/lst-api/src/internal/app"
)

type trainingDirectService struct {
	*TrainingService
}

func (s *trainingDirectService) Create() (*app.Training, error) {
	training := s.Training
	expressions, err := s.NodeRepo.NativeExpressions(training.Slices)
	if err != nil {
		return nil, err
	}

	xCount := len(expressions)

	for i := 0; i < xCount; i++ {
		trnItem := &app.TrainingItem{
			TrainingId:   training.Id,
			ExpressionId: expressions[i].Id,
			Stage:        1,
			Cycle:        1,
			Complete:     false,
		}

		training.Items = append(training.Items, trnItem)
	}

	meta := &app.TrainingMeta{
		StageCount:      1,
		UniqueItemCount: uint(xCount),
		CompleteCount:   0,
	}

	training.Meta = meta

	return s.TrainingRepo.Create(training)
}
