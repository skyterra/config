package internal

import (
	. "config/primitive"
	"errors"
	"io/ioutil"
	"reflect"
	"sync"

	"gopkg.in/yaml.v2"

	"github.com/nacos-group/nacos-sdk-go/vo"
)

const defaultName = "default"

type configIns struct {
	mutex sync.RWMutex

	namespace string       // nacos客户端访问的namespace
	client    INacosClient // nacos客户端

	mixed map[string]interface{}
	files map[string]interface{}
	nacos map[string]*nacosConfig
}

// DailNacos 注册nacos客户端
func (c *configIns) DailNacos(addr, namespace string, opts ...ClientOption) error {
	if c.client != nil {
		return errors.New("nacos client has been init")
	}

	client, err := NewNacosClient(addr, namespace, opts...)
	if err != nil {
		return err
	}

	c.client = client
	c.namespace = namespace
	return nil
}

// RegisterConfig 注册配置文件
func (c *configIns) RegisterFile(file string, v interface{}) error {
	return c.RegisterFileWithName(defaultName, file, v)
}

// RegisterNacos 注册nacos dataID和group
func (c *configIns) RegisterNacos(dataID, group string) error {
	return c.RegisterNacosWithName(defaultName, dataID, group)
}

// RegisterMixed 注册可更新配置
func (c *configIns) RegisterMixed(file, dataID, group string, v IMixedConfig) error {
	return c.RegisterMixedWithName(defaultName, file, dataID, group, v)
}

// RegisterConfig 注册配置文件
func (c *configIns) RegisterFileWithName(name string, file string, v interface{}) error {
	if err := checkType(v); err != nil {
		return err
	}

	if _, exist := c.files[name]; exist {
		return ErrAlreadyRegister
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	// copy出一个新的空对象，由于存储配置信息
	fileConf, err := copyAndUnmarshal(data, v)
	if err != nil {
		return err
	}

	c.files[name] = fileConf
	return nil
}

// RegisterNacos 注册nacos dataID和group
func (c *configIns) RegisterNacosWithName(name, dataID, group string) error {
	if isRegistered(dataID, group) {
		return ErrDataIDAndGroupAlreadyRegister
	}

	content, err := c.client.GetConfig(vo.ConfigParam{DataId: dataID, Group: group})
	if err != nil {
		return err
	}

	if content == "" {
		return ErrNotExistConfig
	}

	err = c.client.ListenConfig(vo.ConfigParam{
		DataId: dataID,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			if namespace != c.namespace {
				return
			}

			for _, conf := range c.nacos {
				if conf.dataID == dataID && conf.group == group {
					conf.mutex.Lock()
					conf.content = data
					conf.mutex.Unlock()
					break
				}
			}
		},
	})

	if err != nil {
		return err
	}

	markRegistered(dataID, group)
	c.nacos[name] = &nacosConfig{
		dataID:  dataID,
		group:   group,
		content: content,
	}

	return nil
}

