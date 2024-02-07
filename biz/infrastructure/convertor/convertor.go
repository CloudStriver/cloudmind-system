package convertor

import (
	notificationmapper "github.com/CloudStriver/cloudmind-system/biz/infrastructure/mapper/notification"
	slidermapper "github.com/CloudStriver/cloudmind-system/biz/infrastructure/mapper/slider"
	gensystem "github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/system"
)

func NotificationInfoToNotificationMapper(in *gensystem.NotificationInfo) *notificationmapper.Notification {
	return &notificationmapper.Notification{
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
		CreateTime:      in.CreateAt.UnixMilli(),
		IsRead:          in.IsRead,
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
