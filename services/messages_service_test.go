package services

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/silvergama/efficientAPI/domain"
	"github.com/silvergama/efficientAPI/utils/errorutils"
	"github.com/stretchr/testify/assert"
)

var (
	tm                   = time.Now()
	getMessageDomain     func(messageId int64) (*domain.Message, errorutils.MessageErr)
	createMessageDomain  func(msg *domain.Message) (*domain.Message, errorutils.MessageErr)
	updateMessageDomain  func(msg *domain.Message) (*domain.Message, errorutils.MessageErr)
	deleteMessageDomain  func(messageId int64) errorutils.MessageErr
	getAllMessagesDomain func() ([]domain.Message, errorutils.MessageErr)
)

type getDBMock struct{}

func (m *getDBMock) Get(messageID int64) (*domain.Message, errorutils.MessageErr) {
	return getMessageDomain(messageID)
}

func (m *getDBMock) Create(msg *domain.Message) (*domain.Message, errorutils.MessageErr) {
	return createMessageDomain(msg)
}

func (m *getDBMock) Update(msg *domain.Message) (*domain.Message, errorutils.MessageErr) {
	return updateMessageDomain(msg)
}

func (m *getDBMock) Delete(messageID int64) errorutils.MessageErr {
	return deleteMessageDomain(messageID)
}

func (m *getDBMock) GetAll() ([]domain.Message, errorutils.MessageErr) {
	return getAllMessagesDomain()
}

func (m *getDBMock) Initialize(string, string, string, string, string, string) *sql.DB {
	return nil
}

///////////////////////////////////////////////////////////
// Start of "GetMessge" tests cases
///////////////////////////////////////////////////////////
func TestMessagesService_GetMessage_Success(t *testing.T) {
	domain.MessageRepo = &getDBMock{} // this is where we swapped the functionality
	getMessageDomain = func(messageId int64) (*domain.Message, errorutils.MessageErr) {
		return &domain.Message{
			ID:        1,
			Title:     "the title",
			Body:      "the body",
			CreatedAt: tm,
		}, nil
	}
	msg, err := MessagesService.GetMessage(1)
	fmt.Println("this is the message: ", msg)
	assert.NotNil(t, msg)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, msg.ID)
	assert.EqualValues(t, "the title", msg.Title)
	assert.EqualValues(t, "the body", msg.Body)
	assert.EqualValues(t, tm, msg.CreatedAt)
}
