package primitive

import "errors"

var (
	ErrAlreadyRegister               = errors.New("name had been registered")
	ErrMustBePointer                 = errors.New("need a pointer type")
	ErrDialNacosFirst                = errors.New("need to dial nacos first")
	ErrCopyException                 = errors.New("copy new mixed object exception")
	ErrDataIDAndGroupAlreadyRegister = errors.New("dataID and group had been registered")
	ErrEmptyConfig                   = errors.New("empty config object")
	ErrConnectFailed                 = errors.New("connect nacos server failed")
	ErrNotExistConfig                = errors.New("not exist configure on nacos server")
)
