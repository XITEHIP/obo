package plugins

import (
	"errors"
	"github.com/xitehip/obo/define"
)

func Register(s *define.Session, key string, handle define.Handle) error {

	if s.Handles == nil {
		s.Handles = make(map[string]define.Handle)
	}
	if _, ok := s.Handles[key]; ok {
		return errors.New("The plugin key " + key + "is exist")
	}
	s.Handles[key] = handle

	return nil
}

func Fire(s *define.Session, message *define.ReceiveMessage, t define.TransmitFun) {
	if len(s.PluginsManager.Handles) > 0 {
		for _, plugin := range s.Handles {
			f := plugin(s, message)
			if f != nil && t != nil {
				f(t)//Message transmit
			}
		}
	}
}
