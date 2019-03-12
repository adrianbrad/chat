package repository

import (
	"database/sql"
	"log"

	"github.com/adrianbrad/chat/model"
)

type ChannelRepository interface {
	GetOne(int) (*model.Channel, error)
}

type dbChannelsRepository struct {
	db                 *sql.DB
	getOneQuery        string
	getAllQuery        string
	createQuery        string
	checkIfExistsQuery string
}

func NewDbChannelsRepository(database *sql.DB) ChannelRepository {
	return &dbChannelsRepository{
		db:                 database,
		getOneQuery:        getOneQuery("Channel", "ChannelID", "Name", "Description"),
		getAllQuery:        getAllQuery("Channel", "ChannelID", "Name", "Description"),
		createQuery:        createOneQuery("Channel", "Name", "Description"),
		checkIfExistsQuery: checkIfExistsQuery("Channel"),
	}
}

func (r dbChannelsRepository) GetOne(id int) (*model.Channel, error) {
	channel := &model.Channel{}
	err := r.db.QueryRow(r.getOneQuery, id).Scan(
		&channel.ID,
		&channel.Name,
		&channel.Description)
	if err != nil {
		log.Println("Error while fetching channel with id", id)
		return channel, err
	}
	return channel, nil
}

func (r dbChannelsRepository) GetAll() (channels []interface{}) {
	rows, err := r.db.Query(r.getAllQuery)
	if err != nil {
		log.Println("Query error: ", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		channel := model.Channel{}
		err = rows.Scan(
			&channel.ID,
			&channel.Name,
			&channel.Description)
		if err != nil {
			log.Println("Mapping error", err)
			return
		}
		channels = append(channels, channel)
	}
	err = rows.Err()
	if err != nil {
		log.Println("Reading rows error:", err)
	}
	return channels
}

func (r dbChannelsRepository) Create(channelI interface{}) (id int, err error) {
	channel := channelI.(model.Channel)
	if err := r.db.QueryRow(
		r.createQuery, channel.Name, channel.Description).
		Scan(&id); err != nil {
		return id, err
	}
	return id, nil
}

func (r dbChannelsRepository) CheckIfExists(channelID int) (exists bool) {
	_ = r.db.QueryRow(r.checkIfExistsQuery, channelID).Scan(&exists)
	return exists
}

func (r dbChannelsRepository) GetAllWhere(cloumn string, value int, limit int) []interface{} {

	return nil
}
