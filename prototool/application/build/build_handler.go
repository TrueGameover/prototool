package build

import (
	"context"
	"fmt"
	"github.com/TrueGameover/prototool/prototool/domain/configuration"
	"github.com/TrueGameover/prototool/prototool/domain/files"
	"github.com/TrueGameover/prototool/prototool/domain/module"
	"github.com/TrueGameover/prototool/prototool/domain/protoc"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

type BuildHandler struct {
	searchFiles            files.SearchFiles
	packageNameGenerator   module.PackageNameGenerator
	protocCommandGenerator protoc.ProtocCommandGenerator
	moveDir                files.MoveDir
}

func NewBuildHandler(
	searchFiles files.SearchFiles,
	packageNameGenerator module.PackageNameGenerator,
	protocCommandGenerator protoc.ProtocCommandGenerator,
	moveDir files.MoveDir,
) BuildHandler {
	return BuildHandler{
		searchFiles:            searchFiles,
		packageNameGenerator:   packageNameGenerator,
		protocCommandGenerator: protocCommandGenerator,
		moveDir:                moveDir,
	}
}

func (h *BuildHandler) Handle(ctx context.Context, appConf configuration.AppConfiguration) error {
	protoFilesList, err := h.searchFiles.Find(ctx, appConf.Main.ProtoRootPath)
	if err != nil {
		return err
	}

	packages := h.packageNameGenerator.Generate(protoFilesList, appConf.Main.ProtoRootPath, appConf.Main.ProtoOutputRootModuleName)

	filesConfigs := make([]protoc.ProtoFileConfig, len(packages), len(packages))
	for i, s := range packages {
		filesConfigs[i] = protoc.ProtoFileConfig{
			Path:       protoFilesList[i],
			ModuleName: s,
		}
	}

	tempDir, err := os.MkdirTemp("/tmp", "prototool_*")
	if err != nil {
		return err
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	generatedArgs, err := h.protocCommandGenerator.Generate(appConf.Main.ProtoRootPath, tempDir, filesConfigs)
	if err != nil {
		return err
	}

	protocPath, err := exec.LookPath(appConf.Protoc.ProtocExecutable)
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, protocPath, generatedArgs...)
	cmd.Dir = appConf.Protoc.ProtocWorkingDir
	response, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(response))
		fmt.Println(err)
		return err
	}

	targetDir, err := filepath.Abs(appConf.Main.ProtoOutputRootPath)
	if err != nil {
		return err
	}
	sourceDir := path.Join(tempDir, appConf.Main.ProtoOutputRootModuleName)
	err = h.moveDir.MoveFolderContent(sourceDir, targetDir)
	if err != nil {
		return err
	}

	return nil
}
