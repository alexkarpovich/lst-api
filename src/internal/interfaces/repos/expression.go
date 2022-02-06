package repos

import (
	"database/sql"
	"log"

	"github.com/alexkarpovich/lst-api/src/internal/domain"
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
	"github.com/alexkarpovich/lst-api/src/internal/interfaces/db"
)

type ExpressionRepo struct {
	db db.DB
}

func NewExpressionRepo(db db.DB) *ExpressionRepo {
	return &ExpressionRepo{db}
}

func (r *ExpressionRepo) Create(obj *domain.Expression) (*domain.Expression, error) {
	stmt := `
		INSERT INTO expressions (lang, value) VALUES(:lang, :value)
		RETURNING id
	`
	rows, err := r.db.Db().NamedQuery(stmt, obj)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		rows.Scan(&obj.Id)
	}

	if err = rows.Close(); err != nil {
		log.Println(err)
	}

	return obj, nil
}

func (r *ExpressionRepo) Get(id *valueobject.ID) (*domain.Expression, error) {
	expression := &domain.Expression{}
	query := `SELECT * FROM expressions WHERE id=$1`
	err := r.db.Db().Get(expression, query, id)
	if err != nil {
		return nil, err
	}

	return expression, nil
}

func (r *ExpressionRepo) Search(langCode string, value string) ([]*domain.Expression, error) {
	var expressions []*domain.Expression
	query := `SELECT * FROM expressions WHERE lang=$1 AND value LIKE $2`

	err := r.db.Db().Select(&expressions, query, langCode, "%"+value+"%")
	if err != nil {
		return nil, err
	}

	return expressions, nil
}

func (r *ExpressionRepo) CreateTranscription(expressionId *valueobject.ID, inTranscription domain.Transcription) (*domain.Transcription, error) {

	transcription := &inTranscription

	tx, err := r.db.Db().Begin()
	if err != nil {
		return nil, err
	}

	query := `SELECT id FROM transcriptions WHERE type=$1 AND value=$2`
	err = tx.QueryRow(query, transcription.Type, transcription.Value).
		Scan(&transcription.Id)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	if transcription.Id == nil {
		query = `INSERT INTO transcriptions (type, value) VALUES ($1, $2) RETURNING id`
		err = tx.QueryRow(query, transcription.Type, transcription.Value).
			Scan(&transcription.Id)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	query = `
		INSERT INTO expression_transcription (expression_id, transcription_id) 
		VALUES ($1, $2)
	`
	_, err = tx.Exec(query, expressionId, transcription.Id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return transcription, nil
}

func (r *ExpressionRepo) GetTranscriptionMap(typeId *valueobject.ID, expressionId *valueobject.ID) (map[string][]*domain.TranscriptionItem, error) {
	query := `
		SELECT e.value, t.id, t.value FROM expressions e
		LEFT JOIN expression_transcription et ON et.expression_id=e.id
		LEFT JOIN transcriptions t ON et.transcription_id=t.id
		WHERE t.type=$1 and (SELECT value FROM expressions WHERE id=$2) LIKE '%' || e.value || '%'
		ORDER BY length(e.value) desc 
	`
	rows, err := r.db.Db().Query(query, typeId, expressionId)
	if err != nil {
		return nil, err
	}

	transcriptionMap := make(map[string][]*domain.TranscriptionItem)

	for rows.Next() {
		var exprValue string
		item := &domain.TranscriptionItem{}
		rows.Scan(&exprValue, &item.Id, &item.Value)
		transcriptionMap[exprValue] = append(transcriptionMap[exprValue], item)
	}

	return transcriptionMap, nil
}
