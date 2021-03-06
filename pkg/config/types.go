package config

type ServerMessage struct {
	IpAddr string
	Port   uint64
}

type ClientMessage struct {
	NamespaceId string
	ServiceName string
	Scope       int
	LogDir      string
	CacheDir    string
	RotateTime  string
	MaxAge      int64
	LogLevel    string
}

type BasicMessage struct {
	InstanceIp             string
	InstancePort           uint64
	SubscribeInstanceCount int
	InstanceClusterName    string
	SubscribeScope         int
	NameServerAddr         string
}

type Config struct {
	Server ServerMessage
	Client ClientMessage
	Basic  BasicMessage
}
