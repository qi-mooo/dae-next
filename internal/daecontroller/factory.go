package daecontroller

import (
	"github.com/daeuniverse/dae/cmd"
	"github.com/daeuniverse/dae/controlapi"
)

func New(runtime cmd.ExternalControllerRuntime, broker *controlapi.LogBroker) (cmd.ExternalControllerHandle, error) {
	provider, err := controlapi.NewDaeProvider(runtime.Version, runtime.Config, runtime.ControlPlane, runtime.ConfigPath, runtime.Loggers...)
	if err != nil {
		return nil, err
	}
	server := controlapi.NewServer(controlapi.ServerConfig{
		Addr:   runtime.Config.Global.ExternalController,
		Secret: runtime.Config.Global.ExternalControllerSecret,
	}, provider, broker)
	if err := server.Start(); err != nil {
		return nil, err
	}
	return server, nil
}
