package internal

import (
	"config/primitive"
	"fmt"
	"sync"
	"time"

	"github.com/nacos-group/nacos-sdk-go/vo"

	"gopkg.in/yaml.v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	addr      = ""
	namespace = ""
	accessKey = ""
	secretKey = ""

	wAccessKey = ""
	wSecretKey = ""

	dataID = "config"
	group  = "config"
)

// GeneralConfig 按照环境区分配置
type GeneralConfig struct {
	CurEnv  string                `yaml:"cur_env"`
	Configs map[string]*Configure `yaml:"envs"`
}

// Configure 具体配置项
type Configure struct {
	Port        int    `yaml:"port"`
	LogLevel    string `yaml:"log_level"`
	ProductName string `yaml:"product_name"`
	GinPProf    bool   `yaml:"gin_pprof"`

	MonkeyMongo MongoConf `yaml:"monkey_mongo"`
	Redis       RedisConf `yaml:"redis"`
	NacosConf   NacosConf `yaml:"nacos"`
}

type NacosConf struct {
	Addr      string `yaml:"addr"`
	Namespace string `yaml:"namespace"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	DataId    string `yaml:"data_id"`
	Group     string `yaml:"group"`
}

// RedisConf redis配置项
type RedisConf struct {
	Host     string `yaml:"host"`
	UserName string `yaml:"user_name"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`

	// MaxIdle 最大空闲连接数
	MaxIdle int `yaml:"max_idle"`

	// IdleTimeout 连接空闲多久会被关闭，单位毫秒
	IdleTimeout int `yaml:"idle_timeout"`
	MaxActive   int `yaml:"max_active"`

	// Wait 最大活动连接到达上限是否要等待
	Wait bool `yaml:"wait"`

	// ReadTimeout 读网络包超时时间,单位毫秒, 0代表没有超时
	ReadTimeout time.Duration `yaml:"read_timeout"`

	// WriteTimeout 写网络包超时时间,单位毫秒,0代表没有超时
	WriteTimeout time.Duration `yaml:"write_timeout"`

	// ConnectTimeout 建连接超时时间,单位毫秒，0代表没有超时时间
	ConnectTimeout time.Duration `yaml:"connect_timeout"`

	// GetConnectTimeout 从连接池拿连接等待时间,,单位毫秒，0代表没有超时时间
	GetConnectTimeout time.Duration `yaml:"get_connect_timeout"`
}

type NacosAllDBConf struct {
	Mongo struct {
		Write struct {
			Monkey NacosDBConf `yaml:"monkey"`
			Whale  NacosDBConf `yaml:"whale"`
		}
		Read struct {
			Monkey NacosDBConf `yaml:"monkey"`
			Whale  NacosDBConf `yaml:"whale"`
		}
	}
}

type NacosDBConf struct {
	Database       string `yaml:"database"`
	Addr           string `yaml:"addr"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	ReadPreference string `yaml:"readPreference"`
}

// MongoConf 配置项
type MongoConf struct {
	DB              string        `yaml:"db"`                 // 数据库名称
	Host            string        `yaml:"host"`               // db连接地址
	Timeout         uint64        `yaml:"timeout"`            // 单个操作执行的最大耗时时长，毫秒
	MaxPoolSize     uint64        `yaml:"max_pool_size"`      // 连接池最大活跃连接数
	MinPoolSize     uint64        `yaml:"min_pool_size"`      // 连接池最小活跃连接数
	MaxConnIdleTime time.Duration `yaml:"max_conn_idle_time"` // 空闲超时时间秒，超时后关闭连接
}

// MongoConf 配置项
type MongoConf2 struct {
	DB              string        `yaml:"db"`                 // 数据库名称
	Host            string        `yaml:"host"`               // db连接地址
	Timeout         uint64        `yaml:"timeout"`            // 单个操作执行的最大耗时时长，毫秒
	MaxPoolSize     uint64        `yaml:"max_pool_size"`      // 连接池最大活跃连接数
	MinPoolSize     uint64        `yaml:"min_pool_size"`      // 连接池最小活跃连接数
	MaxConnIdleTime time.Duration `yaml:"max_conn_idle_time"` // 空闲超时时间秒，超时后关闭连接
}

// MongoConf 配置项
type MongoConf3 struct {
	DB              string        `yaml:"db"`                 // 数据库名称
	Host            string        `yaml:"host"`               // db连接地址
	Timeout         uint64        `yaml:"timeout"`            // 单个操作执行的最大耗时时长，毫秒
	MaxPoolSize     uint64        `yaml:"max_pool_size"`      // 连接池最大活跃连接数
	MinPoolSize     uint64        `yaml:"min_pool_size"`      // 连接池最小活跃连接数
	MaxConnIdleTime time.Duration `yaml:"max_conn_idle_time"` // 空闲超时时间秒，超时后关闭连接
}

