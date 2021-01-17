package grframework

type Server interface {
	Register(funcPath string, h interface{}) //TODO: h is can be umarshal)
	ListenAndBlock(addr string)              // 127.0.0.1
}
