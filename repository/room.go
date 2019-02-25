package repository

import (
	"database/sql"
	"log"

	"github.com/adrianbrad/chat/model"
)

type RoomRepository interface {
	GetOne(int) (model.Room, error)
	GetAll() []model.Room
	Create(model.Room) (int, error)
}

type dbRoomsRepository struct {
	db *sql.DB
}

func NewDbRoomsRepository(database *sql.DB) RoomRepository {
	return &dbRoomsRepository{db: database}
}

func (r dbRoomsRepository) GetOne(id int) (room model.Room, err error) {
	err = r.db.QueryRow(`
	SELECT 
		"RoomID", "Name", "Description"
	FROM "Rooms"
	WHERE "RoomID"=$1
	`, id).Scan(
		&room.ID,
		&room.Name,
		&room.Description)
	if err != nil {
		log.Println("Error while fetching room with id", id)
		return room, err
	}
	return room, nil
}

func (r dbRoomsRepository) GetAll() (rooms []model.Room) {
	rows, err := r.db.Query(`
	SELECT
		"RoomID", "Name", "Description"
	FROM "Rooms"`)
	if err != nil {
		log.Println("Query error: ", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		room := model.Room{}
		err = rows.Scan(
			&room.ID,
			&room.Name,
			&room.Description)
		if err != nil {
			log.Println("Mapping error", err)
			return
		}
		rooms = append(rooms, room)
	}
	err = rows.Err()
	if err != nil {
		log.Println("Reading rows error:", err)
	}
	return rooms
}

func (r dbRoomsRepository) Create(room model.Room) (id int, err error) {
	if err := r.db.QueryRow(`
	INSERT INTO "Rooms"
		("Name", "Description")
	VALUES ($1, $2)
	RETURNING "RoomID"`, room.Name, room.Description).Scan(&id); err != nil {
		return id, err
	}
	return id, nil
}
