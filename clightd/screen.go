package clightd

import (
	"github.com/godbus/dbus/v5"
)

/** Screen API object **/
const (
	screenInterface         = "org.clightd.clightd.Screen"
	screenObjectPath        = "/org/clightd/clightd/Screen"

	screenMethodGetEmitted  = screenInterface + ".GetEmittedBrightness"
)

type ScreenApi interface {
	ClightdApi
	GoGetEmitted(ch chan *dbus.Call) error
	GetEmitted() (float64, error)
}

func NewScreenApi() (ScreenApi, error) {
	err := ensureXorg()
	if err == nil {
		return initialize(screenObjectPath)
	}
	return nil, err
}

func (api api) GoGetEmitted(ch chan *dbus.Call) error {
	call := api.obj.Go(screenMethodGetEmitted,0, ch, xdisplay, xauth)
	return call.Err
}

func (api api) GetEmitted() (float64, error) {
	call := api.obj.Call(screenMethodGetEmitted,0, xdisplay, xauth)
	if call.Err != nil {
		return 0, call.Err
	}
	var val float64
	err := call.Store(&val)
	return val, err
}
