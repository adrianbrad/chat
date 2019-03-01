package repository

import (
	"database/sql"
	"log"

	"github.com/adrianbrad/chat/model"
)

type dbMessagesRepository struct {
	db                     *sql.DB
	getOneQuery            string
	getAllQuery            string
	createQuery            string
	checkIfExistsQuery     string
	getAllWhereRoomIDQuery string
}

func NewDbMessagesRepository(database *sql.DB) Repository {
	return &dbMessagesRepository{
		db:                     database,
		getOneQuery:            getOneQuery("Message", "MessageID", "Content", "RoomID", "UserID"),
		getAllQuery:            getAllQuery("Message", "MessageID", "Content", "RoomID", "UserID"),
		createQuery:            createOneQuery("Message", "Content", "RoomID", "UserID"),
		checkIfExistsQuery:     checkIfExistsQuery("Message"),
		getAllWhereRoomIDQuery: getAllWhereQuery("Message", "RoomID", "CreatedAt", "desc", "*"),
	}
}

func (r dbMessagesRepository) GetOne(id int) (interface{}, error) {
	var message model.Message
	err := r.db.QueryRow(r.getOneQuery, id).Scan(
		&message.ID,
		&message.Content,
		&message.RoomID,
		&message.UserID)
	if err != nil {
		log.Println("Error while fetching message with id", id)
		return message, err
	}
	return message, nil
}

func (r dbMessagesRepository) GetAll() (messages []interface{}) {
	rows, err := r.db.Query(r.getAllQuery)
	if err != nil {
		log.Println("Query error: ", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		message := model.Message{}
		err = rows.Scan(
			&message.ID,
			&message.Content,
			&message.RoomID,
			&message.UserID)
		if err != nil {
			log.Println("Mapping error", err)
			return
		}
		messages = append(messages, message)
	}
	err = rows.Err()
	if err != nil {
		log.Println("Reading rows error:", err)
	}
	return messages
}

func (r dbMessagesRepository) GetAllWhere(cloumn string, value int, limit int) (messages []interface{}) {
	rows, err := r.db.Query(r.getAllWhereRoomIDQuery, value, limit)
	if err != nil {
		log.Println("Query error: ", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		message := model.Message{}
		err = rows.Scan(
			&message.ID,
			&message.Content,
			&message.RoomID,
			&message.UserID,
			&message.CreatedAt)
		if err != nil {
			log.Println("Mapping error", err)
			return
		}
		messages = append(messages, message)
	}
	err = rows.Err()
	if err != nil {
		log.Println("Reading rows error:", err)
	}
	return messages
}

func (r dbMessagesRepository) Create(messageI interface{}) (id int, err error) {
	message := messageI.(model.Message)
	if err := r.db.QueryRow(r.createQuery, message.Content, message.UserID, message.RoomID).Scan(&id); err != nil {
		return id, err
	}
	return id, nil
}

func (r dbMessagesRepository) CheckIfExists(messageID int) (exists bool) {
	_ = r.db.QueryRow(r.checkIfExistsQuery, messageID).Scan(&exists)
	return exists
}
