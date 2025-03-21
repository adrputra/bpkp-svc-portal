package router

import (
	"bpkp-svc-portal/app/client"
	"bpkp-svc-portal/app/config"
	"bpkp-svc-portal/app/controller"
	"bpkp-svc-portal/app/service"

	"github.com/aws/aws-sdk-go/service/s3"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"

	"gorm.io/gorm"
)

type ServiceFactory struct {
	user        service.InterfaceUserService
	role        service.InterfaceRoleService
	param       service.InterfaceParamService
	attendance  service.InterfaceAttendanceService
	institution service.InterfaceInstitutionService
}

type ControllerFactory struct {
	user        controller.InterfaceUserController
	role        controller.InterfaceRoleController
	param       controller.InterfaceParamController
	attendance  controller.InterfaceAttendanceController
	institution controller.InterfaceInstitutionController
}

type ClientFactory struct {
	user        client.InterfaceUserClient
	storage     client.InterfaceStorageClient
	role        client.InterfaceRoleClient
	param       client.InterfaceParamClient
	attendance  client.InterfaceAttendanceClient
	institution client.InterfaceInstitutionClient
}

type Factory struct {
	Service    ServiceFactory
	Controller ControllerFactory
	Client     ClientFactory
}

var factory *Factory

func InitFactory(cfg *config.Config, db *gorm.DB, s3 *s3.S3, redis *redis.Client, mq *amqp.Channel) {
	client := ClientFactory{
		user:        client.NewUserClient(db, cfg),
		storage:     client.NewStorageClient(s3, db),
		role:        client.NewRoleClient(db),
		param:       client.NewParamClient(db, redis),
		attendance:  client.NewAttendanceClient(db),
		institution: client.NewInstitutionClient(db),
	}
	controller := ControllerFactory{
		user:        controller.NewUserController(client.user, client.role, client.param, client.storage),
		role:        controller.NewRoleController(client.role),
		param:       controller.NewParamController(redis, client.param),
		attendance:  controller.NewAttendanceController(client.attendance, client.param),
		institution: controller.NewInstitutionController(client.institution),
	}
	service := ServiceFactory{
		user:        service.NewUserService(controller.user),
		role:        service.NewRoleService(controller.role),
		param:       service.NewParamService(controller.param),
		attendance:  service.NewAttendanceService(controller.attendance),
		institution: service.NewInstitutionService(controller.institution),
	}
	factory = &Factory{
		Service:    service,
		Controller: controller,
		Client:     client,
	}
}
