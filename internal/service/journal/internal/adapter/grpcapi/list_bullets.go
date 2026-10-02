package grpcapi

import (
	"context"

	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *JournalServer) ListBullets(ctx context.Context, req *journalv1.ListBulletsRequest) (*journalv1.ListBulletsResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, codes.InvalidArgument.String())
	}

	bullets, err := s.Usecase.ListBullets(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}

	bulletsProto, err := ToProtoBullets(bullets)
	if err != nil {
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}

	return &journalv1.ListBulletsResponse{
		Bullets: bulletsProto,
	}, nil
}
