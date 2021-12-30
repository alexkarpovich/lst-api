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
	CreateGroup(*app.Group) (*app.Group, error)
	ListGroups(*valueobject.ID) ([]*app.Group, error)
	CreateSlice(*valueobject.ID, *app.Slice) (*app.Slice, error)
	ListSlices(*valueobject.ID) ([]*app.NestedSlice, error)
	AttachMember(*valueobject.ID, *valueobject.ID, app.UserRole) error
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
	h.router.HandleFunc("/me/group/{groupId}/slice", h.CreateSlice()).Methods("POST")
	h.router.HandleFunc("/me/group/{groupId}/slice", h.ListSlices()).Methods("GET")
	h.router.HandleFunc("/me/group/{groupId}/attach-member", h.AttachMember()).Methods("POST")
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

		if group, err = i.groupInteractor.CreateGroup(group); err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, group, http.StatusOK)
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
		profileId := valueobject.ID(groupIdArg)

		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			utils.SendJsonError(w, "Invalid request data", http.StatusBadRequest)
			return
		}

		folder := &app.Slice{
			Name: s.Name,
			Path: s.ParentPath,
		}

		folder, err = i.groupInteractor.CreateSlice((*valueobject.ID)(&profileId), folder)
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

func (i *groupHanlder) AttachMember() http.HandlerFunc {
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

		err = i.groupInteractor.AttachMember(&groupId, s.MemberId, s.Role)
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
