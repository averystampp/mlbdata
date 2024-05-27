package mlb

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/averystampp/sesame"
	"github.com/google/uuid"
	bolt "go.etcd.io/bbolt"
	"golang.org/x/crypto/bcrypt"
)

// bolt bucket schema for users string(uuid):User{}
// bolt bucket schema csrf string(uuid):int64(expires)
type User struct {
	Name     string
	Password string
}

func UserService(rtr *sesame.Router) {
	rtr.Get("/login", LoginView)
	rtr.Post("/login", Login)

	rtr.Get("/user/new", CreateUserView)
	rtr.Post("/user/new", CreateUser)

}

func Login(ctx sesame.Context) error {
	csrfToken := ctx.Request().PostFormValue("_csrf")
	username := ctx.Request().PostFormValue("username")
	password := ctx.Request().PostFormValue("password")
	if csrfToken == "" {
		http.Error(ctx.Response(), "missing token", http.StatusInternalServerError)
		return nil
	}

	if username == "" || password == "" {
		return fmt.Errorf("missing username or password")
	}

	db, err := bolt.Open("../players.db", 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("csrf"))
		token := b.Get([]byte(csrfToken))
		exp, err := strconv.ParseInt(string(token), 10, 64)
		if err != nil {
			return err
		}
		if time.Now().UnixMilli() > exp {
			return fmt.Errorf("session expired please retry logging in")
		}

		b = tx.Bucket([]byte("users"))
		userFromDB := b.Get([]byte(username))
		if string(userFromDB) == "" {
			return fmt.Errorf("user not in database")
		}
		var user User
		err = json.Unmarshal(userFromDB, &user)
		if err != nil {
			return err
		}

		return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	})
	if err != nil {
		return err
	}

	http.Redirect(ctx.Response(), ctx.Request(), "/", http.StatusSeeOther)
	return nil
}

func LoginView(ctx sesame.Context) error {
	db, err := bolt.Open("../players.db", 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()
	tmpl, err := template.ParseFiles("../pages/login.html")
	if err != nil {
		return err
	}
	id := uuid.New().String()
	expires := time.Now().Add(time.Minute * 30).UnixMilli()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("csrf"))
		return b.Put([]byte(id), []byte(strconv.FormatInt(expires, 10)))
	})
	if err != nil {
		return err
	}
	return tmpl.Execute(ctx.Response(), id)
}

func CreateUserView(ctx sesame.Context) error {
	tmpl, err := template.ParseFiles("../pages/newUser.html")
	if err != nil {
		return err
	}

	return tmpl.Execute(ctx.Response(), nil)
}

func CreateUser(ctx sesame.Context) error {
	username := ctx.Request().PostFormValue("username")
	password := ctx.Request().PostFormValue("password")
	if username == "" || password == "" {
		return fmt.Errorf("must have a username and password")
	}

	db, err := bolt.Open("../players.db", 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	b, err := json.Marshal(User{
		Name:     username,
		Password: string(hash),
	})
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		bucket.Put([]byte(username), b)
		return nil
	})

	if err != nil {
		return err
	}
	http.Redirect(ctx.Response(), ctx.Request(), "/", http.StatusSeeOther)
	return nil
}
