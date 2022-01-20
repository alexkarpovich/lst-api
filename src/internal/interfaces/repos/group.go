package repos

import (
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

func (r *GroupRepo) isGroupUntouched(groupId *valueobject.ID) (bool, error) {
	var usersCount, slicesCount int
	query := `
		SELECT
			(SELECT COUNT(user_id) FROM user_group WHERE group_id=$1) as users_count,
			(SELECT coalesce(COUNT(slice_id), 0) FROM group_slice WHERE group_id=$1) as slices_count
	`
	err := r.db.Db().QueryRow(query, groupId).
		Scan(&usersCount, &slicesCount)
	if err != nil {
		return false, err
	}

	return usersCount == 1 && slicesCount == 0, nil

}

func (r *GroupRepo) Create(userId *valueobject.ID, obj app.Group) (*app.Group, error) {
	tx, err := r.db.Db().Begin()
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO groups (name, target_lang, native_lang, status) 
		VALUES($1, $2, $3, $4)
		RETURNING id
	`

	err = tx.QueryRow(query, obj.Name, obj.TargetLangCode, obj.NativeLangCode, obj.Status).
		Scan(&obj.Id)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	query = `INSERT INTO user_group (user_id, group_id, role, status) VALUES($1, $2, $3, $4)`

	_, err = tx.Exec(query, userId, obj.Id, app.UserAdmin, app.MemberActive)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &obj, nil
}

func (r *GroupRepo) Update(obj app.Group) error {
	var query string
	var err error

	isUntouched, err := r.isGroupUntouched(obj.Id)
	if err != nil {
		return err
	}
	if isUntouched {
		query = `UPDATE groups SET name=$1, target_lang=$2, native_lang=$3 WHERE id=$4`

		_, err = r.db.Db().Exec(query, obj.Name, obj.TargetLangCode, obj.NativeLangCode, obj.Id)
	} else {
		query = `UPDATE groups SET name=$1 WHERE id=$2`

		_, err = r.db.Db().Exec(query, obj.Name, obj.Id)
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *GroupRepo) List(userId *valueobject.ID) ([]*app.Group, error) {
	query := `
		SELECT 
			g.id, g.target_lang, g.native_lang, g.name, g.status, g.config,
			(SELECT COUNT(user_id) FROM user_group WHERE group_id=g.id) as users_count,
			(SELECT coalesce(COUNT(node_id), 0) FROM group_node WHERE group_id=g.id) as node_count 
		FROM groups g
		LEFT JOIN user_group ug ON ug.group_id=g.id
		WHERE user_id=$1 AND g.status != $2
		ORDER BY id DESC
	`
	groups := []*app.Group{}
	groupMap := make(map[valueobject.ID]*app.Group)
	rows, err := r.db.Db().Query(query, userId, app.GroupDeleted)
	if err != nil {
		return nil, err
	}

	var usersCount, nodesCount uint

	for rows.Next() {
		group := &app.Group{}
		rows.Scan(&group.Id, &group.TargetLangCode, &group.NativeLangCode, &group.Name, &group.Status, &group.Config, &usersCount, &nodesCount)

		if usersCount == 1 || nodesCount == 0 {
			group.IsUntouched = true
		} else {
			group.IsUntouched = false
		}

		groupMap[*group.Id] = group
		groups = append(groups, group)
	}

	query = `
		SELECT u.id, u.username, ug1.role, ug1.status, ug1.group_id FROM user_group ug
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
		rows.Scan(&member.Id, &member.Username, &member.Role, &member.Status, &groupId)

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

func (r *GroupRepo) FindMemberById(groupId *valueobject.ID, memberId *valueobject.ID) (*app.GroupMember, error) {
	member := &app.GroupMember{
		Id: memberId,
	}

	query := `
		SELECT role, status FROM user_group 
		WHERE group_id=$1 AND user_id=$2
	`
	err := r.db.Db().QueryRow(query, groupId, memberId).
		Scan(&member.Role, &member.Status)
	if err != nil {
		return nil, err
	}

	return member, nil
}

func (r *GroupRepo) FindMemberByNodeId(nodeId *valueobject.ID, memberId *valueobject.ID) (*app.GroupMember, error) {
	member := &app.GroupMember{
		Id: memberId,
	}

	query := `
		SELECT u.username, ug.role, ug.status FROM user_group ug
		LEFT JOIN users u ON u.id=ug.user_id
		LEFT JOIN group_node gn ON gn.group_id=ug.group_id
		WHERE gn.node_id=$1 AND ug.user_id=$2
	`
	err := r.db.Db().QueryRow(query, nodeId, memberId).
		Scan(&member.Username, &member.Role, &member.Status)
	if err != nil {
		return nil, err
	}

	return member, nil
}

func (r *GroupRepo) FindMemberByToken(token string) (*valueobject.ID, *app.GroupMember, error) {
	var groupId *valueobject.ID

	member := &app.GroupMember{}

	query := `
		SELECT group_id, user_id, role, status FROM user_group 
		WHERE token=$1 AND status=$2 AND token_expires_at > NOW()
	`
	err := r.db.Db().QueryRow(query, token, app.MemberPending).
		Scan(&groupId, &member.Id, &member.Role, &member.Status)
	if err != nil {
		return nil, nil, err
	}

	return groupId, member, nil
}

func (r *GroupRepo) AttachUser(groupId *valueobject.ID, member app.GroupMember) error {
	query := `
		INSERT INTO user_group (group_id, user_id, role, status, token, token_expires_at) 
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Db().Exec(query, groupId, member.Id, member.Role, member.Status, member.Token, member.TokenExpiresAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *GroupRepo) DetachMember(groupId *valueobject.ID, userId *valueobject.ID) error {
	query := `DELETE FROM user_group WHERE group_id=$1 AND user_id=$2`

	_, err := r.db.Db().Exec(query, groupId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r *GroupRepo) UpdateMember(groupId *valueobject.ID, member app.GroupMember) error {
	query := `
		UPDATE user_group SET role=$1, status=$2 
		WHERE group_id=$3 AND user_id=$4	
	`

	_, err := r.db.Db().Exec(query, member.Role, member.Status, groupId, member.Id)
	if err != nil {
		return err
	}

	return nil
}
