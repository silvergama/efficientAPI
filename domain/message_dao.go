package domain

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/silvergama/efficient-api/utils/error_formats"
	"github.com/silvergama/efficient-api/utils/error_utils"
)

var (
	MessageRepo messageRepoInterface = &messageRepo{}
)

const (
	queryGetMessage    = "SELECT id, title, body, created_at FROM messages WHERE id=?;"
	queryInsertMessage = "INSERT INTO messages(title, body, created_at) VALUES(?, ?, ?);"
	queryUpdateMessge  = "UPDATE messages SET title=?, body=?, created=? WHERE id=?;"
	queryDeleteMessage = "DELETE FROM messages WHERE id=?;"
	queryGetAllMessage = "SELECT id, tite, body, created_at FROM messages;"
)

type messageRepoInterface interface {
	Get(int64) (*Message, error_utils.MessageErr)
	Create(*Message) (*Message, error_utils.MessageErr)
	Update(*Message) (*Message, error_utils.MessageErr)
	Delete(int64) error_utils.MessageErr
	GetAll() ([]Message, error_utils.MessageErr)
	Initialize(string, string, string, string, string, string) *sql.DB
}

type messageRepo struct {
	db *sql.DB
}

func (mr *messageRepo) Initialize(DbDriver, DbUser, DbPassword, DbPort, DbHost, DbName string) *sql.DB {
	var err error
	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)

	mr.db, err = sql.Open(DbDriver, DBURL)
	if err != nil {
		log.Fatal("This is the error connecting to the database:", err)
	}
	fmt.Printf("We are connected to the %s database", DbDriver)

	return mr.db
}

func NewMessageRepository(db *sql.DB) messageRepoInterface {
	return &messageRepo{
		db: db,
	}
}

type Address struct {
	State        string
	City         string
	Street       string
	Neighborhood string
	ZipCode      string
	StreetNumber string
}

func (mr *messageRepo) Get(messageId int64) (*Message, error_utils.MessageErr) {
	stmt, err := mr.db.Prepare(queryGetMessage)
	if err != nil {
		return nil, error_utils.NewInternalServerError(fmt.Sprintf("Error when trying to prepare message: %s", err.Error()))
	}
	defer stmt.Close()

	var msg Message
	getError := stmt.QueryRow(messageId).Scan(
		&msg.ID,
		&msg.Title,
		&msg.Body,
		&msg.CreatedAt,
	)
	if getError != nil {
		fmt.Println("This is the error man:", getError)
		return nil, error_formats.ParseError(getError)
	}
	return &msg, nil
}

func (mr *messageRepo) GetAll() ([]Message, error_utils.MessageErr) {
	stmt, err := mr.db.Prepare(queryGetAllMessage)
	if err != nil {
		return nil, error_utils.NewInternalServerError(fmt.Sprintf("Error when trying to prepare all messages %s", err.Error()))
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, error_formats.ParseError(err)
	}
	defer rows.Close()

	results := make([]Message, 0)

	for rows.Next() {
		var msg Message
		getError := rows.Scan(
			&msg.ID,
			&msg.Title,
			&msg.Body,
			&msg.CreatedAt,
		)
		if getError != nil {
			return nil, error_utils.NewInternalServerError(fmt.Sprintf("Error when trying to get message %s", err.Error()))
		}
		results = append(results, msg)
	}
	if len(results) == 0 {
		return nil, error_utils.NewNotFoundError("no records found")
	}
	return results, nil
}

func (mr *messageRepo) Create(msg *Message) (*Message, error_utils.MessageErr) {
	fmt.Println("WE REACHED THE DOMAIN")
	stmt, err := mr.db.Prepare(queryInsertMessage)
	if err != nil {
		return nil, error_utils.NewInternalServerError(fmt.Sprintf("error when trying to prepare message to save %s", err.Error()))
	}
	fmt.Println("WE DIDN'T REACH HERE")
	defer stmt.Close()

	insertResult, createErr := stmt.Exec(
		msg.Title, msg.Body, msg.CreatedAt,
	)
	if createErr != nil {
		return nil, error_formats.ParseError(createErr)
	}
	msgId, err := insertResult.LastInsertId()
	if err != nil {
		return nil, error_utils.NewInternalServerError(fmt.Sprintf("error trying to save message %s", err.Error()))
	}
	msg.ID = msgId

	return msg, nil
}

func (mr *messageRepo) Update(msg *Message) (*Message, error_utils.MessageErr) {
	stmt, err := mr.db.Prepare(queryUpdateMessge)
	if err != nil {
		return nil, error_utils.NewInternalServerError(fmt.Sprintf("error when trying to prepare user to save %s", err.Error()))
	}
	defer stmt.Close()

	if _, updErr := stmt.Exec(msg.Title, msg.Body, msg.ID); updErr != nil {
		return nil, error_formats.ParseError(updErr)
	}

	return msg, nil
}

func (mr *messageRepo) Delete(msgId int64) error_utils.MessageErr {
	stmt, err := mr.db.Prepare(queryDeleteMessage)
	if err != nil {
		return error_utils.NewInternalServerError(fmt.Sprintf("error when trying prepare message to delete %s", err.Error()))
	}
	defer stmt.Close()

	if _, err := stmt.Exec(msgId); err != nil {
		return error_formats.ParseError(err)
	}
	return nil
}
