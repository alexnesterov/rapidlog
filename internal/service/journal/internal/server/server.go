package server

import (
	"context"

	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"google.golang.org/grpc"
)

type ListBulletsHandler interface {
	Handle(ctx context.Context, req *journalv1.ListBulletsRequest) (*journalv1.ListBulletsResponse, error)
}

type CreateBulletHandler interface {
	Handle(ctx context.Context, req *journalv1.CreateBulletRequest) (*journalv1.CreateBulletResponse, error)
}

type CompleteBulletHandler interface {
	Handle(ctx context.Context, req *journalv1.CompleteBulletRequest) (*journalv1.CompleteBulletResponse, error)
}

type CancelBulletHandler interface {
	Handle(ctx context.Context, req *journalv1.CancelBulletRequest) (*journalv1.CancelBulletResponse, error)
}

type MigrateBulletHandler interface {
	Handle(ctx context.Context, req *journalv1.MigrateBulletRequest) (*journalv1.MigrateBulletResponse, error)
}

type Handlers struct {
	CreateBullet   CreateBulletHandler
	ListBullets    ListBulletsHandler
	CompleteBullet CompleteBulletHandler
	MigrateBullet  MigrateBulletHandler
	CancelBullet   CancelBulletHandler
}

// JournalService собирает вертикальные слайсы features/* в один gRPC-сервер:
// каждый RPC делегирует в хендлер своего слайса.
type JournalService struct {
	journalv1.UnimplementedJournalServiceServer
	h Handlers
}

func NewJournalService(h Handlers) *JournalService {
	return &JournalService{h: h}
}

func (s *JournalService) ListBullets(ctx context.Context, req *journalv1.ListBulletsRequest) (*journalv1.ListBulletsResponse, error) {
	return s.h.ListBullets.Handle(ctx, req)
}

func (s *JournalService) CreateBullet(ctx context.Context, req *journalv1.CreateBulletRequest) (*journalv1.CreateBulletResponse, error) {
	return s.h.CreateBullet.Handle(ctx, req)
}

func (s *JournalService) CompleteBullet(ctx context.Context, req *journalv1.CompleteBulletRequest) (*journalv1.CompleteBulletResponse, error) {
	return s.h.CompleteBullet.Handle(ctx, req)
}

func (s *JournalService) CancelBullet(ctx context.Context, req *journalv1.CancelBulletRequest) (*journalv1.CancelBulletResponse, error) {
	return s.h.CancelBullet.Handle(ctx, req)
}

func (s *JournalService) MigrateBullet(ctx context.Context, req *journalv1.MigrateBulletRequest) (*journalv1.MigrateBulletResponse, error) {
	return s.h.MigrateBullet.Handle(ctx, req)
}

func Register(srv *grpc.Server, h Handlers) {
	journalService := NewJournalService(h)
	journalv1.RegisterJournalServiceServer(srv, journalService)
}
