package configuration

type AppConfiguration struct {
	Main   MainConfig `toml:"main"`
	Protoc Protoc     `toml:"protoc"`
}

type MainConfig struct {
	// ProtoRootPath path to a folder with proto files
	ProtoRootPath string
	// ProtoOutputRootPath output dir for generated files
	ProtoOutputRootPath string
	// ProtoOutputRootModuleName base module name for proto files in output folder
	ProtoOutputRootModuleName string
}

type Protoc struct {
	// absolute path to protoc binary
	ProtocExecutable string
	// working dir for protoc process
	ProtocWorkingDir string
}
