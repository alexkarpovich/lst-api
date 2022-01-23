package repos

import (
	"github.com/alexkarpovich/lst-api/src/internal/app"
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
	"github.com/alexkarpovich/lst-api/src/internal/interfaces/db"
)

type TrainingRepo struct {
	db db.DB
}

func NewTrainingRepo(db db.DB) *TrainingRepo {
	return &TrainingRepo{db}
}

func (r *TrainingRepo) Create(inTraining app.Training) (*app.Training, error) {
	var query string
	var err error
	training := &inTraining

	tx, err := r.db.Db().Beginx()
	if err != nil {
		return nil, err
	}

	query = `
		INSERT INTO trainings (owner_id, type, nodes)
		VALUES($1, $2, $3)
		RETURNING id
	`
	err = tx.QueryRow(query, training.OwnerId, training.Type, training.Nodes).
		Scan(&training.Id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return training, nil
}

func (r *TrainingRepo) Get(trainingId *valueobject.ID) (*app.Training, error) {
	return nil, nil
}

func (r *TrainingRepo) GetByItemId(itemId *valueobject.ID) (*app.Training, error) {
	return nil, nil
}

func (r *TrainingRepo) HasCreatePermission(userId *valueobject.ID, nodes []valueobject.ID) bool {
	return true
}