func (mc *MongoConf) UpdateAfterRegister() {
	mc.MinPoolSize = 1024
}

func (mc *MongoConf) OnNacosChanged(namespace, group, dataId, data string) error {
	conf := &MongoConf{}
	err := yaml.Unmarshal([]byte(data), conf)
	if err != nil {
		return err
	}

	*mc = *conf
	return nil
}

func (mc *MongoConf2) UpdateAfterRegister() {
	mc.MinPoolSize = 2048
}

func (mc *MongoConf2) OnNacosChanged(namespace, group, dataId, data string) error {
	conf := &MongoConf2{}
	err := yaml.Unmarshal([]byte(data), conf)
	if err != nil {
		return err
	}

	*mc = *conf
	return nil
}

func (mc MongoConf3) UpdateAfterRegister() {
	mc.MinPoolSize = 1024
}

func (mc MongoConf3) OnNacosChanged(namespace, group, dataId, data string) error {
	conf := &MongoConf3{}
	err := yaml.Unmarshal([]byte(data), conf)
	if err != nil {
		return err
	}

	mc = *conf
	return nil
}

func (cc *GeneralConfig) GetConfig() *Configure {
	return cc.Configs[cc.CurEnv]
}

var _ = Describe("Ins", func() {
	Context("file mode", func() {
		It("default & custom name", func() {
			c := NewConfigIns()
			err := c.RegisterFile("default.yaml", &GeneralConfig{})
			Expect(err).Should(Succeed())

			defaultConf := c.GetFileConfig().(*GeneralConfig)
			Expect(defaultConf.GetConfig().ProductName == "default").Should(BeTrue())

			err = c.RegisterFileWithName("app", "app.yaml", &GeneralConfig{})
			Expect(err).Should(Succeed())

			appConf := c.GetFileConfigByName("app").(*GeneralConfig)
			Expect(appConf.GetConfig().ProductName == "app").Should(BeTrue())

			Expect(appConf.GetConfig().ProductName != defaultConf.GetConfig().ProductName).Should(BeTrue())
		})

		It("type check", func() {
			c := NewConfigIns()
			err := c.RegisterFile("default.yaml", GeneralConfig{})
			Expect(err).ShouldNot(Succeed())

			err = c.RegisterFile("default.yaml", nil)
			Expect(err).ShouldNot(Succeed())
		})

		It("file not exist", func() {
			c := NewConfigIns()
			err := c.RegisterFile("unknown.yaml", &GeneralConfig{})
			Expect(err).ShouldNot(Succeed())
		})

		It("duplicate name", func() {
			c := NewConfigIns()
			err := c.RegisterFile("app.yaml", &GeneralConfig{})
			Expect(err).Should(Succeed())

			err = c.RegisterFile("app.yaml", &GeneralConfig{})
			Expect(err).ShouldNot(Succeed())
		})

		It("empty config", func() {
			type UnknownConfig struct {
				Name string `yaml:"name"`
			}

			c := NewConfigIns()
			err := c.RegisterFile("app.yaml", &UnknownConfig{})
			Expect(err).ShouldNot(Succeed())
		})

		It("empty config", func() {
			type UnknownConfig struct {
				Name string `yaml:"name"`
			}

			c := NewConfigIns()
			err := c.RegisterFile("app.yaml", &UnknownConfig{})
			Expect(err).ShouldNot(Succeed())
		})

		It("json config", func() {
			type UnknownConfig struct {
				Name string `yaml:"name"`
			}

			c := NewConfigIns()
			err := c.RegisterFile("default.unknown", &UnknownConfig{})
			Expect(err).ShouldNot(Succeed())
		})
	})

	XContext("nacos", func() {
		It("should be succeed", func() {
			c := NewConfigIns()

			err := c.DailNacos(addr, namespace, primitive.WithSecretKey(secretKey), primitive.WithAccessKey(accessKey))
			Expect(err).Should(Succeed())

			err = c.RegisterNacos(dataID, group)
			Expect(err).Should(Succeed())

			data := c.GetNacosConfig()
			content := data.(string)

			Expect(c.GetNacosConfigByName("app") == nil).Should(BeTrue())

			conf := &MongoConf{}
			err = yaml.Unmarshal([]byte(content), conf)
			Expect(err).Should(Succeed())

			Expect(conf.Host == "mongodb://server:27017/monkey").Should(BeTrue())

			data1, flag := c.GetConfigWithFlag()
			Expect(data1 == data).Should(BeTrue())
			Expect(flag == primitive.OnlyNacos).Should(BeTrue())
		})

		It("listen changed", func() {
			clearRegistered()

			c := NewConfigIns()

			err := c.DailNacos(addr, namespace, primitive.WithSecretKey(wSecretKey), primitive.WithAccessKey(wAccessKey))
			Expect(err).Should(Succeed())

			err = c.RegisterNacos(dataID, group)
			Expect(err).Should(Succeed())

			data := c.GetNacosConfig()
			content := data.(string)

			oriConf := &MongoConf{}
			err = yaml.Unmarshal([]byte(content), oriConf)
			Expect(err).Should(Succeed())

			cc := *oriConf
			cc.Timeout = 10000

			wg := sync.WaitGroup{}
			wg.Add(2)

			go func() {
				for i := 0; i < 5; i++ {
					cc.Timeout = cc.Timeout + 1

					data, _ := yaml.Marshal(cc)
					c.client.PublishConfig(vo.ConfigParam{
						DataId:  dataID,
						Group:   group,
						Content: string(data),
					})
					time.Sleep(500 * time.Millisecond)
				}
				wg.Done()
			}()

			go func() {
				for i := 0; i < 5; i++ {
					time.Sleep(600 * time.Millisecond)
					content := c.GetNacosConfig()
					conf := &MongoConf{}
					err = yaml.Unmarshal([]byte(content.(string)), conf)
					Expect(err).Should(Succeed())

					Expect(conf.Timeout >= oriConf.Timeout).Should(BeTrue())
					fmt.Println(conf.Timeout)
				}
				wg.Done()
			}()

			wg.Wait()

			originData, _ := yaml.Marshal(oriConf)
			c.client.PublishConfig(vo.ConfigParam{
				DataId:  dataID,
				Group:   group,
				Content: string(originData),
			})
		})

		It("empty nacos address", func() {
			clearRegistered()
			c := NewConfigIns()
			err := c.DailNacos("", namespace, primitive.WithTimeoutMs(100))
			Expect(err).ShouldNot(Succeed())
		})

		It("invalidate address", func() {
			clearRegistered()
			c := NewConfigIns()
			err := c.DailNacos("127.0.0.1", namespace, primitive.WithTimeoutMs(100))
			Expect(err).ShouldNot(Succeed())
		})

		It("dail twice", func() {
			clearRegistered()

			c := NewConfigIns()
			err := c.DailNacos(addr, namespace, primitive.WithAccessKey(accessKey), primitive.WithSecretKey(secretKey))
			Expect(err).Should(Succeed())

			err = c.DailNacos(addr, namespace, primitive.WithAccessKey(accessKey), primitive.WithSecretKey(secretKey))
			Expect(err).ShouldNot(Succeed())
		})

		It("register dataID & group twice", func() {
			clearRegistered()

			c := NewConfigIns()

			err := c.DailNacos(addr, namespace, primitive.WithSecretKey(wSecretKey), primitive.WithAccessKey(wAccessKey))
			Expect(err).Should(Succeed())

			err = c.RegisterNacos(dataID, group)
			Expect(err).Should(Succeed())

			err = c.RegisterNacosWithName("app", dataID, group)
			Expect(err).ShouldNot(Succeed())

			err = c.RegisterNacos(dataID, group)
			Expect(err).ShouldNot(Succeed())
		})

		It("register invalidate dataID or group", func() {
			clearRegistered()

			c := NewConfigIns()

			err := c.DailNacos(addr, namespace, primitive.WithSecretKey(wSecretKey), primitive.WithAccessKey(wAccessKey))
			Expect(err).Should(Succeed())

			err = c.RegisterNacos("", group)
			Expect(err).ShouldNot(Succeed())

			err = c.RegisterNacosWithName("app", dataID, "")
			Expect(err).ShouldNot(Succeed())

			err = c.RegisterNacosWithName("app", dataID, "unknown")
			Expect(err).ShouldNot(Succeed())
		})
	})

	XContext("mixed mode", func() {
		It("normal", func() {
			clearRegistered()
			c := NewConfigIns()

			err := c.DailNacos(addr, namespace, primitive.WithSecretKey(secretKey), primitive.WithAccessKey(accessKey))
			Expect(err).Should(Succeed())

			err = c.RegisterMixed("mixed.yaml", dataID, group, &MongoConf{})
			Expect(err).Should(Succeed())

			conf := c.GetMixedConfig().(*MongoConf)
			Expect(conf.MinPoolSize == 1024).Should(BeTrue())
			Expect(conf.MaxPoolSize == 150).Should(BeTrue())
			Expect(conf.Host == "mongodb://server:27017/monkey").Should(BeTrue())

			err = c.RegisterMixedWithName("app", "mixed.yaml", dataID, group, &MongoConf2{})
			Expect(err).Should(Succeed())

			conf2 := c.GetMixedConfigByName("app").(*MongoConf2)
			Expect(conf2.MinPoolSize == 2048).Should(BeTrue())
			Expect(conf2.MaxPoolSize == 150).Should(BeTrue())
			Expect(conf2.Host == "mongodb://server:27017/monkey").Should(BeTrue())

			Expect(conf2.MinPoolSize != conf.MinPoolSize).Should(BeTrue())
		})

		It("no dial nacos", func() {
			clearRegistered()
			c := NewConfigIns()
			err := c.RegisterMixed("mixed.yaml", dataID, group, &MongoConf{})
			Expect(err).ShouldNot(Succeed())
		})

		It("duplicate name", func() {
			clearRegistered()
			c := NewConfigIns()
			err := c.DailNacos(addr, namespace, primitive.WithSecretKey(secretKey), primitive.WithAccessKey(accessKey))
			Expect(err).Should(Succeed())

			err = c.RegisterMixedWithName("app", "mixed.yaml", dataID, group, &MongoConf{})
			Expect(err).Should(Succeed())

			err = c.RegisterMixedWithName("app", "mixed.yaml", dataID, group, &MongoConf{})
			Expect(err).ShouldNot(Succeed())
		})

		It("invalidate type", func() {
			clearRegistered()
			c := NewConfigIns()

			err := c.DailNacos(addr, namespace, primitive.WithSecretKey(secretKey), primitive.WithAccessKey(accessKey))
			Expect(err).Should(Succeed())

			err = c.RegisterMixedWithName("app", "mixed.yaml", dataID, group, MongoConf3{})
			Expect(err).ShouldNot(Succeed())

			err = c.RegisterMixedWithName("app1", "mixed111.yaml", dataID, group, &MongoConf{})
			Expect(err).ShouldNot(Succeed())
		})

		It("listen changed", func() {
			clearRegistered()

			c := NewConfigIns()

			err := c.DailNacos(addr, namespace, primitive.WithSecretKey(wSecretKey), primitive.WithAccessKey(wAccessKey))
			Expect(err).Should(Succeed())

			err = c.RegisterMixed("mixed.yaml", dataID, group, &MongoConf{})
			Expect(err).Should(Succeed())

			oriConf := c.GetMixedConfig().(*MongoConf)

			cc := *oriConf
			cc.Timeout = 10000

			wg := sync.WaitGroup{}
			wg.Add(2)

			go func() {
				for i := 0; i < 5; i++ {
					cc.Timeout = cc.Timeout + 1

					data, _ := yaml.Marshal(cc)
					c.client.PublishConfig(vo.ConfigParam{
						DataId:  dataID,
						Group:   group,
						Content: string(data),
					})
					time.Sleep(500 * time.Millisecond)
				}
				wg.Done()
			}()

			go func() {
				for i := 0; i < 5; i++ {
					time.Sleep(600 * time.Millisecond)
					conf := c.GetMixedConfig().(*MongoConf)

					Expect(conf.Timeout >= oriConf.Timeout).Should(BeTrue())
					fmt.Println(conf.Timeout)
				}
				wg.Done()
			}()

			wg.Wait()

			originData, _ := yaml.Marshal(oriConf)
			c.client.PublishConfig(vo.ConfigParam{
				DataId:  dataID,
				Group:   group,
				Content: string(originData),
			})
		})
	})

	Context("get config", func() {
		XIt("normal", func() {
			c := NewConfigIns()

			err := c.DailNacos(addr, namespace, primitive.WithSecretKey(secretKey), primitive.WithAccessKey(accessKey))
			Expect(err).Should(Succeed())

			err = c.RegisterFile("mixed.yaml", &MongoConf{})
			Expect(err).Should(Succeed())

			err = c.RegisterMixed("mixed.yaml", dataID, group, &MongoConf{})
			Expect(err).Should(Succeed())

			fileConf := c.GetFileConfig().(*MongoConf)
			mixedConf := c.GetMixedConfig().(*MongoConf)
			Expect(fileConf.MinPoolSize != mixedConf.MinPoolSize).Should(BeTrue())

			cc := c.GetConfig().(*MongoConf)
			Expect(cc == mixedConf).Should(BeTrue())

			_, flag := c.GetConfigWithFlag()
			Expect(flag == primitive.Mixed).Should(BeTrue())

			err = c.RegisterFileWithName("app", "mixed.yaml", &MongoConf{})
			Expect(err).Should(Succeed())

			appFc := c.GetFileConfigByName("app").(*MongoConf)
			appCc := c.GetConfigByName("app").(*MongoConf)

			Expect(appFc == appCc).Should(BeTrue())

		})
	})
})
