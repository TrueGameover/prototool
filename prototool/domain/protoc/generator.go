package protoc

import "strings"

type ProtocCommandGenerator struct {
}

func NewProtocCommandGenerator() ProtocCommandGenerator {
	return ProtocCommandGenerator{}
}

type ProtoFileConfig struct {
	Path       string
	ModuleName string
}

func (h *ProtocCommandGenerator) Generate(
	protoFilesRootPath string,
	outputPath string,
	filesConfig []ProtoFileConfig,
) ([]string, error) {
	protocArguments := make([]string, 0)
	protocArguments = append(protocArguments, "--proto_path="+protoFilesRootPath)
	protocArguments = append(protocArguments, "--go_out="+outputPath)

	goOpts := generateOpts(protoFilesRootPath, "--go_opt=M", filesConfig)
	protocArguments = append(protocArguments, goOpts...)

	protocArguments = append(protocArguments, "--go-grpc_out="+outputPath)

	goGrpcOpts := generateOpts(protoFilesRootPath, "--go-grpc_opt=M", filesConfig)
	protocArguments = append(protocArguments, goGrpcOpts...)

	for _, config := range filesConfig {
		relativeRootPath := strings.Replace(config.Path, protoFilesRootPath, "", 1)
		relativeRootPath = strings.Trim(relativeRootPath, "/")

		protocArguments = append(protocArguments, relativeRootPath)
	}

	return protocArguments, nil
}

func generateOpts(protoRootPath string, prefix string, filesConfig []ProtoFileConfig) []string {
	result := make([]string, len(filesConfig), len(filesConfig))

	for i, config := range filesConfig {
		relativeRootPath := strings.Replace(config.Path, protoRootPath, "", 1)
		relativeRootPath = strings.Trim(relativeRootPath, "/")
		result[i] = prefix + relativeRootPath + "=" + config.ModuleName
	}

	return result
}
