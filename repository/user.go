package repository

import (
	"database/sql"
	"log"

	"github.com/adrianbrad/chat/model"
)

type UsersRepository interface {
	GetOne(int) (model.User, error)
	GetAll() []model.User
	Create(model.User) (int, error)
}

type dbUsersRepository struct {
	db *sql.DB
}

func NewDbUsersRepository(database *sql.DB) UsersRepository {
	return &dbUsersRepository{db: database}
}

func (r dbUsersRepository) GetOne(id int) (user model.User, err error) {
	err = r.db.QueryRow(`
	SELECT 
		"UserID", "Name", "RoleID", "UserData"
	FROM "Users"
	WHERE "UserID"=$1
	`, id).Scan(
		&user.ID,
		&user.Name,
		&user.RoleID,
		&user.UserData)
	if err != nil {
		log.Println("Error while fetching user with id", id)
		return user, err
	}
	return user, nil
}

func (r dbUsersRepository) GetAll() (users []model.User) {
	rows, err := r.db.Query(`
	SELECT
		"UserID", "Name", "RoleID", "UserData"
	FROM "Users"`)
	if err != nil {
		log.Println("Query error: ", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		user := model.User{}
		err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.RoleID,
			&user.UserData)
		if err != nil {
			log.Println("Mapping error", err)
			return
		}
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		log.Println("Reading rows error:", err)
	}
	return users
}

func (r dbUsersRepository) Create(user model.User) (id int, err error) {
	if err := r.db.QueryRow(`
	INSERT INTO "Users"
		("Name", "UserData", "RoleID")
	VALUES ($1, $2, $3)
	RETURNING "UserID"`, user.Name, user.UserData, user.RoleID).Scan(&id); err != nil {
		return id, err
	}
	return id, nil
}
