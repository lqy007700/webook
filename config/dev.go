//go:build !k8s

package config

var Config = &config{
	DB: DBConfig{
		DSN: "root:root@tcp(127.0.0.1:13306)/webook",
	},
	Redis: RedisConfig{
		Addr: "127.0.0.1:16379",
	},
}
