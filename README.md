# config 配置管理工具
提供多种方式进行配置管理，简化热更新操作；同时支持多配置文件管理，适用于多应用开发；

- 本地配置文件：只读方式，不支持注册后更新
- nacos配置：只读方式，server端配置变更后，会同步更新本地配置；不支持注册后更新
- 混合配置：读写方式，支持nacos热更新，支持注册后更新

## 本地配置文件
本地配置文件管理，用于变化较少（或不需要热更操作）的配置信息管理

Example

```go
package main

import (
	"fmt"
	"config"
)

type YourConfig struct {
	ProductName string `yaml:"product"`
}

func main() {
	conf := config.NewConfig()
	err := c.RegisterFile("demo.yaml", &YourConfig{})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	conf := c.GetFileConfig().(*YourConfig)
	fmt.Println(conf.ProductName)
	
	// your other code...

	// 多应用场景，可以给配置命名
	err = c.RegisterFileWithName("app", "demo.yaml", &YourConfig{})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	
	conf = c.GetFileConfigByName("app").(*YourConfig)
}

```

## nacos方式
通过nacos从服务端拉取配置，通过监听nacos变化，支持热更新

```go
package main

import "config"

func main() {
	c := config.NewConfig()
	
	err := c.DailNacos(addr, namespace, primitive.WithSecretKey(wSecretKey), primitive.WithAccessKey(wAccessKey))
	if err != nil {
		
	}
	
	// 注册nacos配置时，会自动启动监听策略
	err = c.RegisterNacos(dataID, group)
	if err != nil {
		
	}

	data := c.GetNacosConfig()
	content := data.(string)
	fmt.Println(content)
	
	// 需要获取多个nacos配置
	c.RegisterNacosWithName("app", dataID, group)
	content = c.GetNacosConfigByName("app").(string)
}

```

## 多种混合方式
首先提供默认的文件配置，然后，可以通过环境变量或nacos配置进行更新；通过监听nacos变化，进行热更新

```go
package main

func main() {

type YourConfig struct {
	ProductName string `yaml:"product"`
	LogLevel string `yaml:"log_level"`
}

func (c *YourConfig) UpdateAfterRegister() { 
	// 如果需要在注册后更新配置，可以写在这里 
	c.LogLevel = os.GetEnv("LOG_LEVEL")
}

func (c *YourConfig) OnNacosChanged(namespace, group, dataId, data string) error { 
	// nacos配置发生变化，同步变化到本地
	conf := &YourConfig{}
	err := yaml.Unmarshal([]byte(data), conf)
	if err != nil {
		return err
	}

	*mc = *conf
	return nil
}

func main() {
	c := NewConfigIns()

	err := c.DailNacos(addr, namespace, primitive.WithSecretKey(secretKey), primitive.WithAccessKey(accessKey))
	if err != nil {
	
	}
	
	err = c.RegisterMixed("demo.yaml", dataID, group, &YourConf{})
	if err != nil {
	
	}
	
	conf := c.GetMixedConfig().(*YourConf)
}

```

## 通用读取接口 GetConfig()
> 如果你使用多种模式（注册了本地文件，也注册了nacos，还注册了混合模式），此时，调用GetConfig()读取顺序为 混合模式 -> 文件模式 -> Nacos模式