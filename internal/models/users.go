package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type userModel struct {
	DB *sql.DB
}

type UserModel interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
}

func NewUserModel(db *sql.DB) UserModel {
	return &userModel{DB: db}
}

func (m *userModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (name, email, hashed_password, created)
	VALUE(?, ?, ?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(query, name, email, hashedPassword)
	if err != nil {
		var mySQLErr *mysql.MySQLError
		if errors.As(err, &mySQLErr) {
			if mySQLErr.Number == 1062 && strings.Contains(mySQLErr.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (m *userModel) Authenticate(email, password string) (int, error) {
	query := `SELECT id, hashed_password FROM users WHERE email = ?`

	var user User
	err := m.DB.QueryRow(query, email).Scan(&user.ID, &user.HashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	return user.ID, nil
}

func (m *userModel) Exists(id int) (bool, error) {
	query := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"

	var exists bool
	err := m.DB.QueryRow(query, id).Scan(&exists)
	return exists, err
}
