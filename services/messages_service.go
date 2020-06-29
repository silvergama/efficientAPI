package services

import (
	"time"

	"github.com/silvergama/efficient-api/domain"
	"github.com/silvergama/efficient-api/utils/errorutils"
)

type messagesService struct{}

type messageServiceInterface interface {
	GetMessage(int64) (*domain.Message, errorutils.MessageErr)
	CreateMessage(*domain.Message) (*domain.Message, errorutils.MessageErr)
	UpdateMessage(*domain.Message) (*domain.Message, errorutils.MessageErr)
	DeteleMessage(int64) errorutils.MessageErr
	GetAllMessages() ([]domain.Message, errorutils.MessageErr)
}

func (m *messagesService) GetMessage(msgID int64) (*domain.Message, errorutils.MessageErr) {
	message, err := domain.MessageRepo.Get(msgID)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (m *messagesService) GetAllMessages() ([]domain.Message, errorutils.MessageErr) {
	message, err := domain.MessageRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return message, nil
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

	return message, err
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

func (m *messagesService) DeleteMessage(msgID int64) errorutils.MessageErr {
	msg, err := domain.MessageRepo.Get(msgID)
	if err != nil {
		return err
	}

	deleteErr := domain.MessageRepo.Delete(msg.ID)
	if err != nil {
		return deleteErr
	}

	return nil
}
