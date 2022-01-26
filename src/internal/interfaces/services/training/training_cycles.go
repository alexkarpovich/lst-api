package training

import (
	"math"
	"math/rand"
	"time"

	"github.com/alexkarpovich/lst-api/src/internal/app"
)

const (
	minChunkSize = 7
)

type trainingCyclesService struct {
	*TrainingService
}

func (s *trainingCyclesService) Build(trn app.Training) (*app.Training, error) {
	training := &trn
	expressions, err := s.NodeRepo.NativeExpressions(trn.Slices)
	if err != nil {
		return nil, err
	}

	var stage uint
	xCount := len(expressions)
	stageCount := int(math.Round(math.Log(float64(xCount)/minChunkSize) / math.Log(2)))

	for stage = 1; stage <= uint(stageCount); stage++ {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(expressions), func(i, j int) { expressions[i], expressions[j] = expressions[j], expressions[i] })

		rate := math.Round(float64(xCount) / (minChunkSize * math.Pow(2, float64(stage))))
		cycleCount := int(math.Round(float64(xCount) / rate))
		chunkSize := (xCount + cycleCount - 1) / cycleCount

		for i := 0; i < xCount; i++ {
			cycle := uint(math.Round(float64(i) / float64(chunkSize)))

			trnItem := &app.TrainingItem{
				TrainingId:   trn.Id,
				ExpressionId: expressions[i].Id,
				Stage:        stage,
				Cycle:        cycle,
				Complete:     false,
			}

			training.Items = append(training.Items, trnItem)
		}
	}

	return training, nil
}
