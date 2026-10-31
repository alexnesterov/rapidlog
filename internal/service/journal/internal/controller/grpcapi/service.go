package grpcapi

import (
	"log/slog"

	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/port"
	"google.golang.org/grpc"
)

type JournalService struct {
	journalv1.UnimplementedJournalServiceServer
	usecase port.Usecase
	logger  *slog.Logger
}

func NewJournalService(usecase port.Usecase, logger *slog.Logger) *JournalService {
	return &JournalService{
		usecase: usecase,
		logger:  logger,
	}
}

func RegisterJournalService(srv *grpc.Server, uc port.Usecase, logger *slog.Logger) {
	journalService := NewJournalService(uc, logger)
	journalv1.RegisterJournalServiceServer(srv, journalService)
}
