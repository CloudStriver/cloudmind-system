package provider

import (
	slidermapper "github.com/CloudStriver/cloudmind-system/biz/infrastructure/mapper/slider"
	"github.com/google/wire"

	"github.com/CloudStriver/cloudmind-system/biz/application/service"
	"github.com/CloudStriver/cloudmind-system/biz/infrastructure/config"
	notificationmapper "github.com/CloudStriver/cloudmind-system/biz/infrastructure/mapper/notification"
)

var AllProvider = wire.NewSet(
	ApplicationSet,
	InfrastructureSet,
)

var ApplicationSet = wire.NewSet(
	service.SystemSet,
)

var InfrastructureSet = wire.NewSet(
	config.NewConfig,
	MapperSet,
)

var MapperSet = wire.NewSet(
	notificationmapper.NewNotificationModel,
	slidermapper.NewSliderModel,
)
