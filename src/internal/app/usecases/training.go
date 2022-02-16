package usecases

import (
	"errors"

	"github.com/alexkarpovich/lst-api/src/internal/app"
	"github.com/alexkarpovich/lst-api/src/internal/app/services"
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
	"github.com/alexkarpovich/lst-api/src/internal/interfaces/services/training"
)

type TrainingInteractor struct {
	TrainingRepo    app.TrainingRepo
	NodeRepo        app.NodeRepo
	TrainingService services.TrainingService
}

func NewTrainingInteractor(tr app.TrainingRepo, nr app.NodeRepo, ts services.TrainingService) *TrainingInteractor {
	return &TrainingInteractor{tr, nr, ts}
}

func (i *TrainingInteractor) GetOrCreate(inTraining app.Training) (*app.Training, error) {
	sliceOnlyIds, err := i.NodeRepo.FilterSliceIds(inTraining.Slices)
	if err != nil {
		return nil, err
	}

	if len(sliceOnlyIds) == 0 {
		return nil, errors.New("There must be at least one node.")
	}

	if !i.TrainingRepo.HasCreatePermission(inTraining.OwnerId, sliceOnlyIds) {
		return nil, errors.New("Forbidden, only user which has at least read role can do this.")
	}

	inTraining.Slices = sliceOnlyIds

	if trn, err := i.TrainingRepo.GetBySlices(inTraining); err == nil && trn != nil {
		return trn, nil
	}

	trainingService := training.NewService(i.NodeRepo, i.TrainingRepo, inTraining)
	training, err := trainingService.Create()
	if err != nil {
		return nil, err
	}

	return training, nil
}

func (i *TrainingInteractor) Get(actorId *valueobject.ID, trainingId *valueobject.ID) (*app.Training, error) {
	training, err := i.TrainingRepo.Get(trainingId)
	if err != nil {
		return nil, err
	}

	if *training.OwnerId != *actorId {
		return nil, errors.New("Forbidden, only training owner can do this.")
	}

	return training, nil
}

func (i *TrainingInteractor) List(actorId *valueobject.ID) ([]*app.Training, error) {
	trainings, err := i.TrainingRepo.List(actorId)
	if err != nil {
		return nil, err
	}

	return trainings, nil
}

func (i *TrainingInteractor) Reset(actorId *valueobject.ID, trainingId *valueobject.ID) error {
	training, err := i.TrainingRepo.Get(trainingId)
	if err != nil {
		return err
	}

	if *training.OwnerId != *actorId {
		return errors.New("Forbidden, only training owner can do this.")
	}

	err = i.TrainingRepo.Reset(trainingId)
	if err != nil {
		return err
	}

	return nil
}

func (i *TrainingInteractor) Next(actorId *valueobject.ID, trainingId *valueobject.ID) (*app.TrainingItem, error) {
	trn, err := i.TrainingRepo.Get(trainingId)
	if err != nil {
		return nil, err
	}

	if *trn.OwnerId != *actorId {
		return nil, errors.New("Forbidden, only training owner can do this.")
	}

	trainingService := training.NewService(i.NodeRepo, i.TrainingRepo, *trn)
	trainingItem, err := trainingService.NextItem()
	if err != nil {
		return nil, err
	}

	return trainingItem, nil
}

func (i *TrainingInteractor) GetItem(actorId *valueobject.ID, itemId *valueobject.ID) (*app.TrainingItem, error) {
	training, err := i.TrainingRepo.GetByItemId(itemId)
	if err != nil {
		return nil, err
	}

	if *training.OwnerId != *actorId {
		return nil, errors.New("Forbidden, only training owner can do this.")
	}

	return nil, nil
}

func (i *TrainingInteractor) ItemAnswers(actorId *valueobject.ID, itemId *valueobject.ID) ([]*app.TrainingAnswer, error) {
	trn, err := i.TrainingRepo.GetByItemId(itemId)
	if err != nil {
		return nil, err
	}

	if *trn.OwnerId != *actorId {
		return nil, errors.New("Forbidden, only training owner can do this.")
	}

	trainingService := training.NewService(i.NodeRepo, i.TrainingRepo, *trn)
	answers, err := trainingService.ItemAnswers(itemId)
	if err != nil {
		return nil, err
	}

	return answers, nil
}

func (i *TrainingInteractor) MarkItemAsComplete(actorId *valueobject.ID, itemId *valueobject.ID) error {
	training, err := i.TrainingRepo.GetByItemId(itemId)
	if err != nil {
		return err
	}

	if *training.OwnerId != *actorId {
		return errors.New("Forbidden, only training owner can do this.")
	}

	err = i.TrainingRepo.MarkItemAsComplete(itemId)
	if err != nil {
		return err
	}

	return nil
}
