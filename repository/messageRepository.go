package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/adrianbrad/chat/model"
)

type MessageRepository interface {
	GetOne(int) (*model.Message, error)
	Create(*model.Message) (int, error)
	GetHistory(int, int) []*model.Message
}

type dbMessagesRepository struct {
	db                     *sql.DB
	getOneQuery            string
	getAllQuery            string
	createQuery            string
	checkIfExistsQuery     string
	getAllWhereRoomIDQuery string
}

func NewDbMessagesRepository(database *sql.DB) MessageRepository {
	return &dbMessagesRepository{
		db: database,
		getOneQuery: `
		SELECT "Messages"."MessageID", "Messages"."Content", "Messages"."SentAt", "Messages"."UserID",
			(SELECT array(SELECT "RoomID" FROM "Messages_Rooms" WHERE "Messages_Rooms"."MessageID" = "Messages"."MessageID")) AS "RoomIDs" 
		FROM "Messages"
		WHERE "Messages"."MessageID" = $1`,
		getAllQuery: `
		SELECT "Messages"."MessageID", "Messages"."Content", "Messages"."SentAt", "Messages"."UserID",
			(SELECT array(SELECT "RoomID" FROM "Messages_Rooms" WHERE "Messages_Rooms"."MessageID" = "Messages"."MessageID")) AS "RoomIDs" 
 		FROM "Messages"`,
		createQuery:            createOneQuery("Message", "Content", "UserID"),
		checkIfExistsQuery:     checkIfExistsQuery("Message"),
		getAllWhereRoomIDQuery: getAllWhereQuery("Message", "RoomID", "CreatedAt", "desc", "*"),
	}
}

func (r dbMessagesRepository) GetOne(id int) (*model.Message, error) {
	message := &model.Message{}

	err := r.db.QueryRow(r.getOneQuery, id).Scan(
		&message.ID,
		&message.Content,
		&message.SentAt,
		&message.UserID,
		&message.RoomIDs)
	if err != nil {
		log.Println("Error while fetching message with id", id)
		return nil, err
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
			&message.SentAt,
			&message.UserID,
			&message.RoomIDs,
		)

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

//TODO
//SELECT * ,
//(SELECT array(SELECT "RoomID" FROM "Messages_Rooms" WHERE "mr"."MessageID" = "m"."MessageID" AND "mr"."RoomID" = 2)) AS "RoomIDs"
//FROM "Messages" as "m", "Messages_Rooms" as "mr"
//WHERE "mr"."RoomID" = 1 AND "m"."MessageID" = "mr"."MessageID"
func (r dbMessagesRepository) GetHistory(roomID int, limit int) (messages []*model.Message) {
	rows, err := r.db.Query(`
	SELECT "Messages"."MessageID", "Messages"."Content", "Messages"."UserID", "Messages"."SentAt" ,
	(SELECT array(SELECT "RoomID" FROM "Messages_Rooms" WHERE "Messages_Rooms"."MessageID" = "Messages"."MessageID")) AS "RoomIDs",
	(SELECT "Name" FROM "Users" WHERE "Messages"."UserID" = "Users"."UserID") as "UserName"
	FROM "Messages", "Messages_Rooms"
	WHERE "Messages_Rooms"."RoomID" = $1 AND "Messages"."MessageID" = "Messages_Rooms"."MessageID"
	LIMIT $2
	`, roomID, limit)
	if err != nil {
		log.Println("Query error: ", err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		message := &model.Message{}
		// var ret pq.Int64Array
		err = rows.Scan(
			&message.ID,
			&message.Content,
			&message.UserID,
			&message.SentAt,
			&message.RoomIDs,
			&message.Username)
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

func (r dbMessagesRepository) Create(message *model.Message) (id int, err error) {
	//TODO
	log.Println("to do Save message")
	// tx, err := r.db.Begin()
	// if err != nil {
	// 	log.Println(err)
	// 	return id, err
	// }
	if err := r.db.QueryRow(r.createQuery, message.Content, message.UserID).Scan(&id); err != nil {
		// tx.Rollback()
		log.Println(err)
		return id, err
	}
	fmt.Println(id)
	for _, roomID := range message.RoomIDs {
		if _, err := r.db.Exec(`
		INSERT INTO "Messages_Rooms"
			("MessageID", "RoomID")
		VALUES ($1, $2)
		`, id, roomID); err != nil {
			// tx.Rollback()
			log.Println(err, id, roomID)
			return id, err
		}
	}

	// if err = tx.Commit(); err != nil {
	// 	log.Println(err, id)
	// 	return id, err
	// }
	return id, nil
}

func (r dbMessagesRepository) CheckIfExists(messageID int) (exists bool) {
	_ = r.db.QueryRow(r.checkIfExistsQuery, messageID).Scan(&exists)
	return exists
}
