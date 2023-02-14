package internal

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestInternal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Internal Suite")
}

var _ = AfterSuite(func() {
	// c := NewConfigIns()

	// 	err := c.DailNacos(addr, namespace, primitive.WithSecretKey(wSecretKey), primitive.WithAccessKey(wAccessKey))
	// 	Expect(err).Should(Succeed())
	//
	// 	c.client.PublishConfig(vo.ConfigParam{
	// 		DataId: dataID,
	// 		Group:  group,
	// 		Content: `db: monkey
	// host: mongodb://server:27017/monkey
	// timeout: 10000
	// max_pool_size: 150
	// min_pool_size: 1024
	// max_conn_idle_time: 300`,
	// 	})
})
