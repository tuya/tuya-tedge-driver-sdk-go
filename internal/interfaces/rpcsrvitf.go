package interfaces

import (
	"context"

	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/dpmodel"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/cache"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/clients"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/config"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/proto"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/worker"
	"github.com/tuya/tuya-tedge-driver-sdk-go/thingmodel"
)

type RPCServerItf interface {
	Serve() error
}

type DriverCommonItf interface {
	GetContext() (context.Context, context.CancelFunc)
	GetRPCConfig() config.RPCConfig
	GetLogger() commons.TedgeLogger
	ChangeLogLevel(level string) error
	SetCloudStatus(status bool)
	GetDevCache() *cache.DeviceCache

	GetAppHandler() clients.AppCallBack
	UpdateAppAddress(ctx context.Context, req *proto.AppBaseAddress) error
}

type DpDriverService interface {
	DriverCommonItf

	GetDriver() dpmodel.DPModelDriver
	GetPdCache() *cache.DPModelProductCache
}

type TyDriverService interface {
	DriverCommonItf

	GetDriver() thingmodel.ThingModelDriver
	GetPdCache() *cache.ThingModelProductCache
	GetTMWorkerPool() *worker.TMWorkerPool
}
