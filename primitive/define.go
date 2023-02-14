package primitive

import (
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
)

var (
	WithNamespaceId  = constant.WithNamespaceId
	WithAccessKey    = constant.WithAccessKey
	WithSecretKey    = constant.WithSecretKey
	WithEndpoint     = constant.WithEndpoint
	WithCustomLogger = constant.WithCustomLogger
	WithTimeoutMs    = constant.WithTimeoutMs
	WithUserName     = constant.WithUsername
	WithPassword     = constant.WithPassword
	WithLogLevel     = constant.WithLogLevel
)

type (
	ClientOption = constant.ClientOption
	INacosClient = config_client.IConfigClient
)

type Flag uint8

const (
	Unknown   Flag = iota // 未知配置
	OnlyFile              // 文件模式
	OnlyNacos             // nacos模式
	Mixed                 // 混合模式
)

type IMixedConfig interface {
	UpdateAfterRegister()                                       // 注册成功过，调用该函数进行更新操作
	OnNacosChanged(namespace, group, dataId, data string) error // nacos有变更时触发该函数
}
