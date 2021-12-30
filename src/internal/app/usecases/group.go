package usecases

import (
	"log"
	"strings"

	"github.com/alexkarpovich/lst-api/src/internal/app"
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
)

type GroupInteractor struct {
	GroupRepo app.GroupRepo
	SliceRepo app.SliceRepo
}

func NewGroupInteractor(gr app.GroupRepo, fr app.SliceRepo) *GroupInteractor {
	return &GroupInteractor{gr, fr}
}

func (i *GroupInteractor) CreateGroup(obj *app.Group) (*app.Group, error) {
	obj.Status = app.GroupActive
	obj.Name = strings.TrimSpace(obj.Name)

	profile, err := i.GroupRepo.Create(obj)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return profile, nil
}

func (i *GroupInteractor) ListGroups(userId *valueobject.ID) ([]*app.Group, error) {
	profiles, err := i.GroupRepo.List(userId)
	if err != nil {
		log.Println(err)
	}

	return profiles, nil
}

func (i *GroupInteractor) MarkGroupAsDeleted(groupId *valueobject.ID) error {
	err := i.GroupRepo.MarkAsDeleted(groupId)
	if err != nil {
		log.Println(err)
	}

	return nil
}

func (i *GroupInteractor) AttachMember(groupId *valueobject.ID, userId *valueobject.ID, role app.UserRole) error {
	err := i.GroupRepo.AttachMember(groupId, userId, role)
	if err != nil {
		log.Println(err)
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
