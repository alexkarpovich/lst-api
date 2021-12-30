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
	query := `
		SELECT g.* FROM groups g
		LEFT JOIN user_group ug ON ug.group_id=g.id
		WHERE user_id=$1
	`

	groups := []*app.Group{}
	groupMap := make(map[valueobject.ID]*app.Group)
	rows, err := r.db.Db().Query(query, userId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		group := &app.Group{}
		rows.Scan(&group.Id, &group.TargetLangCode, &group.NativeLangCode, &group.Name, &group.Status)

		groupMap[*group.Id] = group
		groups = append(groups, group)
	}

	query = `
		SELECT u.id, u.username, ug1.role, ug1.group_id FROM user_group ug
		LEFT JOIN user_group ug1 ON ug1.group_id=ug.group_id
		LEFT JOIN users u ON u.id=ug1.user_id
		where ug.user_id=$1;
	`

	rows, err = r.db.Db().Query(query, userId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var groupId valueobject.ID
		member := &app.GroupMember{}
		rows.Scan(&member.Id, &member.Username, &member.Role, &groupId)

		if group, ok := groupMap[groupId]; ok {
			group.Members = append(group.Members, member)
		}
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

func (r *GroupRepo) AttachMember(groupId *valueobject.ID, userId *valueobject.ID, role app.UserRole) error {
	stmt := `INSERT INTO user_group (group_id, user_id, role) VALUES ($1, $2, $3)`

	_, err := r.db.Db().Exec(stmt, groupId, userId, role)
	if err != nil {
		return err
	}

	return nil
}

func (r *GroupRepo) DetachMember(groupId *valueobject.ID, userId *valueobject.ID) error {
	stmt := `DELETE FROM user_group WHERE group_id=$1 AND user_id=$2`

	_, err := r.db.Db().Exec(stmt, app.GroupDeleted, groupId, userId)
	if err != nil {
		return err
	}

	return nil
}
