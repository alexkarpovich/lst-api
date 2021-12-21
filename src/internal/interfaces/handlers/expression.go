package handlers

import (
	"github.com/gorilla/mux"
)

type ExpressionInteractor interface {
}

type expressionHanlder struct {
	BaseHanlder
	expressionInteractor ExpressionInteractor
}

func ConfigureExpressionHandler(ei ExpressionInteractor, r *mux.Router) {
	// h := &expressionHanlder{
	// 	BaseHanlder: BaseHanlder{
	// 		router: r,
	// 	},
	// 	expressionInteractor: ei,
	// }

	// h.router.HandleFunc("/", h.CreateGroup()).Methods("POST")
	// h.router.HandleFunc("/me/group", h.ListGroups()).Methods("GET")
	// h.router.HandleFunc("/me/group/{groupId}/slice", h.CreateSlice()).Methods("POST")
	// h.router.HandleFunc("/me/group/{groupId}/slice", h.ListSlices()).Methods("GET")
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

// 		folder := &app.Slice{
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
