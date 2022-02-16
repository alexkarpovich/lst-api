package repos

import (
	"fmt"

	"github.com/alexkarpovich/lst-api/src/internal/app"
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
	"github.com/alexkarpovich/lst-api/src/internal/interfaces/db"
	"github.com/jmoiron/sqlx"
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
		INSERT INTO trainings (owner_id, type, transcription_type, slices)
		VALUES($1, $2, $3, $4)
		RETURNING id
	`
	err = tx.QueryRow(query, training.OwnerId, training.Type, training.TranscriptionTypeId, pq.Array(training.Slices)).
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

func (r *TrainingRepo) Reset(trainingId *valueobject.ID) error {
	query := `
		UPDATE training_items SET complete=FALSE
		WHERE training_id=$1
	`
	_, err := r.db.Db().Exec(query, trainingId)
	if err != nil {
		return err
	}

	return nil
}

func (r *TrainingRepo) getMeta(trainingId *valueobject.ID) (*app.TrainingMeta, error) {
	query := `
		SELECT (SELECT COUNT(id) FROM training_items WHERE training_id=$1 AND stage=1) itemCount,
			(SELECT COUNT(id) FROM training_items WHERE training_id=$1 AND complete=TRUE) completeCount,
			(SELECT COUNT(DISTINCT stage) FROM training_items WHERE training_id=$1) stageCount
	`
	meta := &app.TrainingMeta{}
	err := r.db.Db().QueryRow(query, trainingId).
		Scan(&meta.UniqueItemCount, &meta.CompleteCount, &meta.StageCount)
	if err != nil {
		return nil, err
	}
	return meta, nil
}

func (r *TrainingRepo) Get(trainingId *valueobject.ID) (*app.Training, error) {
	query := `
		SELECT id, owner_id, type, transcription_type, slices FROM trainings
		WHERE id=$1
	`
	training := &app.Training{}
	sliceArr := pq.Int64Array{}
	err := r.db.Db().QueryRow(query, trainingId).
		Scan(&training.Id, &training.OwnerId, &training.Type, &training.TranscriptionTypeId, &sliceArr)
	if err != nil {
		return nil, err
	}

	for _, sliceId := range sliceArr {
		training.Slices = append(training.Slices, valueobject.ID(sliceId))
	}

	training.Meta, err = r.getMeta(trainingId)
	if err != nil {
		return nil, err
	}

	return training, nil
}

func (r *TrainingRepo) List(ownerId *valueobject.ID) ([]*app.Training, error) {
	query := `
		SELECT id, type, transcription_type, slices FROM trainings
		WHERE owner_id=$1
	`
	trainings := []*app.Training{}
	rows, err := r.db.Db().Queryx(query, ownerId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		sliceArr := pq.Int64Array{}
		training := &app.Training{OwnerId: ownerId}
		rows.Scan(&training.Id, &training.Type, &training.TranscriptionTypeId, &sliceArr)
		trainings = append(trainings, training)

		for _, sliceId := range sliceArr {
			training.Slices = append(training.Slices, valueobject.ID(sliceId))
		}

		training.Meta, err = r.getMeta(training.Id)
		if err != nil {
			return nil, err
		}
	}

	return trainings, nil
}

func (r *TrainingRepo) GetByItemId(itemId *valueobject.ID) (*app.Training, error) {
	query := `
		SELECT id, owner_id, type, transcription_type, slices FROM trainings
		WHERE id = (SELECT training_id FROM training_items WHERE id=$1)
	`
	training := &app.Training{}
	sliceArr := pq.Int64Array{}
	err := r.db.Db().QueryRow(query, itemId).
		Scan(&training.Id, &training.OwnerId, &training.Type, &training.TranscriptionTypeId, &sliceArr)
	if err != nil {
		return nil, err
	}

	for _, sliceId := range sliceArr {
		training.Slices = append(training.Slices, valueobject.ID(sliceId))
	}

	training.Meta, err = r.getMeta(training.Id)
	if err != nil {
		return nil, err
	}

	return training, nil
}

func (r *TrainingRepo) GetBySlices(inTraining app.Training) (*app.Training, error) {
	training := &inTraining
	query := `
		SELECT id FROM trainings
		WHERE owner_id=? AND type=? AND transcription_type=? AND slices = array[?]::smallint[]
	`
	query, args, err := sqlx.In(query, inTraining.OwnerId, inTraining.Type, inTraining.TranscriptionTypeId, inTraining.Slices)
	query = r.db.Db().Rebind(query)
	err = r.db.Db().QueryRow(query, args...).
		Scan(&training.Id)
	if err != nil {
		return nil, err
	}

	training.Meta, err = r.getMeta(training.Id)
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
		ORDER BY random()
		LIMIT 1
	`
	trainingItem := &app.TrainingItem{
		TrainingId: trainingId,
	}
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

func (r *TrainingRepo) ItemAnswers(itemId *valueobject.ID) ([]*app.TrainingAnswer, error) {
	training, err := r.GetByItemId(itemId)
	if err != nil {
		return nil, err
	}
	query := `
		SELECT e.id, e.value, t.id FROM expressions e
		LEFT JOIN translations t ON t.target_id=e.id
		LEFT JOIN node_translation nt ON nt.translation_id=t.id
		LEFT JOIN training_items ti ON ti.expression_id=t.native_id
		WHERE ti.id=? AND nt.node_id IN (?)
	`
	answers := []*app.TrainingAnswer{}
	query, args, err := sqlx.In(query, itemId, training.Slices)
	query = r.db.Db().Rebind(query)
	rows, err := r.db.Db().Query(query, args...)
	if err != nil {
		return nil, err
	}

	var translationId *valueobject.ID
	for rows.Next() {
		answer := &app.TrainingAnswer{}
		rows.Scan(&answer.Id, &answer.Value, &translationId)

		query = `
			SELECT t.id, t.value FROM transcriptions t
			LEFT JOIN translation_transcription tt ON tt.transcription_id=t.id
			WHERE tt.translation_id=$1 AND t.type=$2
		`
		transcriptions := []*app.Transcription{}
		err = r.db.Db().Select(&transcriptions, query, translationId, training.TranscriptionTypeId)

		answer.Transcriptions = transcriptions

		answers = append(answers, answer)
	}

	return answers, nil
}

func (r *TrainingRepo) MarkItemAsComplete(itemId *valueobject.ID) error {
	query := `
		UPDATE training_items SET complete=TRUE
		WHERE id=$1
	`
	_, err := r.db.Db().Exec(query, itemId)
	if err != nil {
		return err
	}

	return nil
}

func (r *TrainingRepo) HasCreatePermission(userId *valueobject.ID, nodes []valueobject.ID) bool {
	return true
}
