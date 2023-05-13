package main

import (
	"errors"
	"github.com/BurntSushi/toml"
	"github.com/TrueGameover/prototool/prototool/application/build"
	"github.com/TrueGameover/prototool/prototool/domain/configuration"
	"github.com/TrueGameover/prototool/prototool/domain/files"
	"github.com/TrueGameover/prototool/prototool/domain/module"
	"github.com/TrueGameover/prototool/prototool/domain/protoc"
	"github.com/gookit/gcli/v3"
	"os"
	"os/exec"
	"path/filepath"
)

var appConfiguration configuration.AppConfiguration

func main() {
	gcliApp := gcli.NewApp()
	gcliApp.Version = "0.0.1"
	gcliApp.Desc = "tool for handling proto files without pain"

	err := tryGetConfig()
	if err != nil {
		panic(err)
	}

	// commands
	gcliApp.AddCommand(newBuildCommand())

	gcliApp.Run(nil)
}

func tryGetConfig() error {
	_, err := toml.DecodeFile("prototool.toml", &appConfiguration)
	if err != nil {
		return err
	}

	_, err = os.Stat("prototool.local.toml")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err == nil {
		_, err = toml.DecodeFile("prototool.local.toml", &appConfiguration)
		if err != nil {
			return err
		}
	}

	err = tryAutoResolveProtocConf()
	if err != nil {
		return err
	}

	return nil
}

func tryAutoResolveProtocConf() error {
	protocConf := appConfiguration.Protoc

	if len(protocConf.ProtocExecutable) == 0 {
		path, err := exec.LookPath("protoc")
		if err != nil {
			return err
		}

		protocConf.ProtocExecutable = path
	}

	if len(protocConf.ProtocWorkingDir) == 0 {
		abs, err := filepath.Abs("")
		if err != nil {
			return err
		}

		protocConf.ProtocWorkingDir = abs
	}

	appConfiguration.Protoc = protocConf

	return nil
}

func newBuildCommand() *gcli.Command {
	return &gcli.Command{
		Name: "build",
		Desc: "build proto files",
		Func: func(c *gcli.Command, args []string) error {
			buildHandler := build.NewBuildHandler(
				files.NewSearchFiles(),
				module.NewPackageNameGenerator(),
				protoc.NewProtocCommandGenerator(),
				files.NewMoveDir(),
			)

			err := buildHandler.Handle(c.Ctx, appConfiguration)
			if err != nil {
				return err
			}

			return nil
		},
		Help: "",
	}
}
