package grpcapi

import (
	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/domain/port"
)

type JournalServer struct {
	journalv1.UnimplementedJournalServiceServer
	Usecase port.BulletService
}
