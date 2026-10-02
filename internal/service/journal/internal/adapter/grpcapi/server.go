package grpcapi

import (
	"fmt"

	journalv1 "github.com/alexnesterov/rapidlog-api/api/gen/journal/v1"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/domain/entity"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/domain/port"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type JournalServer struct {
	journalv1.UnimplementedJournalServiceServer
	Usecase port.BulletService
}

func ToProtoBulletType(t entity.BulletType) (journalv1.BulletType, error) {
	switch t {
	case entity.BulletTask:
		return journalv1.BulletType_BULLET_TYPE_TASK, nil
	case entity.BulletEvent:
		return journalv1.BulletType_BULLET_TYPE_EVENT, nil
	case entity.BulletNote:
		return journalv1.BulletType_BULLET_TYPE_NOTE, nil
	default:
		return journalv1.BulletType_BULLET_TYPE_UNSPECIFIED, fmt.Errorf("unknown bullet type: %v", t)
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
	case entity.SignifierMigrated:
		return journalv1.Signifier_SIGNIFIER_MIGRATED, nil
	case entity.SignifierScheduled:
		return journalv1.Signifier_SIGNIFIER_SCHEDULED, nil
	default:
		return journalv1.Signifier_SIGNIFIER_UNSPECIFIED, fmt.Errorf("unknown signifier: %v", s)
	}
}

func ToProtoBullet(bullet *entity.Bullet) (*journalv1.Bullet, error) {
	typeProto, err := ToProtoBulletType(bullet.Type)
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
		Type:      typeProto,
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
