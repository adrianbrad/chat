package repository

import (
	"database/sql"
	"log"

	"github.com/adrianbrad/chat/model"
)

type dbRolesRepository struct {
	db                 *sql.DB
	getOneQuery        string
	getAllQuery        string
	createQuery        string
	checkIfExistsQuery string
}

// func NewDbRolesRepository(database *sql.DB) Repository {
// 	return &dbRolesRepository{
// 		db:                 database,
// 		getOneQuery:        getOneQuery("Role", "RoleID", "Name", "Description"),
// 		getAllQuery:        getAllQuery("Role", "RoleID", "Name", "Description"),
// 		createQuery:        createOneQuery("Role", "Name", "Description"),
// 		checkIfExistsQuery: checkIfExistsQuery("Role"),
// 	}
// }

func (r dbRolesRepository) GetOne(id int) (interface{}, error) {
	var role model.Role
	err := r.db.QueryRow(r.getOneQuery, id).Scan(
		&role.ID,
		&role.Name,
		&role.Description)
	if err != nil {
		log.Println("Error while fetching role with id", id)
		return role, err
	}
	return role, nil
}

func (r dbRolesRepository) GetAll() (roles []interface{}) {
	rows, err := r.db.Query(r.getAllQuery)
	if err != nil {
		log.Println("Query error: ", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		role := model.Role{}
		err = rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description)
		if err != nil {
			log.Println("Mapping error", err)
			return
		}
		roles = append(roles, role)
	}
	err = rows.Err()
	if err != nil {
		log.Println("Reading rows error:", err)
	}
	return roles
}

func (r dbRolesRepository) GetAllWhere(cloumn string, value int, limit int) []interface{} {

	return nil
}

func (r dbRolesRepository) Create(roleI interface{}) (id int, err error) {
	role := roleI.(model.Role)
	if err := r.db.QueryRow(
		r.createQuery, role.Name, role.Description).
		Scan(&id); err != nil {
		return id, err
	}
	return id, nil
}

func (r dbRolesRepository) CheckIfExists(roleID int) (exists bool) {
	_ = r.db.QueryRow(r.checkIfExistsQuery, roleID).Scan(&exists)
	return exists
}
