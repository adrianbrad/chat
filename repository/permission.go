package repository

import (
	"database/sql"
	"log"

	"github.com/adrianbrad/chat/model"
)

type dbPermissionsRepository struct {
	db                 *sql.DB
	getOneQuery        string
	getAllQuery        string
	createQuery        string
	checkIfExistsQuery string
}

func NewDbPermissionsRepository(database *sql.DB) Repository {
	return &dbPermissionsRepository{
		db:                 database,
		getOneQuery:        getOneQuery("Permission", "PermissionID", "Name", "Description"),
		getAllQuery:        getAllQuery("Permission", "PermissionID", "Name", "Description"),
		createQuery:        createOneQuery("Permission", "Name", "Description"),
		checkIfExistsQuery: checkIfExistsQuery("Permission"),
	}
}

func (r dbPermissionsRepository) GetOne(id int) (interface{}, error) {
	var permission model.Permission
	err := r.db.QueryRow(r.getOneQuery, id).Scan(
		&permission.ID,
		&permission.Name,
		&permission.Description)
	if err != nil {
		log.Println("Error while fetching permission with id", id)
		return permission, err
	}
	return permission, nil
}

func (r dbPermissionsRepository) GetAll() (permissions []interface{}) {
	rows, err := r.db.Query(r.getAllQuery)
	if err != nil {
		log.Println("Query error: ", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		permission := model.Permission{}
		err = rows.Scan(
			&permission.ID,
			&permission.Name,
			&permission.Description)
		if err != nil {
			log.Println("Mapping error", err)
			return
		}
		permissions = append(permissions, permission)
	}
	err = rows.Err()
	if err != nil {
		log.Println("Reading rows error:", err)
	}
	return permissions
}

func (r dbPermissionsRepository) GetAllWhere(cloumn string, value int, limit int) []interface{} {

	return nil
}

func (r dbPermissionsRepository) Create(permissionI interface{}) (id int, err error) {
	permission := permissionI.(model.Permission)
	if err := r.db.QueryRow(
		r.createQuery, permission.Name, permission.Description).
		Scan(&id); err != nil {
		return id, err
	}
	return id, nil
}

func (r dbPermissionsRepository) CheckIfExists(permissionID int) (exists bool) {
	_ = r.db.QueryRow(r.checkIfExistsQuery, permissionID).Scan(&exists)
	return exists
}
