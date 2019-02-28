package repository

import (
	"database/sql"
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
	_, err = r.db.Exec(r.addOrUpdateUserChannelQuery, userID, channelID, "now()")
	return err
}
