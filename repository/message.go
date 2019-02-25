package repository

import (
	"database/sql"
	"log"

	"github.com/adrianbrad/chat/model"
)

type MessagesRepository interface {
	GetOne(int) (model.Message, error)
	GetAll() []model.Message
	Create(model.Message) (int, error)
}

type dbMessagesRepository struct {
	db *sql.DB
}

func NewDbMessagesRepository(database *sql.DB) MessagesRepository {
	return &dbMessagesRepository{db: database}
}

func (r dbMessagesRepository) GetOne(id int) (message model.Message, err error) {
	err = r.db.QueryRow(`
	SELECT 
		"MessageID", "Content", "RoomID", "UserID"
	FROM "Messages"
	WHERE "MessageID"=$1
	`, id).Scan(
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

func (r dbMessagesRepository) GetAll() (messages []model.Message) {
	rows, err := r.db.Query(`
	SELECT
		"MessageID", "Content", "RoomID", "UserID"
	FROM "Messages"`)
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

func (r dbMessagesRepository) Create(message model.Message) (id int, err error) {
	if err := r.db.QueryRow(`
	INSERT INTO "Messages"
		("Content", "UserID", "RoomID")
	VALUES ($1, $2, $3)
	RETURNING "MessageID"`, message.Content, message.UserID, message.RoomID).Scan(&id); err != nil {
		return id, err
	}
	return id, nil
}
