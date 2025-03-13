package services

import (
	"encoding/json"
	"log"
	"normaladmin/backend/internal/models"
	"normaladmin/backend/pkg/mq"
)

type NoticeService struct {
	mq *mq.RabbitMQ
}

func NewNoticeService(mq *mq.RabbitMQ) *NoticeService {
	return &NoticeService{mq: mq}
}

func (s *NoticeService) SendNotice(notice *models.SystemNotice) error {
	message, err := json.Marshal(notice)
	if err != nil {
		return err
	}
	return s.mq.Publish("notices", message)
}

func (s *NoticeService) StartNoticeConsumer() {
	messages, err := s.mq.Consume("notices")
	if err != nil {
		log.Fatalf("Failed to start notice consumer: %v", err)
	}

	go func() {
		for msg := range messages {
			var notice models.SystemNotice
			if err := json.Unmarshal(msg.Body, &notice); err != nil {
				log.Printf("Failed to parse notice: %v", err)
				continue
			}
			// 处理通知（如保存到数据库、发送邮件等）
			log.Printf("Received notice: %+v", notice)
		}
	}()
}
