package repository

import (
	"database/sql"
	"log"

	"github.com/adrianbrad/chat/model"
)

type PermissionRepository interface {
	GetOne(int) (model.Permission, error)
	GetAll() []model.Permission
	Create(model.Permission) (int, error)
}

type dbPermissionsRepository struct {
	db *sql.DB
}

func NewDbPermissionsRepository(database *sql.DB) PermissionRepository {
	return &dbPermissionsRepository{db: database}
}

func (r dbPermissionsRepository) GetOne(id int) (permission model.Permission, err error) {
	err = r.db.QueryRow(`
	SELECT 
		"PermissionID", "Name", "Description"
	FROM "Permissions"
	WHERE "PermissionID"=$1
	`, id).Scan(
		&permission.ID,
		&permission.Name,
		&permission.Description)
	if err != nil {
		log.Println("Error while fetching permission with id", id)
		return permission, err
	}
	return permission, nil
}

func (r dbPermissionsRepository) GetAll() (permissions []model.Permission) {
	rows, err := r.db.Query(`
	SELECT
		"PermissionID", "Name", "Description"
	FROM "Permissions"`)
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

func (r dbPermissionsRepository) Create(permission model.Permission) (id int, err error) {
	if err := r.db.QueryRow(`
	INSERT INTO "Permissions"
		("Name", "Description")
	VALUES ($1, $2)
	RETURNING "PermissionID"`, permission.Name, permission.Description).Scan(&id); err != nil {
		return id, err
	}
	return id, nil
}
