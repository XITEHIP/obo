package plugins

import (
	"github.com/xitehip/obo/define"
)

type PluginProviderInterface interface {
	Register(*define.Session)
}