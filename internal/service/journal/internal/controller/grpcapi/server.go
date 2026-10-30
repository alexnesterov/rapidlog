package grpcapi

import (
	"log/slog"

	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/port"
)

type JournalServer struct {
	journalv1.UnimplementedJournalServiceServer
	Usecase port.Usecase
	Logger  *slog.Logger
}
