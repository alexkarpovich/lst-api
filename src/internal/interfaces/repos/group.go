package repos

import (
	"log"

	"github.com/alexkarpovich/lst-api/src/internal/app"
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
	"github.com/alexkarpovich/lst-api/src/internal/interfaces/db"
)

type GroupRepo struct {
	db db.DB
}

func NewGroupRepo(db db.DB) *GroupRepo {
	return &GroupRepo{db}
}

func (r *GroupRepo) Create(obj *app.Group) (*app.Group, error) {
	var id *valueobject.ID

	stmt := `
		INSERT INTO groups (name, target_lang, native_lang, status) 
		VALUES(:name, :target_lang, :native_lang, :status)`

	rows, err := r.db.Db().NamedQuery(stmt, obj)

	if err != nil {
		return nil, err
	}

	if rows.Next() {
		rows.Scan(id)
	}

	if err = rows.Close(); err != nil {
		// but what should we do if there's an error?
		log.Println(err)
	}

	obj.Id = id

	return obj, nil
}

func (r *GroupRepo) List(userId *valueobject.ID) ([]*app.Group, error) {
	stmt := `
		SELECT g.* FROM groups g
		LEFT JOIN user_group ug ON ug.group_id=g.id
		WHERE user_id=$1
	`

	groups := []*app.Group{}
	err := r.db.Db().Select(&groups, stmt, userId)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (r *GroupRepo) MarkAsDeleted(groupId *valueobject.ID) error {
	stmt := `UPDATE groups SET status=$1 WHERE id=$2`

	_, err := r.db.Db().Exec(stmt, app.GroupDeleted, groupId)
	if err != nil {
		return err
	}

	return nil
}

func (r *GroupRepo) AttachUser(groupId *valueobject.ID, userId *valueobject.ID, role app.UserRole) error {
	stmt := `INSERT INTO user_group (group_id, user_id, role) VALUES ($1, $2, $3)`

	_, err := r.db.Db().Exec(stmt, groupId, userId, role)
	if err != nil {
		return err
	}

	return nil
}

func (r *GroupRepo) DetachUser(groupId *valueobject.ID, userId *valueobject.ID) error {
	stmt := `DELETE FROM user_group WHERE group_id=$1 AND user_id=$2`

	_, err := r.db.Db().Exec(stmt, app.GroupDeleted, groupId, userId)
	if err != nil {
		return err
	}

	return nil
}