// RegisterMixedWithName 注册可更新配置
func (c *configIns) RegisterMixedWithName(name, file, dataID, group string, v IMixedConfig) error {
	if err := checkType(v); err != nil {
		return err
	}

	if _, exist := c.mixed[name]; exist {
		return ErrAlreadyRegister
	}

	if c.client == nil {
		return ErrDialNacosFirst
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	// copy出一个新的空对象，由于存储配置信息
	cp, err := copyAndUnmarshal(data, v)
	if err != nil {
		return err
	}

	mixedConf, ok := cp.(IMixedConfig)
	if !ok {
		return ErrCopyException
	}

	// 从nacos拉取配置信息
	content, err := c.client.GetConfig(vo.ConfigParam{DataId: dataID, Group: group})
	if err != nil {
		return err
	}

	err = mixedConf.OnNacosChanged(c.namespace, group, dataID, content)
	if err != nil {
		return err
	}

	err = c.client.ListenConfig(vo.ConfigParam{
		DataId:   dataID,
		Group:    group,
		OnChange: c.onChange,
	})

	if err != nil {
		return err
	}

	mixedConf.UpdateAfterRegister()

	c.mutex.Lock()
	c.mixed[name] = mixedConf
	c.mutex.Unlock()

	return nil
}

// GetConfig 读取顺序：混合模式 -> 文件模式 -> Nacos模式
func (c *configIns) GetConfig() interface{} {
	conf, _ := c.GetConfigWithFlagByName(defaultName)
	return conf
}

// GetFileConfig 获取文件模式的配置信息
//  为了提升性能，返回值为引用，外部禁止修改
func (c *configIns) GetFileConfig() interface{} {
	return c.GetFileConfigByName(defaultName)
}

// GetNacosConfig 获取nacos模式的配置信息
func (c *configIns) GetNacosConfig() interface{} {
	return c.GetNacosConfigByName(defaultName)
}

// GetMixedConfig 获取混合模式的配置信息
//  为了提升性能，返回值为引用，外部禁止修改
func (c *configIns) GetMixedConfig() interface{} {
	return c.GetMixedConfigByName(defaultName)
}

// GetConfigWithFlagByName 通过名字获取配置信息
//  为了提升性能，返回值为引用，外部禁止修改
func (c *configIns) GetConfigByName(name string) interface{} {
	v, _ := c.GetConfigWithFlagByName(name)
	return v
}

// GetFileConfigByName 获取配置文件中的信息
//  为了提升性能，返回值为引用，外部禁止修改
func (c *configIns) GetFileConfigByName(name string) interface{} {
	if v, exist := c.files[name]; exist {
		return v
	}

	return nil
}

// GetNacosConfigByName 获取nacos模式指定名称的配置信息，值为string
func (c *configIns) GetNacosConfigByName(name string) interface{} {
	conf, exist := c.nacos[name]
	if c.client == nil || !exist {
		return nil
	}

	conf.mutex.RLock()
	content := conf.content
	conf.mutex.RUnlock()

	return content
}

// GetMixedConfigByName 获取混合模式下指定名称的配置信息
func (c *configIns) GetMixedConfigByName(name string) interface{} {
	c.mutex.RLock()
	v, exist := c.mixed[name]
	c.mutex.RUnlock()

	if exist {
		return v
	}

	return nil
}

// GetConfigWithFlag 读取顺序：混合模式 -> 文件模式 -> Nacos模式
func (c *configIns) GetConfigWithFlag() (interface{}, Flag) {
	return c.GetConfigWithFlagByName(defaultName)
}

// GetConfigWithFlagByName 通过名字获取配置信息
func (c *configIns) GetConfigWithFlagByName(name string) (interface{}, Flag) {
	type getter struct {
		f    func(string) interface{}
		flag Flag
	}

	var getters = []getter{
		{c.GetMixedConfigByName, Mixed},
		{c.GetFileConfigByName, OnlyFile},
		{c.GetNacosConfigByName, OnlyNacos},
	}

	for _, gt := range getters {
		v := gt.f(name)
		if v != nil {
			return v, gt.flag
		}
	}

	return nil, Unknown
}

// onChange nacos配置发生变化触发该回调函数
func (c *configIns) onChange(namespace, group, dataID, data string) {
	for _, obj := range c.mixed {
		conf, ok := obj.(IMixedConfig)
		if !ok {
			continue
		}

		c.mutex.Lock()
		conf.OnNacosChanged(namespace, group, dataID, data)
		c.mutex.Unlock()
	}
}

// copyAndUnmarshal 复制一个新对象，然后在进行反序列化
func copyAndUnmarshal(data []byte, v interface{}) (interface{}, error) {
	empty := reflect.New(reflect.Indirect(reflect.ValueOf(v)).Type())
	emptyObject := empty.Interface()

	cp := reflect.New(reflect.Indirect(reflect.ValueOf(v)).Type())
	cpObject := cp.Interface()

	err := yaml.Unmarshal(data, cpObject)
	if err != nil {
		return nil, err
	}

	if reflect.DeepEqual(cpObject, emptyObject) {
		return nil, ErrEmptyConfig
	}

	return cpObject, nil
}

// 检查输入类型是否为指针
func checkType(a interface{}) error {
	if reflect.ValueOf(a).Kind() != reflect.Ptr {
		return ErrMustBePointer
	}

	return nil
}

func NewConfigIns() *configIns {
	return &configIns{
		files: make(map[string]interface{}),
		mixed: make(map[string]interface{}),
		nacos: make(map[string]*nacosConfig),
	}
}
