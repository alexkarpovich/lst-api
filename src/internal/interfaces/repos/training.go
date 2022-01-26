package repos

import (
	"fmt"

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
		INSERT INTO trainings (owner_id, type, slices)
		VALUES($1, $2, $3)
		RETURNING id
	`
	err = tx.QueryRow(query, training.OwnerId, training.Type, pq.Array(training.Slices)).
		Scan(&training.Id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	query = fmt.Sprintf(`
		INSERT INTO training_items (training_id, expression_id, stage, cycle, complete)
		VALUES (%d, :expression_id, :stage, :cycle, :complete)
	`, *training.Id)
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
		SELECT id, owner_id, type, slices FROM trainings
		WHERE id=$1
	`
	training := &app.Training{}
	sliceArr := pq.Int64Array{}
	err := r.db.Db().QueryRow(query, trainingId).
		Scan(&training.Id, &training.OwnerId, &training.Type, &sliceArr)
	if err != nil {
		return nil, err
	}

	for sliceId := range sliceArr {
		training.Slices = append(training.Slices, valueobject.ID(sliceId))
	}

	return training, nil
}

func (r *TrainingRepo) GetByItemId(itemId *valueobject.ID) (*app.Training, error) {
	query := `
		SELECT id, owner_id, type, slices FROM trainings
		WHERE id = (SELECT training_id FROM training_items WHERE id=$1)
	`
	training := &app.Training{}
	err := r.db.Db().Get(&training, query, itemId)
	if err != nil {
		return nil, err
	}

	return training, nil
}

func (r *TrainingRepo) NextItem(trainingId *valueobject.ID) (*app.TrainingItem, error) {
	query := `
		SELECT ti.id, ti.expression_id, ti.stage, ti.cycle, e.value FROM training_items ti
		LEFT JOIN expressions e ON e.id=ti.expression_id
		WHERE training_id=$1 AND complete=FALSE AND cycle = (
			SELECT MIN(cycle) FROM training_items WHERE training_id=$1 AND complete=FALSE
		)
		ORDER BY RANDOM()
		LIMIT 1
	`
	trainingItem := &app.TrainingItem{}
	expr := &app.TrainingExpression{}
	err := r.db.Db().QueryRow(query, trainingId).
		Scan(&trainingItem.Id, &trainingItem.ExpressionId, &trainingItem.Stage, &trainingItem.Cycle, &expr.Value)
	if err != nil {
		return nil, err
	}
	expr.Id = trainingItem.ExpressionId
	trainingItem.Expression = expr

	return trainingItem, nil
}

func (r *TrainingRepo) HasCreatePermission(userId *valueobject.ID, nodes []valueobject.ID) bool {
	return true
}
