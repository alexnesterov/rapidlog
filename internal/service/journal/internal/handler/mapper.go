package handler

import (
	"fmt"

	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToProtoKind(k domain.Kind) (journalv1.Kind, error) {
	switch k {
	case domain.KindTask:
		return journalv1.Kind_KIND_TASK, nil
	case domain.KindEvent:
		return journalv1.Kind_KIND_EVENT, nil
	case domain.KindNote:
		return journalv1.Kind_KIND_NOTE, nil
	default:
		return journalv1.Kind_KIND_UNSPECIFIED, fmt.Errorf("unknown bullet kind: %v", k)
	}
}

func ToProtoSignifier(s domain.Signifier) (journalv1.Signifier, error) {
	switch s {
	case domain.SignifierOpen:
		return journalv1.Signifier_SIGNIFIER_OPEN, nil
	case domain.SignifierCompleted:
		return journalv1.Signifier_SIGNIFIER_COMPLETED, nil
	case domain.SignifierCancelled:
		return journalv1.Signifier_SIGNIFIER_CANCELLED, nil
	default:
		return journalv1.Signifier_SIGNIFIER_UNSPECIFIED, fmt.Errorf("unknown signifier: %v", s)
	}
}

func ToProtoBullet(bullet *domain.Bullet) (*journalv1.Bullet, error) {
	kindProto, err := ToProtoKind(bullet.Kind)
	if err != nil {
		return nil, err
	}
	signifierProto, err := ToProtoSignifier(bullet.Signifier)
	if err != nil {
		return nil, err
	}
	return &journalv1.Bullet{
		Id:        bullet.ID.String(),
		UserId:    bullet.UserID.String(),
		Content:   bullet.Content,
		Kind:      kindProto,
		Signifier: signifierProto,
		CreatedAt: timestamppb.New(bullet.CreatedAt),
		UpdatedAt: timestamppb.New(bullet.UpdatedAt),
	}, nil
}

func ToProtoBullets(bullets []*domain.Bullet) ([]*journalv1.Bullet, error) {
	result := make([]*journalv1.Bullet, 0, len(bullets))
	for _, bullet := range bullets {
		bulletProto, err := ToProtoBullet(bullet)
		if err != nil {
			return nil, err
		}
		result = append(result, bulletProto)
	}
	return result, nil
}
