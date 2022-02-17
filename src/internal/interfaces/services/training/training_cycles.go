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

func (s *trainingCyclesService) Create() (*app.Training, error) {
	training := s.Training
	translations, err := s.NodeRepo.TranslationsBySlices(training.Slices)
	if err != nil {
		return nil, err
	}

	var stage uint
	xCount := len(translations)
	stageCount := uint(math.Round(math.Log(float64(xCount)/minChunkSize)/math.Log(2))) + 1

	for stage = 1; stage <= uint(stageCount); stage++ {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(translations), func(i, j int) { translations[i], translations[j] = translations[j], translations[i] })

		rate := math.Round(float64(xCount) / (minChunkSize * math.Pow(2, float64(stage))))
		chunkSize := (float64(xCount) + rate - 1) / rate

		for i := 0; i < xCount; i++ {
			cycle := uint(math.Round(float64(i) / float64(chunkSize)))

			trnItem := &app.TrainingItem{
				TrainingId:    s.Training.Id,
				TranslationId: translations[i].Id,
				Stage:         stage,
				Cycle:         cycle,
				Complete:      false,
			}

			training.Items = append(training.Items, trnItem)
		}
	}

	meta := &app.TrainingMeta{
		StageCount:      stageCount,
		UniqueItemCount: uint(xCount),
		CompleteCount:   0,
	}

	training.Meta = meta

	return s.TrainingRepo.Create(training)
}
