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
	CreateGroup(*valueobject.ID, app.Group) (*app.Group, error)
	UpdateGroup(*valueobject.ID, app.Group) error
	ListGroups(*valueobject.ID) ([]*app.Group, error)
	MarkGroupAsDeleted(*valueobject.ID, *valueobject.ID) error
	CreateSlice(*valueobject.ID, app.Node) (*app.Node, error)
	ListSlices(*valueobject.ID) ([]*app.FlatNode, error)
	InviteUser(*valueobject.ID, *valueobject.ID, *valueobject.ID) error
	ConfirmInvitation(*valueobject.ID, string) error
	DetachMember(*valueobject.ID, *valueobject.ID) error
	UpdateMemberRole(*valueobject.ID, *valueobject.ID, app.GroupMember) error
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

	h.router.HandleFunc("/me/groups", h.CreateGroup()).Methods("POST")
	h.router.HandleFunc("/me/groups", h.ListGroups()).Methods("GET")
	h.router.HandleFunc("/me/groups/confirm-invitation/{token}", h.ConfirmInvitation()).Methods("POST")
	h.router.HandleFunc("/me/groups/{groupId}", h.UpdateGroup()).Methods("POST")
	h.router.HandleFunc("/me/groups/{groupId}", h.DeleteGroup()).Methods("DELETE")
	h.router.HandleFunc("/me/groups/{groupId}/nodes", h.CreateSlice()).Methods("POST")
	h.router.HandleFunc("/me/groups/{groupId}/nodes", h.ListSlices()).Methods("GET")
	h.router.HandleFunc("/me/groups/{groupId}/invite-user/{userId}", h.InviteUser()).Methods("POST")
	h.router.HandleFunc("/me/groups/{groupId}/detach-member/{memberId}", h.DetachMember()).Methods("POST")
	h.router.HandleFunc("/me/groups/{groupId}/update-role", h.UpdateMemberRole()).Methods("POST")
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

		inGroup := app.Group{
			Name:           s.Name,
			TargetLangCode: s.TargetLangCode,
			NativeLangCode: s.NativeLangCode,
		}

		group, err := i.groupInteractor.CreateGroup(user.Id, inGroup)
		if err != nil {
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

		inGroup := app.Group{
			Id:             &groupId,
			Name:           s.Name,
			TargetLangCode: s.TargetLangCode,
			NativeLangCode: s.NativeLangCode,
		}

		if err = i.groupInteractor.UpdateGroup(user.Id, inGroup); err != nil {
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
		Type       app.NodeType `json:"type"`
		Name       string       `json:"name"`
		ParentPath string       `json:"parentPath"`
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

		inNode := app.Node{
			Type: s.Type,
			Name: s.Name,
			Path: s.ParentPath,
		}

		node, err := i.groupInteractor.CreateSlice((*valueobject.ID)(&groupId), inNode)
		if err != nil {
			utils.SendJsonError(w, "Create slice error", http.StatusBadRequest)
			return
		}

		utils.SendJson(w, node, http.StatusOK)
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

		err = i.groupInteractor.InviteUser(user.Id, &groupId, &userId)
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

func (i *groupHanlder) UpdateMemberRole() http.HandlerFunc {
	type request struct {
		MemberId *valueobject.ID `json:"memberId"`
		Role     app.UserRole    `json:"role"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var s request

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

		user := utils.LoggedInUser(r)
		if user == nil {
			log.Printf("error group list user context")
			return
		}

		member := app.GroupMember{
			Id:   s.MemberId,
			Role: s.Role,
		}

		err = i.groupInteractor.UpdateMemberRole(user.Id, &groupId, member)
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, "Success", http.StatusOK)
	}
}
