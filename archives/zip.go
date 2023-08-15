package archives

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func UnzipFile(zipFile, destinationDir string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		// 构建解压后的文件路径
		targetPath := filepath.Join(destinationDir, file.Name)

		if file.FileInfo().IsDir() {
			// 如果是目录，创建对应的目录
			os.MkdirAll(targetPath, os.ModePerm)
			continue
		}

		// 如果是文件，创建并复制文件内容
		if err := os.MkdirAll(filepath.Dir(targetPath), os.ModePerm); err != nil {
			return err
		}

		// 打开ZIP文件中的文件
		srcFile, err := file.Open()
		if err != nil {
			return err
		}
		defer srcFile.Close()

		// 创建目标文件
		dstFile, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		// 复制文件内容
		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			return err
		}
	}

	return nil
}
