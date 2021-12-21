package repos

import (
	"log"

	"github.com/alexkarpovich/lst-api/src/internal/app"
	"github.com/alexkarpovich/lst-api/src/internal/domain/valueobject"
	"github.com/alexkarpovich/lst-api/src/internal/interfaces/db"
)

type UserRepo struct {
	db db.DB
}

func NewUserRepo(db db.DB) *UserRepo {
	return &UserRepo{db}
}

func (r *UserRepo) Create(obj *app.User) (*app.User, error) {
	stmt := `
		INSERT INTO users (username, email, encrypted_password, first_name, last_name, token, token_expires_at, status) 
		VALUES(:username, :email, :encrypted_password, :first_name, :last_name, :token, :token_expires_at, :status)`

	rows, err := r.db.Db().NamedQuery(stmt, obj)

	if err != nil {
		return nil, err
	}

	if rows.Next() {
		rows.Scan(obj.Id)
	}

	if err = rows.Close(); err != nil {
		// but what should we do if there's an error?
		log.Println(err)
	}

	return obj, nil
}

func (r *UserRepo) Get(userId *valueobject.ID) (*app.User, error) {
	stmt := `SELECT * FROM users WHERE id=$1`
	user := &app.User{}

	if err := r.db.Db().Get(user, stmt, userId); err != nil {
		log.Println(err)
		return nil, err
	}

	return user, nil
}

func (r *UserRepo) FindByEmail(email valueobject.EmailAddress) (*app.User, error) {
	stmt := `SELECT * FROM users WHERE email=$1`
	user := &app.User{}

	if err := r.db.Db().Get(user, stmt, email); err != nil {
		log.Println(err)
		return nil, err
	}

	return user, nil
}

func (r *UserRepo) FindByToken(token string) (*app.User, error) {
	stmt := `
		SELECT * FROM users 
		WHERE token=$1 AND token_expires_at > NOW() AND status=$2`
	user := &app.User{}

	if err := r.db.Db().Get(user, stmt, token, app.UserUnconfirmed); err != nil {
		log.Println(err)
		return nil, err
	}

	return user, nil
}

func (r *UserRepo) Delete(userId *valueobject.ID) error {
	stmt := `UPDATE users SET status=$1 WHERE id=$2`

	if _, err := r.db.Db().Exec(stmt, app.UserDeleted, userId); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (r *UserRepo) Update(obj *app.User) error {
	stmt := `
		UPDATE users SET username=:username, email=:email, encrypted_password=:encrypted_password, 
						 first_name=:first_name, last_name=:last_name, status=:status`

	if _, err := r.db.Db().NamedExec(stmt, obj); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
