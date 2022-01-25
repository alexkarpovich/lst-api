package training

import (
	"github.com/alexkarpovich/lst-api/src/internal/app"
)

type trainingThroughService struct {
	*TrainingService
}

func (s *trainingThroughService) Build(trn app.Training) (*app.Training, error) {
	training := &trn
	expressions, err := s.NodeRepo.NativeExpressions(trn.Slices)
	if err != nil {
		return nil, err
	}

	xCount := len(expressions)

	for i := 0; i < xCount; i++ {
		trnItem := &app.TrainingItem{
			TrainingId:   trn.Id,
			ExpressionId: expressions[i].Id,
			Stage:        0,
			Cycle:        0,
			Complete:     false,
		}

		training.Items = append(training.Items, trnItem)
	}

	return training, nil
}
