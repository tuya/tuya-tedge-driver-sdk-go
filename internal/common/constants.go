package common

import "time"

const (
	GRPCTimeout          = 5 * time.Second
	GRPCHttpTimeout      = 10 * time.Second
	DefaultConfiguration = "/etc/driver/res/configuration.toml"
)

const (
	HTTP_API_EG_SYNCDATA_RESULT                   = "tuya.industry.base.eg.syncdata.result.post"
	HTTP_API_EDGE_GATEWAY_SYNCDATA_RESULT         = "tuya.industry.base.edge.gateway.syncdata.type.result.post"
	HTTP_API_EDGE_UPLOAD_TOKEN                    = "tuya.industry.base.edge.file.upload.token"
	HTTP_API_EDGE_DEVICE_QUERY_ALL                = "tuya.industry.base.edge.device.query.all"
	HTTP_API_EG_DEVICE_CID_QUERY                  = "tuya.industry.base.eg.device.cid.query"
	HTTP_API_EDGE_IPC_MODIFY_SKILL                = "tuya.industry.base.edge.ipc.modify.skill"
	HTTP_API_EDGE_FILE_UPLOAD_ALG                 = "tuya.industry.base.edge.file.upload.alg"
	HTTP_API_DEVICE_IPC_TY_SUB_STORAGE_SECRET_GET = "tuya.device.ipc.ty.sub.storage.secret.get "
	HTTP_API_EDGE_IPC_QUERY_DETAIL                = "tuya.industry.base.edge.ipc.query.detail"
	HTTP_DEVICE_IPC_STORAGE_SECRET_GET            = "tuya.device.ipc.storage.secret.get"
	HTTP_DEVICE_SKILL_UPDATE                      = "tuya.device.skill.update"
	HTTP_DEVICE_TY_SUB_STOR_CONFIG_GET            = "tuya.device.ty.sub.stor.config.get"
)

const (
	DeviceFromUnknown = iota
	DeviceFromWeb
	DeviceFromDriver
	DeviceActive
)
