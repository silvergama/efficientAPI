package services

import (
	"database/sql"
	"time"

	"github.com/silvergama/efficient-api/domain"
	"github.com/silvergama/efficient-api/utils/error_utils"
)

var (
	tm                   = time.Now()
	getMessageDomain     func(messageId int64) (*domain.Message, error_utils.MessageErr)
	createMessageDomain  func(msg *domain.Message) (*domain.Message, error_utils.MessageErr)
	updateMessageDomain  func(msg *domain.Message) (*domain.Message, error_utils.MessageErr)
	deleteMessageDomain  func(messageId int64) error_utils.MessageErr
	getAllMessagesDomain func() ([]domain.Message, error_utils.MessageErr)
)

type getDBMock struct{}

func (m *getDBMock) Get(messageID int64) (*domain.Message, error_utils.MessageErr) {
	return getMessageDomain(messageID)
}

func (m *getDBMock) Create(msg *domain.Message) (*domain.Message, error_utils.MessageErr) {
	return createMessageDomain(msg)
}

func (m *getDBMock) Update(msg *domain.Message) (*domain.Message, error_utils.MessageErr) {
	return updateMessageDomain(msg)
}

func (m *getDBMock) Delete(messageID int64) error_utils.MessageErr {
	return deleteMessageDomain(messageID)
}

func (m *getDBMock) GetAll() ([]domain.Message, error_utils.MessageErr) {
	return getAllMessagesDomain()
}

func (m *getDBMock) Initialize(string, string, string, string, string, string) *sql.DB {
	return nil
}
