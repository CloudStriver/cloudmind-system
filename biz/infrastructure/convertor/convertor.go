package convertor

import (
	notificationmapper "github.com/CloudStriver/cloudmind-system/biz/infrastructure/mapper/notification"
	gensystem "github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/system"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NotificationToNotificationMapper(in *gensystem.Notification) *notificationmapper.Notification {
	oid, _ := primitive.ObjectIDFromHex(in.NotificationId)
	return &notificationmapper.Notification{
		ID:              oid,
		TargetUserId:    in.TargetUserId,
		SourceUserId:    in.SourceUserId,
		SourceContentId: in.SourceContentId,
		Type:            in.Type,
		TargetType:      in.TargetType,
		Text:            in.Text,
		IsRead:          in.IsRead,
	}
}

func NotificationMapperToNotification(in *notificationmapper.Notification) *gensystem.Notification {
	return &gensystem.Notification{
		NotificationId:  in.ID.Hex(),
		TargetUserId:    in.TargetUserId,
		SourceUserId:    in.SourceUserId,
		SourceContentId: in.SourceContentId,
		TargetType:      in.TargetType,
		Type:            in.Type,
		Text:            in.Text,
		CreateAt:        in.CreateAt.UnixMilli(),
		IsRead:          in.IsRead,
	}
}

func NotificationFilterOptionsToMapper(in *gensystem.NotificationFilterOptions) *notificationmapper.FilterOptions {
	if in == nil {
		return &notificationmapper.FilterOptions{}
	}
	return &notificationmapper.FilterOptions{
		OnlyUserId:     in.OnlyUserId,
		OnlyType:       in.OnlyType,
		OnlyTargetType: in.OnlyTargetType,
		OnlyFirstId:    in.OnlyFirstId,
		OnlyLastId:     in.OnlyLastId,
		OnlyIsRead:     in.OnlyIsRead,
	}
}
