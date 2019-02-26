package repository

import (
	"database/sql"
	"log"

	"github.com/adrianbrad/chat/model"
)

type dbUsersRepository struct {
	db                 *sql.DB
	getOneQuery        string
	getAllQuery        string
	createQuery        string
	checkIfExistsQuery string
}

func NewDbUsersRepository(database *sql.DB) Repository {
	return &dbUsersRepository{
		db:                 database,
		getOneQuery:        getOneQuery("User", "UserID", "Name", "RoleID", "UserData"),
		getAllQuery:        getAllQuery("User", "UserID", "Name", "RoleID", "UserData"),
		createQuery:        createOneQuery("User", "Name", "UserData", "RoleID"),
		checkIfExistsQuery: checkIfExistsQuery("User"),
	}
}

func (r dbUsersRepository) GetOne(id int) (interface{}, error) {
	var user model.User
	err := r.db.QueryRow(r.getOneQuery, id).Scan(
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

func (r dbUsersRepository) CheckIfExists(userID int) (exists bool) {
	_ = r.db.QueryRow(r.checkIfExistsQuery, userID).Scan(&exists)
	return exists
}

func (r dbUsersRepository) GetAll() (users []interface{}) {
	rows, err := r.db.Query(r.getAllQuery)
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

func (r dbUsersRepository) Create(userI interface{}) (id int, err error) {
	user := userI.(model.User)
	if err := r.db.QueryRow(
		r.createQuery, user.Name, user.UserData, user.RoleID).
		Scan(&id); err != nil {
		return id, err
	}
	return id, nil
}
