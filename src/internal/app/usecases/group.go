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
	SliceRepo app.SliceRepo
	UserRepo  app.UserRepo
	Email     services.EmailService
}

func NewGroupInteractor(gr app.GroupRepo, fr app.SliceRepo, ur app.UserRepo, es services.EmailService) *GroupInteractor {
	return &GroupInteractor{gr, fr, ur, es}
}

func (i *GroupInteractor) CreateGroup(admin *valueobject.ID, obj *app.Group) (*app.Group, error) {
	obj.Status = app.GroupActive
	obj.Name = strings.TrimSpace(obj.Name)

	group, err := i.GroupRepo.Create(admin, obj)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return group, nil
}

func (i *GroupInteractor) UpdateGroup(obj *app.Group) error {
	obj.Name = strings.TrimSpace(obj.Name)

	err := i.GroupRepo.Update(obj)
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

func (i *GroupInteractor) InviteUser(groupId *valueobject.ID, userId *valueobject.ID) error {
	member := app.GroupMember{
		Id:             userId,
		Role:           app.UserReader,
		Status:         app.MemberPending,
		Token:          pkg.RandomString(128),
		TokenExpiresAt: time.Now().Add(3 * 24 * time.Hour),
	}

	err := i.GroupRepo.AttachUser(groupId, member)
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

func (i *GroupInteractor) CreateSlice(groupId *valueobject.ID, s *app.Slice) (*app.Slice, error) {
	s.Visibility = app.SlicePrivate
	s.Name = strings.TrimSpace(s.Name)

	slice, err := i.SliceRepo.Create(groupId, s)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return slice, nil
}

func (i *GroupInteractor) ListSlices(groupId *valueobject.ID) ([]*app.NestedSlice, error) {
	folders, err := i.SliceRepo.List(groupId)
	if err != nil {
		return nil, err
	}

	return folders, nil
}
