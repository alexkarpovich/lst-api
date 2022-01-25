package repos

import (
	"github.com/alexkarpovich/lst-api/src/internal/app"
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
	"github.com/alexkarpovich/lst-api/src/internal/interfaces/db"
	"github.com/lib/pq"
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
	err = tx.QueryRow(query, training.OwnerId, training.Type, pq.Array(training.Slices)).
		Scan(&training.Id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	query = `
		INSERT INTO training_items (training_id, expression_id, stage, cycle, completed)
		VALUES (:training_id, :expression_id, :stage, :cycle, :completed)
	`
	_, err = tx.NamedExec(query, training.Items)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return training, nil
}

func (r *TrainingRepo) Get(trainingId *valueobject.ID) (*app.Training, error) {
	query := `
		SELECT id, owner_id, type, nodes FROM trainings
		WHERE id=$1
	`
	training := &app.Training{}
	err := r.db.Db().Get(&training, query, trainingId)
	if err != nil {
		return nil, err
	}

	return training, nil
}

func (r *TrainingRepo) GetByItemId(itemId *valueobject.ID) (*app.Training, error) {
	query := `
		SELECT id, owner_id, type, nodes FROM trainings
		WHERE id = (SELECT training_id FROM training_items WHERE id=$1)
	`
	training := &app.Training{}
	err := r.db.Db().Get(&training, query, itemId)
	if err != nil {
		return nil, err
	}

	return training, nil
}

func (r *TrainingRepo) HasCreatePermission(userId *valueobject.ID, nodes []valueobject.ID) bool {
	return true
}
