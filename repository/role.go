package repository

import (
	"database/sql"
	"log"

	"github.com/adrianbrad/chat/model"
)

type RoleRepository interface {
	GetOne(int) (model.Role, error)
	GetAll() []model.Role
	Create(model.Role) (int, error)
}

type dbRolesRepository struct {
	db *sql.DB
}

func NewDbRolesRepository(database *sql.DB) RoleRepository {
	return &dbRolesRepository{db: database}
}

func (r dbRolesRepository) GetOne(id int) (role model.Role, err error) {
	err = r.db.QueryRow(`
	SELECT 
		"RoleID", "Name", "Description"
	FROM "Roles"
	WHERE "RoleID"=$1
	`, id).Scan(
		&role.ID,
		&role.Name,
		&role.Description)
	if err != nil {
		log.Println("Error while fetching role with id", id)
		return role, err
	}
	return role, nil
}

func (r dbRolesRepository) GetAll() (roles []model.Role) {
	rows, err := r.db.Query(`
	SELECT
		"RoleID", "Name", "Description"
	FROM "Roles"`)
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

func (r dbRolesRepository) Create(role model.Role) (id int, err error) {
	if err := r.db.QueryRow(`
	INSERT INTO "Roles"
		("Name", "Description")
	VALUES ($1, $2)
	RETURNING "RoleID"`, role.Name, role.Description).Scan(&id); err != nil {
		return id, err
	}
	return id, nil
}
