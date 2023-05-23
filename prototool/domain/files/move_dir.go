package files

import (
	"io"
	"os"
	"path"
)

type MoveDir struct {
}

func NewMoveDir() MoveDir {
	return MoveDir{}
}

func (h *MoveDir) MoveFolderContent(sourcePath string, targetPath string) error {
	dirEntities, err := os.ReadDir(sourcePath)
	if err != nil {
		return err
	}

	stat, err := os.Stat(targetPath)
	if (err != nil && os.IsNotExist(err)) || !stat.IsDir() {
		err = os.Mkdir(targetPath, 0755)
		if err != nil {
			return err
		}
	}

	for _, dirEntity := range dirEntities {
		if dirEntity.IsDir() {
			nextSourcePath := path.Join(sourcePath, dirEntity.Name())
			nextTargetPath := path.Join(targetPath, dirEntity.Name())

			err = h.MoveFolderContent(nextSourcePath, nextTargetPath)
			if err != nil {
				return err
			}

		} else {
			sourceFilePath := path.Join(sourcePath, dirEntity.Name())
			targetFilePath := path.Join(targetPath, dirEntity.Name())

			err = cloneFile(sourceFilePath, targetFilePath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func cloneFile(source string, target string) error {
	sourceFile, err := os.OpenFile(source, os.O_RDONLY, 0444)
	if err != nil {
		return err
	}
	defer func() {
		_ = sourceFile.Close()
	}()

	targetFile, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer func() {
		_ = targetFile.Close()
	}()

	_, err = io.Copy(targetFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}
