package repos

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	"github.com/alexkarpovich/lst-api/src/internal/app"
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
	"github.com/alexkarpovich/lst-api/src/internal/interfaces/db"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type NodeRepo struct {
	db db.DB
}

func NewNodeRepo(db db.DB) *NodeRepo {
	return &NodeRepo{db}
}

func (r *NodeRepo) Create(groupId *valueobject.ID, obj app.Node) (*app.Node, error) {
	var query string

	tx, err := r.db.Db().Beginx()
	if err != nil {
		return nil, err
	}

	query = `
		INSERT INTO nodes (type, name, visibility) 
		VALUES(:type, :name, :visibility)
		RETURNING id
	`

	rows, err := tx.NamedQuery(query, obj)

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

	query = `INSERT INTO group_node (group_id, node_id, path) VALUES ($1, $2, $3)`
	_, err = tx.Exec(query, groupId, obj.Id, obj.Path)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &obj, nil
}

func (r *NodeRepo) Get(nodeId *valueobject.ID) (*app.Node, error) {
	var err error
	var query string
	var node app.Node

	query = `
		SELECT n.*, gn.path FROM nodes n
		LEFT JOIN group_node gn ON gn.node_id=n.id
		WHERE n.id=$1
	`

	err = r.db.Db().Get(&node, query, nodeId)
	if err != nil {
		return nil, err
	}

	query = `
		SELECT t.id, t.target_id, e.value, t.comment, nt.created_at FROM translations t
		LEFT JOIN node_translation nt ON nt.translation_id=t.id
		LEFT JOIN expressions e ON e.id=t.native_id
		WHERE nt.node_id=$1
		ORDER BY nt.created_at DESC
	`

	rows, err := r.db.Db().Query(query, nodeId)
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
		SELECT e.id, e.value, ne.created_at FROM expressions e
		LEFT JOIN node_expression ne ON ne.expression_id=e.id
		WHERE ne.node_id=$1
		ORDER BY ne.created_at DESC;
	`

	rowsx, err := r.db.Db().Queryx(query, nodeId)
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

	node.Expressions = expressions

	return &node, nil
}

func (r *NodeRepo) View(ids []valueobject.ID) (*app.NodeView, error) {
	var err error
	var query string
	var nodeView app.NodeView

	query = `
		SELECT t.id, t.target_id, e.value, t.comment, nt.created_at FROM translations t
		LEFT JOIN node_translation nt ON nt.translation_id=t.id
		LEFT JOIN expressions e ON e.id=t.native_id
		WHERE nt.node_id IN (?)
		ORDER BY nt.created_at DESC
	`
	query, args, err := sqlx.In(query, ids)
	query = r.db.Db().Rebind(query)
	rows, err := r.db.Db().Query(query, args...)
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
		SELECT e.id, e.value, ne.created_at FROM expressions e
		LEFT JOIN node_expression ne ON ne.expression_id=e.id
		WHERE ne.node_id IN (?)
		ORDER BY ne.created_at DESC;
	`

	query, args, err = sqlx.In(query, ids)
	query = r.db.Db().Rebind(query)
	rows, err = r.db.Db().Query(query, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		expr := &app.Expression{}
		rows.Scan(&expr.Id, &expr.Value, &expr.CreatedAt)

		if translations, ok := exprTranslMap[*expr.Id]; ok {
			expr.Translations = translations
		}

		expressions = append(expressions, expr)
	}

	nodeView.Expressions = expressions

	return &nodeView, nil
}

func (r *NodeRepo) GetGroupByNode(nodeId *valueobject.ID) (*app.Group, error) {
	var group app.Group
	query := `
		SELECT g.* FROM groups g
		LEFT JOIN group_node gn ON gn.group_id=g.id
		WHERE gn.node_id=$1
	`
	err := r.db.Db().Get(&group, query, nodeId)
	if err != nil {
		return nil, err
	}

	return &group, nil
}

