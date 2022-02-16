package training

import (
	"github.com/alexkarpovich/lst-api/src/internal/app"
	"github.com/alexkarpovich/lst-api/src/internal/app/services"
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
)

type TrainingService struct {
	Training     app.Training
	NodeRepo     app.NodeRepo
	TrainingRepo app.TrainingRepo
}

func NewService(nr app.NodeRepo, tr app.TrainingRepo, trn app.Training) services.TrainingService {
	if trn.Type == app.TrainingDirect {
		return &trainingDirectService{&TrainingService{trn, nr, tr}}
	}
	if trn.Type == app.TrainingCycles {
		return &trainingCyclesService{&TrainingService{trn, nr, tr}}
	}

	return nil
}

func (s *TrainingService) NextItem() (*app.TrainingItem, error) {
	return s.TrainingRepo.NextItem(s.Training.Id)
}

func (s *TrainingService) ItemAnswers(itemId *valueobject.ID) ([]*app.TrainingAnswer, error) {
	return s.TrainingRepo.ItemAnswers(itemId)
}
