package handlers

import (
	"net/http"
	"strconv"

	"github.com/alexkarpovich/lst-api/src/internal/domain"
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
	"github.com/alexkarpovich/lst-api/src/internal/utils"
	"github.com/gorilla/mux"
)

type ExpressionInteractor interface {
	Search(string, string) ([]*domain.Expression, error)
	GetTranscriptionParts(*valueobject.ID, *valueobject.ID) (map[string][]*domain.TranscriptionItem, error)
}

type expressionHanlder struct {
	BaseHanlder
	expressionInteractor ExpressionInteractor
}

func ConfigureExpressionHandler(ei ExpressionInteractor, r *mux.Router) {
	h := &expressionHanlder{
		BaseHanlder: BaseHanlder{
			router: r,
		},
		expressionInteractor: ei,
	}

	h.router.
		HandleFunc("/x", h.Search()).
		Queries("lang", "{[a-z]{2}}").
		Queries("search", "{.+}").
		Methods("GET")
	h.router.
		HandleFunc("/x/{expression_id}/transcription-parts", h.GetTranscriptionParts()).
		Queries("type", "{\\d+}").
		Methods("GET")
	// h.router.HandleFunc("/me/group/{groupId}/slice", h.CreateSlice()).Methods("POST")
	// h.router.HandleFunc("/me/group/{groupId}/slice", h.ListSlices()).Methods("GET")
}

func (i *expressionHanlder) Search() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		langCode := r.FormValue("lang")
		searchValue := r.FormValue("search")

		expressions, err := i.expressionInteractor.Search(langCode, searchValue)
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, expressions, http.StatusOK)
	}
}

func (i *expressionHanlder) GetTranscriptionParts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		typeIdArg, err := strconv.Atoi(r.FormValue("type"))
		if err != nil {
			utils.SendJsonError(w, "Invalid type id", http.StatusBadRequest)
			return
		}
		typeId := valueobject.ID(typeIdArg)

		vars := mux.Vars(r)
		expressionIdArg, err := strconv.Atoi(vars["expression_id"])
		if err != nil {
			utils.SendJsonError(w, "Invalid group id", http.StatusBadRequest)
			return
		}
		expressionId := valueobject.ID(expressionIdArg)

		parts, err := i.expressionInteractor.GetTranscriptionParts(&expressionId, &typeId)
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, parts, http.StatusOK)
	}
}

// func (i *groupHanlder) CreateGroup() http.HandlerFunc {
// 	type request struct {
// 		Name         string          `json:"name"`
// 		TargetLangId *valueobject.ID `json:"targetLangId"`
// 		NativeLangId *valueobject.ID `json:"nativeLangId"`
// 	}

// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var s request
// 		var err error

// 		if err = json.NewDecoder(r.Body).Decode(&s); err != nil {
// 			utils.SendJsonError(w, "Invalid request data", http.StatusBadRequest)
// 			return
// 		}

// 		user := utils.LoggedInUser(r)
// 		if user == nil {
// 			log.Println("error profile create context")
// 			return
// 		}

// 		group := &app.Group{
// 			Name:         s.Name,
// 			TargetLangId: s.TargetLangId,
// 			NativeLangId: s.NativeLangId,
// 		}

// 		if group, err = i.groupInteractor.CreateGroup(group); err != nil {
// 			utils.SendJsonError(w, err, http.StatusBadRequest)
// 			return
// 		}

// 		utils.SendJson(w, group, http.StatusOK)
// 	}
// }

// func (i *groupHanlder) ListGroups() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		user := utils.LoggedInUser(r)
// 		if user == nil {
// 			log.Printf("error group list user context")
// 			return
// 		}

// 		groups, err := i.groupInteractor.ListGroups(user.Id)
// 		if err != nil {
// 			utils.SendJsonError(w, err, http.StatusBadRequest)
// 			return
// 		}

// 		utils.SendJson(w, groups, http.StatusOK)
// 	}
// }

// func (i *groupHanlder) CreateSlice() http.HandlerFunc {
// 	type request struct {
// 		Name       string `json:"name"`
// 		ParentPath string `json:"parentPath"`
// 	}

// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var s request
// 		var err error

// 		vars := mux.Vars(r)
// 		groupIdArg, err := strconv.Atoi(vars["groupId"])
// 		if err != nil {
// 			utils.SendJsonError(w, "Invalid group id", http.StatusBadRequest)
// 			return
// 		}
// 		profileId := valueobject.ID(groupIdArg)

// 		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
// 			utils.SendJsonError(w, "Invalid request data", http.StatusBadRequest)
// 			return
// 		}

// 		folder := &app.Node{
// 			Name: s.Name,
// 			Path: s.ParentPath,
// 		}

// 		folder, err = i.groupInteractor.CreateSlice((*valueobject.ID)(&profileId), folder)
// 		if err != nil {
// 			utils.SendJsonError(w, "Create folder error", http.StatusBadRequest)
// 			return
// 		}

// 		utils.SendJson(w, folder, http.StatusOK)
// 	}
// }

// func (i *groupHanlder) ListSlices() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		vars := mux.Vars(r)
// 		groupIdArg, err := strconv.Atoi(vars["groupId"])
// 		if err != nil {
// 			utils.SendJsonError(w, "Invalid profile id", http.StatusBadRequest)
// 			return
// 		}
// 		groupId := valueobject.ID(groupIdArg)

// 		slices, err := i.groupInteractor.ListSlices((*valueobject.ID)(&groupId))
// 		if err != nil {
// 			utils.SendJsonError(w, "List folder error", http.StatusBadRequest)
// 			return
// 		}

// 		utils.SendJson(w, slices, http.StatusOK)
// 	}
// }
