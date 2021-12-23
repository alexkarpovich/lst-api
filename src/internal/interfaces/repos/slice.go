package repos

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	"github.com/alexkarpovich/lst-api/src/internal/app"
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
	"github.com/alexkarpovich/lst-api/src/internal/interfaces/db"
	"github.com/lib/pq"
)

type SliceRepo struct {
	db db.DB
}

func NewSliceRepo(db db.DB) *SliceRepo {
	return &SliceRepo{db}
}

func (r *SliceRepo) Create(groupId *valueobject.ID, obj *app.Slice) (*app.Slice, error) {
	var stmt string

	tx, err := r.db.Db().Beginx()
	if err != nil {
		return nil, err
	}

	stmt = `
		INSERT INTO slices (name, visibility) 
		VALUES(:name, :visibility)
		RETURNING id
	`

	rows, err := tx.NamedQuery(stmt, obj)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if rows.Next() {
		rows.Scan(&obj.Id)
	}

	if err = rows.Close(); err != nil {
		log.Println(err)
	}

	stmt = `INSERT INTO group_slice (group_id, slice_id, path) VALUES ($1, $2, $3)`
	_, err = tx.Exec(stmt, groupId, obj.Id, obj.Path)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return obj, nil
}

func (r *SliceRepo) Get(sliceId *valueobject.ID) (*app.Slice, error) {
	var err error
	var query string
	var slice app.Slice

	query = `
		SELECT s.*, gs.path FROM slices s
		LEFT JOIN group_slice gs ON gs.slice_id=s.id
		WHERE s.id=$1
	`

	err = r.db.Db().Get(&slice, query, sliceId)
	if err != nil {
		return nil, err
	}

	query = `
		SELECT t.id, t.target_id, e.value, t.comment, st.created_at FROM translations t
		LEFT JOIN slice_translation st ON st.translation_id=t.id
		LEFT JOIN expressions e ON e.id=t.native_id
		WHERE st.slice_id=$1
		ORDER BY st.created_at DESC
	`

	rows, err := r.db.Db().Query(query, sliceId)
	if err != nil {
		return nil, err
	}

	var targetExprId valueobject.ID
	exprTranslMap := make(map[valueobject.ID][]*app.Translation)

	for rows.Next() {
		tr := &app.Translation{}
		err = rows.Scan(&tr.Id, &targetExprId, &tr.Value, &tr.Comment, &tr.CreatedAt)

		exprTranslMap[targetExprId] = append(exprTranslMap[targetExprId], tr)
	}

	if err != nil {
		return nil, err
	}

	expressions := []*app.Expression{}

	query = `
		SELECT e.id, e.value, se.created_at FROM expressions e
		LEFT JOIN slice_expression se ON se.expression_id=e.id
		WHERE se.slice_id=$1
		ORDER BY se.created_at DESC;
	`

	rowsx, err := r.db.Db().Queryx(query, sliceId)
	if err != nil {
		return nil, err
	}

	for rowsx.Next() {
		expr := &app.Expression{}
		rowsx.StructScan(&expr)

		if translations, ok := exprTranslMap[*expr.Id]; ok {
			expr.Translations = translations
		}

		expressions = append(expressions, expr)
	}

	slice.Expressions = expressions

	return &slice, nil
}

func (r *SliceRepo) GetGroupBySlice(sliceId *valueobject.ID) (*app.Group, error) {
	var group app.Group
	query := `
		SELECT g.* FROM groups g
		LEFT JOIN group_slice gs ON gs.group_id=g.id
		WHERE gs.slice_id=$1
	`
	err := r.db.Db().Get(&group, query, sliceId)
	if err != nil {
		return nil, err
	}

	return &group, nil
}

func (r *SliceRepo) List(groupId *valueobject.ID) ([]*app.NestedSlice, error) {
	var result []*app.NestedSlice
	path_map := make(map[valueobject.ID]*app.NestedSlice)
	slices := []*app.NestedSlice{}

	err := r.db.Db().Select(&slices, `
		SELECT id, name, visibility, (SELECT COUNT(expression_id) FROM slice_expression se WHERE se.slice_id=s.id) as count, gs.path as path FROM slices s
		LEFT JOIN group_slice gs ON gs.slice_id=s.id
		WHERE gs.group_id=$1
		GROUP BY s.id, gs.path
		ORDER BY gs.path
	`, groupId)

	if err != nil {
		return nil, err
	}

	for _, slice := range slices {
		if len(slice.Path) == 0 {
			result = append(result, slice)
		} else {
			nodes := strings.Split(slice.Path, ".")
			parentIdx, _ := strconv.Atoi(nodes[len(nodes)-1])
			parentId := valueobject.ID(parentIdx)

			if s, ok := path_map[parentId]; ok {
				s.Children = append(s.Children, slice)
			}
		}

		path_map[*slice.Id] = slice
	}

	return result, nil
}

