package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("filesystems", map[string]any{
		// Default Filesystem Disk
		//
		// Here you may specify the default filesystem disk that should be used
		// by the framework. The "local" disk, as well as a variety of cloud
		// based disks are available to your application. Just store away!
		"default": config.Env("FILESYSTEM_DISK", "local"),

		// Filesystem Disks
		//
		// Here you may configure as many filesystem "disks" as you wish, and you
		// may even configure multiple disks of the same driver. Defaults have
		// been set up for each driver as an example of the required values.
		//
		// Supported Drivers: "local", "s3", "oss", "cos", "custom"
		// local : set path to local storage/app
		// uploads : set path to local public/uploads
		// module  : set path to local module -- for module exec path
		"disks": map[string]any{
			"local": map[string]any{
				"driver": "local",
				"root":   "storage/app",
				"url":    config.Env("APP_URL").(string) + "/storage",
			},
			"uploads": map[string]any{
				"driver": "local",
				"root":   "public/uploads",
				"url":    config.Env("APP_URL").(string) + "/public/uploads",
			},
			"module": map[string]any{
				"driver": "local",
				"root":   "module",
				"url":    config.Env("APP_URL").(string) + "/module",
			},
		},
	})
}
