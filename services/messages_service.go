package services

import (
	"time"

	"github.com/silvergama/efficientAPI/domain"
	"github.com/silvergama/efficientAPI/utils/errorutils"
)

var (
	MessagesService messageServiceInterface = &messagesService{}
)

type messagesService struct{}

type messageServiceInterface interface {
	GetMessage(int64) (*domain.Message, errorutils.MessageErr)
	CreateMessage(*domain.Message) (*domain.Message, errorutils.MessageErr)
	UpdateMessage(*domain.Message) (*domain.Message, errorutils.MessageErr)
	DeleteMessage(int64) errorutils.MessageErr
	GetAllMessages() ([]domain.Message, errorutils.MessageErr)
}

func (m *messagesService) GetMessage(msgId int64) (*domain.Message, errorutils.MessageErr) {
	message, err := domain.MessageRepo.Get(msgId)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (m *messagesService) GetAllMessages() ([]domain.Message, errorutils.MessageErr) {
	messages, err := domain.MessageRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (m *messagesService) CreateMessage(message *domain.Message) (*domain.Message, errorutils.MessageErr) {
	if err := message.Validate(); err != nil {
		return nil, err
	}
	message.CreatedAt = time.Now()
	message, err := domain.MessageRepo.Create(message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (m *messagesService) UpdateMessage(message *domain.Message) (*domain.Message, errorutils.MessageErr) {
	if err := message.Validate(); err != nil {
		return nil, err
	}
	current, err := domain.MessageRepo.Get(message.ID)
	if err != nil {
		return nil, err
	}
	current.Title = message.Title
	current.Body = message.Body

	updateMsg, err := domain.MessageRepo.Update(current)
	if err != nil {
		return nil, err
	}
	return updateMsg, nil
}

func (m *messagesService) DeleteMessage(msgId int64) errorutils.MessageErr {
	msg, err := domain.MessageRepo.Get(msgId)
	if err != nil {
		return err
	}
	deleteErr := domain.MessageRepo.Delete(msg.ID)
	if deleteErr != nil {
		return deleteErr
	}
	return nil
}
