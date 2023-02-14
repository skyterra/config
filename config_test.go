package config_test

import (
	"config"
	"config/primitive"
	"fmt"
	"time"

	"gopkg.in/yaml.v2"

	. "github.com/onsi/gomega"

	. "github.com/onsi/ginkgo"
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

// MongoConf 配置项
type MongoConf struct {
	DB              string        `yaml:"db"`                 // 数据库名称
	Host            string        `yaml:"host"`               // db连接地址
	Timeout         uint64        `yaml:"timeout"`            // 单个操作执行的最大耗时时长，毫秒
	MaxPoolSize     uint64        `yaml:"max_pool_size"`      // 连接池最大活跃连接数
	MinPoolSize     uint64        `yaml:"min_pool_size"`      // 连接池最小活跃连接数
	MaxConnIdleTime time.Duration `yaml:"max_conn_idle_time"` // 空闲超时时间秒，超时后关闭连接
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

func (c *GeneralConfig) UpdateAfterRegister() {
	c.Configs[c.CurEnv].LogLevel = "info"
}

func (c *GeneralConfig) OnNacosChanged(namespace, group, dataId, data string) error {
	allConf := &NacosAllDBConf{}
	err := yaml.Unmarshal([]byte(data), allConf)
	if err != nil {
		return err
	}

	mongoMonkeyConf := allConf.Mongo.Write.Monkey
	monkeyDB := mongoMonkeyConf.Database
	monkeyHost := fmt.Sprintf("mongodb://%s:%s@%s/%s?readPreference=%s", mongoMonkeyConf.User, mongoMonkeyConf.Password, mongoMonkeyConf.Addr, mongoMonkeyConf.Database, mongoMonkeyConf.ReadPreference)

	c.Configs[c.CurEnv].MonkeyMongo.DB = monkeyDB
	c.Configs[c.CurEnv].MonkeyMongo.Host = monkeyHost
	return nil
}

const (
	addr      = ""
	namespace = ""
	accessKey = ""
	secretKey = ""

	dataID = ""
	group  = ""

	dbDataID = ""
	dbGroup  = ""
)

var _ = Describe("Api", func() {
	Context("file mode", func() {
		It("default name & custom name", func() {
			c := config.NewConfig()

			err := c.RegisterFile("demo.yaml", &GeneralConfig{})
			Expect(err).Should(Succeed())

			conf := c.GetFileConfig()
			Expect(conf.(*GeneralConfig).Configs["local"].ProductName == "config").Should(BeTrue())

			_, flag := c.GetConfigWithFlag()
			Expect(flag == primitive.OnlyFile).Should(BeTrue())

			err = c.RegisterFileWithName("app", "demo.yaml", &GeneralConfig{})
			Expect(err).Should(Succeed())

			conf = c.GetFileConfigByName("app")
			Expect(conf.(*GeneralConfig).Configs["local"].ProductName == "config").Should(BeTrue())

			_, flag = c.GetConfigWithFlag()
			Expect(flag == primitive.OnlyFile).Should(BeTrue())
		})

		It("no pointer type", func() {
			c := config.NewConfig()

			err := c.RegisterFile("demo.yaml", GeneralConfig{})
			Expect(err).ShouldNot(Succeed())
		})
	})

	XContext("nacos mode", func() {
		It("should be succeed", func() {
			c := config.NewConfig()

			err := c.DailNacos(addr, namespace, primitive.WithSecretKey(secretKey), primitive.WithAccessKey(accessKey))
			Expect(err).Should(Succeed())

			err = c.RegisterNacos(dataID, group)
			Expect(err).Should(Succeed())

			data := c.GetConfig()
			content := data.(string)
			Expect(content != "").Should(BeTrue())

			data1, flag := c.GetConfigWithFlag()
			Expect(data1 == data).Should(BeTrue())
			Expect(flag == primitive.OnlyNacos).Should(BeTrue())
		})

		It("not dial nacos", func() {
			c := config.NewConfig()
			data := c.GetConfig()
			Expect(data == nil).Should(BeTrue())
		})
	})

	XContext("mixed mode", func() {
		It("should be succeed", func() {
			c := config.NewConfig()

			err := c.DailNacos(addr, namespace, primitive.WithSecretKey(secretKey), primitive.WithAccessKey(accessKey))
			Expect(err).Should(Succeed())

			err = c.RegisterFile("demo.yaml", &GeneralConfig{})
			Expect(err).Should(Succeed())

			err = c.RegisterMixed("demo.yaml", dbDataID, dbGroup, &GeneralConfig{})
			Expect(err).Should(Succeed())

			// time.Sleep(10 * time.Second)

			fileConf := c.GetFileConfig()
			fileHost := fileConf.(*GeneralConfig).Configs["local"].MonkeyMongo.Host

			conf := c.GetConfig()
			Expect(conf.(*GeneralConfig).Configs["local"].MonkeyMongo.Host != fileHost).Should(BeTrue())
			Expect(conf.(*GeneralConfig).Configs["local"].LogLevel == "info").Should(BeTrue())

			fmt.Println(conf.(*GeneralConfig).Configs["local"])
		})
	})

})
