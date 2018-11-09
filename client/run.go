package main

import "github.com/xitehip/obo/plugins"
import "github.com/xitehip/obo/plugins/filehelper"
import (
	"github.com/xitehip/obo/plugins/group"
	"github.com/xitehip/obo/pipe"
)

func main() {

	var plugins []plugins.PluginProviderInterface

	fh := &filehelper.AutoSendService{}
	plugins = append(plugins, fh)

	group := &group.CustomerService{}
	plugins = append(plugins, group)

	pipe := pipe.New().AttachPlugins(plugins)

	pipe.Run()







}
