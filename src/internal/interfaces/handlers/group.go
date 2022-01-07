package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/alexkarpovich/lst-api/src/internal/app"
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
	"github.com/alexkarpovich/lst-api/src/internal/utils"
	"github.com/gorilla/mux"
)

type GroupInteractor interface {
	CreateGroup(*valueobject.ID, *app.Group) (*app.Group, error)
	UpdateGroup(*app.Group) error
	ListGroups(*valueobject.ID) ([]*app.Group, error)
	MarkGroupAsDeleted(*valueobject.ID, *valueobject.ID) error
	CreateSlice(*valueobject.ID, *app.Slice) (*app.Slice, error)
	ListSlices(*valueobject.ID) ([]*app.NestedSlice, error)
	InviteUser(*valueobject.ID, *valueobject.ID) error
	ConfirmInvitation(*valueobject.ID, string) error
	DetachMember(*valueobject.ID, *valueobject.ID) error
}

type groupHanlder struct {
	BaseHanlder
	groupInteractor GroupInteractor
}

func ConfigureGroupHandler(pi GroupInteractor, r *mux.Router) {
	h := &groupHanlder{
		BaseHanlder: BaseHanlder{
			router: r,
		},
		groupInteractor: pi,
	}

	h.router.HandleFunc("/me/group", h.CreateGroup()).Methods("POST")
	h.router.HandleFunc("/me/group", h.ListGroups()).Methods("GET")
	h.router.HandleFunc("/me/group/confirm-invitation/{token}", h.ConfirmInvitation()).Methods("POST")
	h.router.HandleFunc("/me/group/{groupId}", h.UpdateGroup()).Methods("POST")
	h.router.HandleFunc("/me/group/{groupId}", h.DeleteGroup()).Methods("DELETE")
	h.router.HandleFunc("/me/group/{groupId}/slice", h.CreateSlice()).Methods("POST")
	h.router.HandleFunc("/me/group/{groupId}/slice", h.ListSlices()).Methods("GET")
	h.router.HandleFunc("/me/group/{groupId}/invite-user/{userId}", h.InviteUser()).Methods("POST")
	h.router.HandleFunc("/me/group/{groupId}/detach-member/{memberId}", h.DetachMember()).Methods("POST")
}

func (i *groupHanlder) CreateGroup() http.HandlerFunc {
	type request struct {
		Name           string `json:"name"`
		TargetLangCode string `json:"targetLangCode"`
		NativeLangCode string `json:"nativeLangCode"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var s request
		var err error

		if err = json.NewDecoder(r.Body).Decode(&s); err != nil {
			utils.SendJsonError(w, "Invalid request data", http.StatusBadRequest)
			return
		}

		user := utils.LoggedInUser(r)
		if user == nil {
			log.Println("error group create context")
			return
		}

		group := &app.Group{
			Name:           s.Name,
			TargetLangCode: s.TargetLangCode,
			NativeLangCode: s.NativeLangCode,
		}

		if group, err = i.groupInteractor.CreateGroup(user.Id, group); err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, group, http.StatusOK)
	}
}

func (i *groupHanlder) UpdateGroup() http.HandlerFunc {
	type request struct {
		Name           string `json:"name"`
		TargetLangCode string `json:"targetLangCode"`
		NativeLangCode string `json:"nativeLangCode"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var s request
		var err error

		vars := mux.Vars(r)
		groupIdArg, err := strconv.Atoi(vars["groupId"])
		if err != nil {
			utils.SendJsonError(w, "Invalid group id", http.StatusBadRequest)
			return
		}
		groupId := valueobject.ID(groupIdArg)

		if err = json.NewDecoder(r.Body).Decode(&s); err != nil {
			utils.SendJsonError(w, "Invalid request data", http.StatusBadRequest)
			return
		}

		user := utils.LoggedInUser(r)
		if user == nil {
			log.Println("error group update context")
			return
		}

		group := &app.Group{
			Id:             &groupId,
			Name:           s.Name,
			TargetLangCode: s.TargetLangCode,
			NativeLangCode: s.NativeLangCode,
		}

		if err = i.groupInteractor.UpdateGroup(group); err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, "Success", http.StatusOK)
	}
}

