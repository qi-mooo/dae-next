package daecontroller

import (
	daecmd "github.com/daeuniverse/dae/cmd"
	daecontrolapi "github.com/daeuniverse/dae/controlapi"

	nextcontrolapi "github.com/qi-mooo/dae-next/internal/controllerapi"
)

func New(runtime daecmd.ExternalControllerRuntime) (daecmd.ExternalControllerHandle, error) {
	provider, err := daecontrolapi.NewDaeProvider(runtime.Version, runtime.Config, runtime.ControlPlane, runtime.ConfigPath, runtime.Loggers...)
	if err != nil {
		return nil, err
	}
	logBroker, _ := runtime.LogSink.(*nextcontrolapi.LogBroker)
	server := nextcontrolapi.NewServer(nextcontrolapi.ServerConfig{
		Addr:   runtime.Config.Global.ExternalController,
		Secret: runtime.Config.Global.ExternalControllerSecret,
	}, provider, logBroker)
	if err := server.Start(); err != nil {
		return nil, err
	}
	return server, nil
}
