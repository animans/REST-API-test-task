package infrastructure

import (
	"log"
	"testing"
	"time"

	"github.com/animans/REST-API-test-task/domain"
	"github.com/google/uuid"
)

func TestOpen(t *testing.T) {
	s := NewServiceRepoPG()
	err := s.Open()
	if err != nil {
		log.Fatal(err)
	}
	s.Close()
}

func TestSave(t *testing.T) {
	s := NewServiceRepoPG()
	s.Open()

	uuid, _ := uuid.NewRandom()
	ser := domain.NewService("Yandex Plus", 400, uuid, time.Now())
	s.Save(ser)
	s.Close()
}
