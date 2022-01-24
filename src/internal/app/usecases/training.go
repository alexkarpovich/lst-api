package usecases

import (
	"errors"

	"github.com/alexkarpovich/lst-api/src/internal/app"
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
)

type TrainingInteractor struct {
	TrainingRepo app.TrainingRepo
}

func NewTrainingInteractor(tr app.TrainingRepo) *TrainingInteractor {
	return &TrainingInteractor{tr}
}

func (i *TrainingInteractor) Create(inTraining app.Training) (*app.Training, error) {
	if !i.TrainingRepo.HasCreatePermission(inTraining.OwnerId, inTraining.Nodes) {
		return nil, errors.New("Forbidden, only user which has at least read role can do this.")
	}

	if len(inTraining.Nodes) == 0 {
		return nil, errors.New("There must be at least one node.")
	}

	training, err := i.TrainingRepo.Create(inTraining)
	if err != nil {
		return nil, err
	}

	return training, nil
}

func (i *TrainingInteractor) Reset(actorId *valueobject.ID, trainingId *valueobject.ID) error {
	training, err := i.TrainingRepo.Get(trainingId)
	if err != nil {
		return err
	}

	if training.OwnerId != actorId {
		return errors.New("Forbidden, only training owner can do this.")
	}

	return nil
}

func (i *TrainingInteractor) Next(actorId *valueobject.ID, trainingId *valueobject.ID) (*app.TrainingItem, error) {
	training, err := i.TrainingRepo.Get(trainingId)
	if err != nil {
		return nil, err
	}

	if training.OwnerId != actorId {
		return nil, errors.New("Forbidden, only training owner can do this.")
	}

	return nil, nil
}

func (i *TrainingInteractor) GetItem(actorId *valueobject.ID, itemId *valueobject.ID) (*app.TrainingItem, error) {
	training, err := i.TrainingRepo.GetByItemId(itemId)
	if err != nil {
		return nil, err
	}

	if training.OwnerId != actorId {
		return nil, errors.New("Forbidden, only training owner can do this.")
	}

	return nil, nil
}

func (i *TrainingInteractor) Complete(actorId *valueobject.ID, itemId *valueobject.ID) error {
	training, err := i.TrainingRepo.GetByItemId(itemId)
	if err != nil {
		return err
	}

	if training.OwnerId != actorId {
		return errors.New("Forbidden, only training owner can do this.")
	}

	return nil
}
