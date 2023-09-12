package config

import (
	"github.com/goravel/framework/facades"
)

type Module struct {
	Name string
	IP   string
	Port int32
}

func init() {
	config := facades.Config()
	//Module config
	config.Add(
		"module",
		map[string]any{
			// Module uploads path
			"uploads_path": config.Env("MODULE_UPLOADS_PATH", "/public/uploads"),
			// Module exec path.
			"module_path":         config.Env("MODULE_PATH", "/module"),
			"module_port_start":   config.Env("MODULE_PORT_START", 3010), // module port from
			"module_port_current": 3010,                                  //current max port
		},
	)
}
