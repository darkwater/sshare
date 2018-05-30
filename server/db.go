package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/darkwater/sshare/common"
	"github.com/golang/protobuf/proto"
	_ "github.com/mattn/go-sqlite3"
)

// Invite represents a single invite, handed out by a user or the system
type Invite struct {
	ID     int
	Sender int
}

// User represents a single user in the system
type User struct {
	ID               int
	URL              string
	InvitedBy        int
	InvitesAvailable int
}

var db *sql.DB

func dbInit() error {
	var err error
	db, err = sql.Open("sqlite3", "./sshare.db")
	if err != nil {
		return err
	}

	// TODO: Only do this when the database was newly created to help detect broken/outdated databases
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS invites (
		code        STRING  PRIMARY KEY,
		sender      INTEGER
	);
	CREATE TABLE IF NOT EXISTS users (
		userid            INTEGER PRIMARY KEY ASC,
		url               STRING  UNIQUE,
		invited_by        INTEGER,
		invites_available INTEGER
	);
	CREATE TABLE IF NOT EXISTS public_keys (
		hash        STRING  PRIMARY KEY,
		userid      INTEGER,
		data        BLOB
	);
	`)
	if err != nil {
		return err
	}

	return nil
}

func dbGetInvite(code string) (invite *Invite, err error) {
	stmt, err := db.Prepare(`SELECT rowid, sender FROM invites WHERE code = ?`)
	if err != nil {
		return
	}
	defer stmt.Close()

	row := stmt.QueryRow(code)
	invite = &Invite{}
	err = row.Scan(&invite.ID, &invite.Sender)
	if err != nil {
		return
	}

	return
}

func dbUseInvite(invite *Invite, user *User) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}

	// Delete invite
	stmt, err := tx.Prepare(`DELETE FROM invites WHERE rowid = ?`)
	if err != nil {
		tx.Rollback()
		panic("invalid query")
	}
	defer stmt.Close()

	result, err := stmt.Exec(invite.ID)
	// TODO: handle sql.ErrNoRows?
	if err != nil {
		tx.Rollback()
		return
	}
	affected, _ := result.RowsAffected()
	if affected != 1 {
		// Application error? dbGetInvite should've been called first. Race condition possible, however.
		// Can this even happen without ErrNoRows?
		tx.Rollback()
		return errors.New("invite not found")
	}

	// Create user
	stmt, err = tx.Prepare(`INSERT INTO users (url, invited_by) VALUES (?, ?)`)
	if err != nil {
		tx.Rollback()
		panic("invalid query")
	}
	defer stmt.Close()

	result, err = stmt.Exec(user.URL, user.InvitedBy)
	if err != nil {
		tx.Rollback()
		panic("failed to insert user")
	}
	uid, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		panic("failed to get new user id")
	}
	user.ID = int(uid)

	err = tx.Commit()
	if err != nil {
		panic("transaction failed")
	}

	return
}

func dbAddPublicKey(user *User, key *common.PublicKey) (err error) {
	stmt, err := db.Prepare(`INSERT INTO public_keys (hash, userid, data) VALUES (?, ?, ?)`)
	if err != nil {
		panic("invalid query")
	}
	defer stmt.Close()

	hash := common.HashPublicKey(key.FromProtobuf())
	data, err := proto.Marshal(key)
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec(fmt.Sprintf("%x", hash), user.ID, data)
	if err != nil {
		panic(err)
	}

	return
}
