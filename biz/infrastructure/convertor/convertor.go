package convertor

import (
	notificationmapper "github.com/CloudStriver/cloudmind-system/biz/infrastructure/mapper/notification"
	slidermapper "github.com/CloudStriver/cloudmind-system/biz/infrastructure/mapper/slider"
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

func SliderToSliderMapper(in *gensystem.Slider) *slidermapper.Slider {
	oid, _ := primitive.ObjectIDFromHex(in.SliderId)
	return &slidermapper.Slider{
		ID:       oid,
		ImageUrl: in.ImageUrl,
		LinkUrl:  in.LinkUrl,
		Type:     in.Type,
		IsPublic: in.IsPublic,
	}
}

func SliderMapperToSlider(in *slidermapper.Slider) *gensystem.Slider {
	return &gensystem.Slider{
		SliderId:   in.ID.Hex(),
		ImageUrl:   in.ImageUrl,
		LinkUrl:    in.LinkUrl,
		Type:       in.Type,
		IsPublic:   in.IsPublic,
		CreateTime: in.CreateAt.UnixMilli(),
		UpdateTime: in.UpdateAt.UnixMilli(),
	}
}

func SliderFilterOptionsToMapper(in *gensystem.SliderFilterOptions) *slidermapper.FilterOptions {
	if in == nil {
		return &slidermapper.FilterOptions{}
	}
	return &slidermapper.FilterOptions{
		OnlyType:     in.OnlyType,
		OnlyIsPublic: in.OnlyIsPublic,
	}
}