func (r *SliceRepo) AttachExpression(sliceId *valueobject.ID, expression *app.Expression) (*app.Expression, error) {
	var err error
	var query string

	group, err := r.GetGroupBySlice(sliceId)
	if err != nil {
		return nil, err
	}

	tx, err := r.db.Db().Begin()

	if expression.Id == nil {
		query = `SELECT id FROM expressions WHERE value=$1 AND lang=$2`
		err = tx.QueryRow(query, expression.Value, group.TargetLangCode).
			Scan(&expression.Id)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, err
			}
		}

		if expression.Id == nil {
			query = `
				INSERT INTO expressions (value, lang) VALUES($1, $2)
				RETURNING id
			`
			err = tx.QueryRow(query, expression.Value, group.TargetLangCode).
				Scan(&expression.Id)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	query = `
		INSERT INTO slice_expression (slice_id, expression_id)
		VALUES ($1, $2)
		RETURNING created_at
	`
	err = tx.QueryRow(query, sliceId, expression.Id).
		Scan(&expression.CreatedAt)
	if err != nil {
		pqErr := err.(*pq.Error)

		// If relation slice-expression already exists then return success
		if pqErr.Code == "23505" {
			return expression, nil
		}

		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return expression, nil
}

func (r *SliceRepo) DetachExpression(sliceId *valueobject.ID, expressionId *valueobject.ID) error {
	var query string

	tx, err := r.db.Db().Begin()
	if err != nil {
		return err
	}

	query = `
		DELETE FROM slice_translation 
		WHERE slice_id=$1 AND translation_id IN (
			SELECT id FROM translations WHERE target_id=$2 AND type=(SELECT id FROM object_types WHERE name='expression')
		)
	`
	_, err = tx.Exec(query, sliceId, expressionId)
	if err != nil {
		tx.Rollback()
		return err
	}

	query = `DELETE FROM slice_expression WHERE slice_id=$1 AND expression_id=$2`
	_, err = tx.Exec(query, sliceId, expressionId)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (r *SliceRepo) AttachTranslation(sliceId *valueobject.ID, expressionId *valueobject.ID, translation *app.Translation) (*app.Translation, error) {
	var query string
	var err error

	group, err := r.GetGroupBySlice(sliceId)
	if err != nil {
		return nil, err
	}

	tx, err := r.db.Db().Begin()

	if translation.Id == nil {
		var nativeId *valueobject.ID = nil

		query = `SELECT id FROM expressions WHERE lang=$1 AND value=$2`
		err = tx.QueryRow(query, group.NativeLangCode, translation.Value).
			Scan(&nativeId)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, err
			}
		}

		if nativeId == nil {
			query = `INSERT INTO expressions (lang, value) VALUES($1, $2) RETURNING id`
			err := tx.QueryRow(query, group.NativeLangCode, translation.Value).
				Scan(&nativeId)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		query = `
			SELECT id FROM translations
			WHERE type=(SELECT id FROM object_types WHERE name='expression') AND target_id=$1 AND native_id=$2`
		err = tx.QueryRow(query, expressionId, nativeId).
			Scan(&translation.Id)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, err
			}
		}

		if translation.Id == nil {
			query = `
				INSERT INTO translations (type, target_id, native_id, comment)
				SELECT id, $1, $2, $3 FROM object_types WHERE name='expression'
				RETURNING id
			`
			err = tx.QueryRow(query, expressionId, nativeId, translation.Comment).
				Scan(&translation.Id)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	query = `INSERT INTO slice_translation (slice_id, translation_id) VALUES($1, $2) RETURNING created_at`
	err = tx.QueryRow(query, sliceId, translation.Id).
		Scan(&translation.CreatedAt)
	if err != nil {
		pqErr := err.(*pq.Error)

		// If relation slice-expression already exists then return success
		if pqErr.Code == "23505" {
			return translation, nil
		}

		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return translation, nil
}

func (r *SliceRepo) DetachTranslation(sliceId *valueobject.ID, translationId *valueobject.ID) error {
	query := `DELETE FROM slice_translation WHERE slice_id=$1 AND translation_id=$2`
	_, err := r.db.Db().Exec(query, sliceId, translationId)
	if err != nil {
		return err
	}

	return nil
}

func (r *SliceRepo) AttachText(sliceId *valueobject.ID, text *app.Text) (*app.Text, error) {
	var err error
	var query string

	group, err := r.GetGroupBySlice(sliceId)
	if err != nil {
		return nil, err
	}

	tx, err := r.db.Db().Begin()

	if text.Id == nil {
		query = `
			INSERT INTO texts (author_id, title, content, lang) VALUES($1, $2, $3, $4)
			RETURNING id, created_at
		`
		err = tx.QueryRow(query, text.AuthorId, text.Title, text.Content, group.TargetLangCode).
			Scan(&text.Id, &text.CreatedAt)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	query = `UPDATE slices SET text_id=$1 WHERE id=$2`
	_, err = tx.Exec(query, text.Id, sliceId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return text, nil
}

func (r *SliceRepo) DetachText(sliceId *valueobject.ID, textId *valueobject.ID) error {
	query := `UPDATE slices SET text_id=NULL WHERE id=$1`
	_, err := r.db.Db().Exec(query, sliceId)
	if err != nil {
		return err
	}

	return nil
}