func (r *NodeRepo) List(groupId *valueobject.ID) ([]*app.FlatNode, error) {
	var result []*app.FlatNode
	path_map := make(map[valueobject.ID]*app.FlatNode)
	nodes := []*app.FlatNode{}

	err := r.db.Db().Select(&nodes, `
		SELECT id, type, name, visibility, (
			SELECT COUNT(expression_id) FROM node_expression ne 
			LEFT JOIN group_node cgn ON cgn.node_id=ne.node_id 
			WHERE cgn.group_id=$1 AND (n.type=0 AND index(cgn.path, CASE WHEN gn.path='' THEN concat(gn.node_id) ELSE concat(gn.path,'.',gn.node_id) END::ltree) <> -1) OR (n.type=1 AND ne.node_id=n.id)
			) as count, 
			gn.path as path FROM nodes n
		LEFT JOIN group_node gn ON gn.node_id=n.id
		WHERE gn.group_id=$1
		GROUP BY n.id, gn.node_id, gn.path
		ORDER BY gn.path
	`, groupId)

	if err != nil {
		return nil, err
	}

	for _, node := range nodes {
		if len(node.Path) == 0 {
			result = append(result, node)
		} else {
			nodes := strings.Split(node.Path, ".")
			parentIdx, _ := strconv.Atoi(nodes[len(nodes)-1])
			parentId := valueobject.ID(parentIdx)

			if s, ok := path_map[parentId]; ok {
				s.Children = append(s.Children, node)
			}
		}

		path_map[*node.Id] = node
	}

	return result, nil
}

func (r *NodeRepo) AttachExpression(nodeId *valueobject.ID, expression app.Expression) (*app.Expression, error) {
	var err error
	var query string

	group, err := r.GetGroupByNode(nodeId)
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
		INSERT INTO node_expression (node_id, expression_id)
		VALUES ($1, $2)
		RETURNING created_at
	`
	err = tx.QueryRow(query, nodeId, expression.Id).
		Scan(&expression.CreatedAt)
	if err != nil {
		pqErr := err.(*pq.Error)

		// If relation slice-expression already exists then return success
		if pqErr.Code == "23505" {
			return &expression, nil
		}

		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &expression, nil
}

func (r *NodeRepo) DetachExpression(nodeId *valueobject.ID, expressionId *valueobject.ID) error {
	var query string

	tx, err := r.db.Db().Begin()
	if err != nil {
		return err
	}

	query = `
		DELETE FROM node_translation 
		WHERE node_id=$1 AND translation_id IN (
			SELECT id FROM translations WHERE target_id=$2 AND type=(SELECT id FROM object_types WHERE name='expression')
		)
	`
	_, err = tx.Exec(query, nodeId, expressionId)
	if err != nil {
		tx.Rollback()
		return err
	}

	query = `DELETE FROM node_expression WHERE node_id=$1 AND expression_id=$2`
	_, err = tx.Exec(query, nodeId, expressionId)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (r *NodeRepo) AttachTranslation(nodeId *valueobject.ID, expressionId *valueobject.ID, translation app.Translation) (*app.Translation, error) {
	var query string
	var err error

	group, err := r.GetGroupByNode(nodeId)
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

	query = `INSERT INTO node_translation (node_id, translation_id) VALUES($1, $2) RETURNING created_at`
	err = tx.QueryRow(query, nodeId, translation.Id).
		Scan(&translation.CreatedAt)
	if err != nil {
		pqErr := err.(*pq.Error)

		// If relation slice-expression already exists then return success
		if pqErr.Code == "23505" {
			return &translation, nil
		}

		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &translation, nil
}

func (r *NodeRepo) DetachTranslation(nodeId *valueobject.ID, translationId *valueobject.ID) error {
	query := `DELETE FROM node_translation WHERE node_id=$1 AND translation_id=$2`
	_, err := r.db.Db().Exec(query, nodeId, translationId)
	if err != nil {
		return err
	}

	return nil
}

func (r *NodeRepo) AttachText(nodeId *valueobject.ID, text app.Text) (*app.Text, error) {
	var err error
	var query string

	group, err := r.GetGroupByNode(nodeId)
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

	query = `UPDATE nodes SET text_id=$1 WHERE id=$2`
	_, err = tx.Exec(query, text.Id, nodeId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &text, nil
}

func (r *NodeRepo) DetachText(nodeId *valueobject.ID, textId *valueobject.ID) error {
	query := `UPDATE nodes SET text_id=NULL WHERE id=$1`
	_, err := r.db.Db().Exec(query, nodeId)
	if err != nil {
		return err
	}

	return nil
}
