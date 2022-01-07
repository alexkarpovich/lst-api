package handlers

import (
	"log"
	"net/http"

	"github.com/alexkarpovich/lst-api/src/internal/app"
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
	"github.com/alexkarpovich/lst-api/src/internal/utils"
	"github.com/gorilla/mux"
)

type UserInteractor interface {
	Get(*valueobject.ID) (*app.User, error)
	FindByUsername(string) (*app.User, error)
}

type userHandler struct {
	BaseHanlder
	userInteractor UserInteractor
}

func ConfigureUserHandler(ui UserInteractor, r *mux.Router) {
	h := &userHandler{
		BaseHanlder: BaseHanlder{
			router: r,
		},
		userInteractor: ui,
	}

	h.router.HandleFunc("/me", h.Get()).Methods("GET")
	h.router.HandleFunc("/users/{username:[a-zA-Z0..9]+}", h.FindByUsername()).Methods("GET")
	// h.router.HandleFunc("/me/group", h.ListGroups()).Methods("GET")
	// h.router.HandleFunc("/me/group/{groupId}/slice", h.CreateSlice()).Methods("POST")
	// h.router.HandleFunc("/me/group/{groupId}/slice", h.ListSlices()).Methods("GET")
}

func (i *userHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := utils.LoggedInUser(r)
		if user == nil {
			log.Println("error user get context")
			return
		}

		utils.SendJson(w, user, http.StatusOK)
	}
}

func (i *userHandler) FindByUsername() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]

		user := utils.LoggedInUser(r)
		if user == nil {
			utils.SendJsonError(w, "You don't have permissions.", http.StatusBadRequest)
			return
		}

		usr, err := i.userInteractor.FindByUsername(username)
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, usr, http.StatusOK)
	}
}
