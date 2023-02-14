package internal

import (
	. "config/primitive"
	"fmt"
	"strings"
	"sync"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

const (
	defaultLogDir   = "/tmp/nacos/log"
	defaultCacheDir = "/tmp/nacos/cache"
	defaultLogLevel = "debug"
)

var registeredDataIDAndGroup = map[string]struct{}{}

type nacosConfig struct {
	mutex sync.RWMutex

	dataID  string
	group   string
	content string
}

// isRegistered 检查dataID和group是否已经注册过
func isRegistered(dataID, group string) bool {
	_, exist := registeredDataIDAndGroup[fmt.Sprintf("%s:%s", dataID, group)]
	return exist
}

// markRegistered 标记dataID和group已经注册过
func markRegistered(dataID, group string) {
	registeredDataIDAndGroup[fmt.Sprintf("%s:%s", dataID, group)] = struct{}{}
}

// clearRegistered 清理已经注册过的dataID和group
func clearRegistered() {
	registeredDataIDAndGroup = map[string]struct{}{}
}

// NewNacosClient 创建Nacos客户端
//  e.g
//  c := NewNacosClient("localhost:8080", namespace, WithAccessKey("accessKey"), WithSecretKey("secretKey"))
func NewNacosClient(addr, namespace string, opts ...ClientOption) (INacosClient, error) {
	cc := constant.NewClientConfig(
		constant.WithEndpoint(addr),
		constant.WithNamespaceId(namespace),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir(defaultLogDir),
		constant.WithCacheDir(defaultCacheDir), // 断网情况下，GetConfig()会读取缓存数据
		constant.WithLogLevel(defaultLogLevel),
		constant.WithLogStdout(true),
	)

	for _, opt := range opts {
		opt(cc)
	}

	client, err := clients.NewConfigClient(vo.NacosClientParam{ClientConfig: cc})
	if err != nil {
		return nil, err
	}

	// 用于nacos sdk不支持ping操作，这个只能通过PublishConfig发起请求，如果返回的err中
	// 包含"server list is empty"，表示无法登录nacos server
	_, err = client.PublishConfig(vo.ConfigParam{
		DataId:  "echo",
		Group:   "echo",
		Content: "ping",
	})

	if err != nil && strings.Contains(err.Error(), "server list is empty") {
		return nil, ErrConnectFailed
	}

	return client, nil
}
