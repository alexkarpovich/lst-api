package usecases

import (
	"errors"

	"github.com/alexkarpovich/lst-api/src/internal/app"
	"github.com/alexkarpovich/lst-api/src/internal/app/services"
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
)

type TrainingInteractor struct {
	TrainingRepo    app.TrainingRepo
	NodeRepo        app.NodeRepo
	TrainingService services.TrainingService
}

func NewTrainingInteractor(tr app.TrainingRepo, nr app.NodeRepo, ts services.TrainingService) *TrainingInteractor {
	return &TrainingInteractor{tr, nr, ts}
}

func (i *TrainingInteractor) Create(inTraining app.Training) (*app.Training, error) {
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

	trnWithItems, err := i.TrainingService.Build(inTraining)
	if err != nil {
		return nil, err
	}

	training, err := i.TrainingRepo.Create(*trnWithItems)
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
	training, err := i.TrainingRepo.Get(trainingId)
	if err != nil {
		return nil, err
	}

	if *training.OwnerId != *actorId {
		return nil, errors.New("Forbidden, only training owner can do this.")
	}

	trainingItem, err := i.TrainingRepo.NextItem(trainingId)
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
