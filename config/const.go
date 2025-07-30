package config

type EnvType string

const (
	EnvTypeProd EnvType = "prod"
	EnvTypeDev EnvType = "dev"
	EnvTypeTest	EnvType = "test"
	EnvTypeUnitTest EnvType = "unit_test"
)

// error type
type ErrType uint32

const (
	ErrTypeTimeout ErrType = iota
)

var ErrTypeName = map[ErrType]string{
	ErrTypeTimeout: "timeout",
}

var ErrTypeMsg = map[ErrType]string{
	ErrTypeTimeout: "请求超时",
}
