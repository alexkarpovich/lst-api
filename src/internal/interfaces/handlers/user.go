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
