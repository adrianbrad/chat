package repository

import (
	"database/sql"
	"time"
)

type UsersChannelsRepository interface {
	AddOrUpdateUserToChannel(int, int) error
}

type usersChannelsRepository struct {
	db                          *sql.DB
	addOrUpdateUserChannelQuery string
}

func NewDbUsersChannelsRepository(database *sql.DB) UsersChannelsRepository {
	return &usersChannelsRepository{
		db:                          database,
		addOrUpdateUserChannelQuery: AddOrUpdateManyToMany("User", "Channel", "Joined"),
	}
}

func (r usersChannelsRepository) AddOrUpdateUserToChannel(userID int, channelID int) (err error) {
	t := time.Now()
	ts := t.Format("2006-01-02 15:04:05")
	_, err = r.db.Exec(r.addOrUpdateUserChannelQuery, userID, channelID, ts)
	return err
}
