package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/alexkarpovich/lst-api/src/internal/app"
	"github.com/alexkarpovich/lst-api/src/internal/app/usecases"
	"github.com/alexkarpovich/lst-api/src/internal/utils"
	"github.com/gorilla/mux"
)

type AuthInteractor interface {
	CreateRegistrant(*usecases.Registrant) (*app.User, error)
	ConfirmEmail(token string) error
	FindUserByEmailPassword(email string, password string) (*app.User, error)
}

type authHanlder struct {
	BaseHanlder
	authInteractor AuthInteractor
}

func ConfigureAuthHandler(ai AuthInteractor, r *mux.Router) {
	i := &authHanlder{
		BaseHanlder: BaseHanlder{
			router: r,
		},
		authInteractor: ai,
	}

	i.router.HandleFunc("/signup", i.Signup()).Methods("POST")
	i.router.HandleFunc("/signup/{token}", i.ConfirmSignup()).Methods("POST")
	i.router.HandleFunc("/login", i.Login()).Methods("POST")
}

func (i *authHanlder) Signup() http.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var s request

		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			utils.SendJsonError(w, "Invalid request data", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		registrant := &usecases.Registrant{
			Username: s.Username,
			Email:    s.Email,
			Password: s.Password,
		}

		if _, err := i.authInteractor.CreateRegistrant(registrant); err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		// TODO: send signup email

		utils.SendJson(w, "Success", http.StatusOK)
	}
}

func (i *authHanlder) ConfirmSignup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		if err := i.authInteractor.ConfirmEmail(vars["token"]); err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, "Success", http.StatusOK)
	}
}

func (i *authHanlder) Login() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var s request
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
		}
		defer r.Body.Close()

		user, err := i.authInteractor.FindUserByEmailPassword(s.Email, s.Password)
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		token := utils.NewToken(user.Id)

		utils.SendJson(w, map[string]string{"token": token}, http.StatusOK)
	}
}
