package config

import (
	"config/internal"
	. "config/primitive"
)

type IConfig interface {
	DailNacos(addr, namespace string, opts ...ClientOption) error

	RegisterFile(file string, v interface{}) error
	RegisterFileWithName(name, file string, v interface{}) error

	RegisterMixed(file, dataID, group string, v IMixedConfig) error
	RegisterMixedWithName(name, file, dataID, group string, v IMixedConfig) error

	RegisterNacos(dataID, group string) error
	RegisterNacosWithName(name, dataID, group string) error

	GetConfig() interface{}
	GetFileConfig() interface{}
	GetNacosConfig() interface{}
	GetMixedConfig() interface{}

	GetConfigByName(name string) interface{}
	GetFileConfigByName(name string) interface{}
	GetNacosConfigByName(name string) interface{}
	GetMixedConfigByName(name string) interface{}

	GetConfigWithFlag() (interface{}, Flag)
	GetConfigWithFlagByName(name string) (interface{}, Flag)
}

func NewConfig() IConfig {
	return internal.NewConfigIns()
}
