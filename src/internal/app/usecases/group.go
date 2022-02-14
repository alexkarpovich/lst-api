package usecases

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/alexkarpovich/lst-api/src/internal/app"
	"github.com/alexkarpovich/lst-api/src/internal/app/services"
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
	"github.com/alexkarpovich/lst-api/src/pkg"
)

type GroupInteractor struct {
	GroupRepo app.GroupRepo
	NodeRepo  app.NodeRepo
	UserRepo  app.UserRepo
	Email     services.EmailService
}

func NewGroupInteractor(gr app.GroupRepo, fr app.NodeRepo, ur app.UserRepo, es services.EmailService) *GroupInteractor {
	return &GroupInteractor{gr, fr, ur, es}
}

func (i *GroupInteractor) CreateGroup(actorId *valueobject.ID, obj app.Group) (*app.Group, error) {
	obj.Status = app.GroupActive
	obj.Name = strings.TrimSpace(obj.Name)

	group, err := i.GroupRepo.Create(actorId, obj)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	group, err = i.GroupRepo.Get(group.Id)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (i *GroupInteractor) UpdateGroup(actorId *valueobject.ID, obj app.Group) error {
	var err error

	member, err := i.GroupRepo.FindMemberById(obj.Id, obj.Id)
	if err != nil {
		return err
	}

	if member.Role == app.UserReader {
		return errors.New("Forbidden, only non-reader member can edit group.")
	}

	obj.Name = strings.TrimSpace(obj.Name)

	err = i.GroupRepo.Update(obj)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (i *GroupInteractor) ListGroups(userId *valueobject.ID) ([]*app.Group, error) {
	groups, err := i.GroupRepo.List(userId)
	if err != nil {
		log.Println(err)
	}

	return groups, nil
}

func (i *GroupInteractor) MarkGroupAsDeleted(userId *valueobject.ID, groupId *valueobject.ID) error {
	err := i.GroupRepo.MarkAsDeleted(groupId)
	if err != nil {
		log.Println(err)
	}

	return nil
}

func (i *GroupInteractor) InviteUser(actorId *valueobject.ID, groupId *valueobject.ID, userId *valueobject.ID) error {
	var err error

	actor, err := i.GroupRepo.FindMemberById(groupId, actorId)
	if err != nil {
		return err
	}

	if actor.Role != app.UserAdmin {
		return errors.New("Forbidden, only admin can invite member.")
	}

	member := app.GroupMember{
		Id:             userId,
		Role:           app.UserReader,
		Status:         app.MemberPending,
		Token:          pkg.RandomString(128),
		TokenExpiresAt: time.Now().Add(3 * 24 * time.Hour),
	}

	err = i.GroupRepo.AttachUser(groupId, member)
	if err != nil {
		log.Println(err)
	}

	user, err := i.UserRepo.Get(userId)
	if err != nil {
		return err
	}

	go i.Email.SendGroupInvitation(user.Email, member.Token)

	return nil
}

func (i *GroupInteractor) ConfirmInvitation(userId *valueobject.ID, token string) error {
	groupId, member, err := i.GroupRepo.FindMemberByToken(token)
	if err != nil {
		return err
	}

	if *userId != *member.Id {
		return errors.New("You don't have permissions to accept invitation.")
	}

	member.Status = app.MemberActive

	err = i.GroupRepo.UpdateMember(groupId, *member)
	if err != nil {
		return err
	}

	return nil
}

func (i *GroupInteractor) DetachMember(groupId *valueobject.ID, userId *valueobject.ID) error {
	err := i.GroupRepo.DetachMember(groupId, userId)
	if err != nil {
		log.Println(err)
	}

	return nil
}

func (i *GroupInteractor) UpdateMemberRole(actorId *valueobject.ID, groupId *valueobject.ID, member app.GroupMember) error {
	var err error

	actor, err := i.GroupRepo.FindMemberById(groupId, actorId)
	if err != nil {
		return err
	}

	if actor.Role != app.UserAdmin {
		return errors.New("Forbidden, only admin of a group can change member roles.")
	}

	mbr, err := i.GroupRepo.FindMemberById(groupId, member.Id)
	if err != nil {
		return err
	}

	mbr.Role = member.Role

	err = i.GroupRepo.UpdateMember(groupId, member)
	if err != nil {
		log.Println(err)
	}

	return nil
}

func (i *GroupInteractor) CreateNode(groupId *valueobject.ID, s app.Node) (*app.Node, error) {
	s.Visibility = app.NodePrivate
	s.Name = strings.TrimSpace(s.Name)

	slice, err := i.NodeRepo.Create(groupId, s)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	slice, err = i.NodeRepo.Get(slice.Id)
	if err != nil {
		return nil, err
	}

	return slice, nil
}

func (i *GroupInteractor) ListNodes(groupId *valueobject.ID) ([]*app.FlatNode, error) {
	folders, err := i.NodeRepo.List(groupId)
	if err != nil {
		return nil, err
	}

	return folders, nil
}

func (i *GroupInteractor) DeleteNode(groupId *valueobject.ID, nodeId *valueobject.ID) error {
	err := i.GroupRepo.DeleteNode(groupId, nodeId)
	if err != nil {
		return err
	}

	return nil
}
