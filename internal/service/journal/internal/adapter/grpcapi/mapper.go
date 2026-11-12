package grpcapi

import (
	"fmt"

	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToProtoKind(k entity.Kind) (journalv1.Kind, error) {
	switch k {
	case entity.KindTask:
		return journalv1.Kind_KIND_TASK, nil
	case entity.KindEvent:
		return journalv1.Kind_KIND_EVENT, nil
	case entity.KindNote:
		return journalv1.Kind_KIND_NOTE, nil
	default:
		return journalv1.Kind_KIND_UNSPECIFIED, fmt.Errorf("unknown bullet kind: %v", k)
	}
}

func ToProtoSignifier(s entity.Signifier) (journalv1.Signifier, error) {
	switch s {
	case entity.SignifierOpen:
		return journalv1.Signifier_SIGNIFIER_OPEN, nil
	case entity.SignifierCompleted:
		return journalv1.Signifier_SIGNIFIER_COMPLETED, nil
	case entity.SignifierCancelled:
		return journalv1.Signifier_SIGNIFIER_CANCELLED, nil
	default:
		return journalv1.Signifier_SIGNIFIER_UNSPECIFIED, fmt.Errorf("unknown signifier: %v", s)
	}
}

func ToProtoBullet(bullet *entity.Bullet) (*journalv1.Bullet, error) {
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

func ToProtoBullets(bullets []*entity.Bullet) ([]*journalv1.Bullet, error) {
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
