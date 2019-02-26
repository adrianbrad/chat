package repository

import (
	"database/sql"
	"log"

	"github.com/adrianbrad/chat/model"
)

type dbRoomsRepository struct {
	db                 *sql.DB
	getOneQuery        string
	getAllQuery        string
	createQuery        string
	checkIfExistsQuery string
}

func NewDbRoomsRepository(database *sql.DB) Repository {
	return &dbRoomsRepository{
		db:                 database,
		getOneQuery:        getOneQuery("Room", "RoomID", "Name", "Description", "ChannelID"),
		getAllQuery:        getAllQuery("Room", "RoomID", "Name", "Description", "ChannelID"),
		createQuery:        createOneQuery("Room", "Name", "Description"),
		checkIfExistsQuery: checkIfExistsQuery("Room"),
	}
}

func (r dbRoomsRepository) GetOne(id int) (interface{}, error) {
	var room model.Room
	err := r.db.QueryRow(r.checkIfExistsQuery, id).Scan(
		&room.ID,
		&room.Name,
		&room.Description,
		&room.ChannelID)
	if err != nil {
		log.Println("Error while fetching room with id", id)
		return room, err
	}
	return room, nil
}

func (r dbRoomsRepository) GetAll() (rooms []interface{}) {
	rows, err := r.db.Query(r.getAllQuery)
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
			&room.Description,
			&room.ChannelID)
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

func (r dbRoomsRepository) Create(roomI interface{}) (id int, err error) {
	room := roomI.(model.Room)
	if err := r.db.QueryRow(r.createQuery, room.Name, room.Description, room.ChannelID).Scan(&id); err != nil {
		return id, err
	}
	return id, nil
}

func (r dbRoomsRepository) CheckIfExists(roomID int) (exists bool) {
	_ = r.db.QueryRow(checkIfExistsQuery("Room"), roomID).Scan(&exists)
	return exists
}