func (i *groupHanlder) DeleteGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		vars := mux.Vars(r)
		groupIdArg, err := strconv.Atoi(vars["groupId"])
		if err != nil {
			utils.SendJsonError(w, "Invalid group id", http.StatusBadRequest)
			return
		}
		groupId := valueobject.ID(groupIdArg)

		user := utils.LoggedInUser(r)
		if user == nil {
			log.Println("error group update context")
			return
		}

		if err = i.groupInteractor.MarkGroupAsDeleted(user.Id, &groupId); err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, "Success", http.StatusOK)
	}
}

func (i *groupHanlder) ListGroups() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := utils.LoggedInUser(r)
		if user == nil {
			log.Printf("error group list user context")
			return
		}

		groups, err := i.groupInteractor.ListGroups(user.Id)
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, groups, http.StatusOK)
	}
}

func (i *groupHanlder) CreateSlice() http.HandlerFunc {
	type request struct {
		Name       string `json:"name"`
		ParentPath string `json:"parentPath"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var s request
		var err error

		vars := mux.Vars(r)
		groupIdArg, err := strconv.Atoi(vars["groupId"])
		if err != nil {
			utils.SendJsonError(w, "Invalid group id", http.StatusBadRequest)
			return
		}
		groupId := valueobject.ID(groupIdArg)

		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			utils.SendJsonError(w, "Invalid request data", http.StatusBadRequest)
			return
		}

		folder := &app.Slice{
			Name: s.Name,
			Path: s.ParentPath,
		}

		folder, err = i.groupInteractor.CreateSlice((*valueobject.ID)(&groupId), folder)
		if err != nil {
			utils.SendJsonError(w, "Create slice error", http.StatusBadRequest)
			return
		}

		utils.SendJson(w, folder, http.StatusOK)
	}
}

func (i *groupHanlder) ListSlices() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		groupIdArg, err := strconv.Atoi(vars["groupId"])
		if err != nil {
			utils.SendJsonError(w, "Invalid group id", http.StatusBadRequest)
			return
		}
		groupId := valueobject.ID(groupIdArg)

		slices, err := i.groupInteractor.ListSlices((*valueobject.ID)(&groupId))
		if err != nil {
			utils.SendJsonError(w, "List slice error", http.StatusBadRequest)
			return
		}

		utils.SendJson(w, slices, http.StatusOK)
	}
}

func (i *groupHanlder) InviteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		groupIdArg, err := strconv.Atoi(vars["groupId"])
		if err != nil {
			utils.SendJsonError(w, "Invalid group id", http.StatusBadRequest)
			return
		}
		groupId := valueobject.ID(groupIdArg)

		userIdArg, err := strconv.Atoi(vars["userId"])
		if err != nil {
			utils.SendJsonError(w, "Invalid user id", http.StatusBadRequest)
			return
		}
		userId := valueobject.ID(userIdArg)

		user := utils.LoggedInUser(r)
		if user == nil {
			log.Printf("error group list user context")
			return
		}

		err = i.groupInteractor.InviteUser(&groupId, &userId)
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, "Success", http.StatusOK)
	}
}

func (i *groupHanlder) ConfirmInvitation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		token := vars["token"]

		user := utils.LoggedInUser(r)
		if user == nil {
			utils.SendJsonError(w, "You don't have permissions to confirm invitation.", http.StatusBadRequest)
			return
		}

		err := i.groupInteractor.ConfirmInvitation(user.Id, token)
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, "Success", http.StatusOK)
	}
}

func (i *groupHanlder) DetachMember() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		groupIdArg, err := strconv.Atoi(vars["groupId"])
		if err != nil {
			utils.SendJsonError(w, "Invalid group id", http.StatusBadRequest)
			return
		}
		groupId := valueobject.ID(groupIdArg)

		memberIdArg, err := strconv.Atoi(vars["memberId"])
		if err != nil {
			utils.SendJsonError(w, "Invalid member id", http.StatusBadRequest)
			return
		}
		memberId := valueobject.ID(memberIdArg)

		user := utils.LoggedInUser(r)
		if user == nil {
			log.Printf("error group list user context")
			return
		}

		err = i.groupInteractor.DetachMember(&groupId, &memberId)
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, "Success", http.StatusOK)
	}
}
