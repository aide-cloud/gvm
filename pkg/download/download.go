package download

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// DownloadFile 下载文件
func FetchFile(url, destPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 创建目标文件
	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// ExtractTarGz 解压 tar.gz 文件
func ExtractGoSdkTarGzFile(srcPath, destPath string) error {
	// 打开 tar.gz 文件
	file, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	// 遍历 tar 文件
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // 读取完毕
		}
		if err != nil {
			return err
		}

		// 去掉 go/ 根目录前缀
		relativePath := strings.TrimPrefix(header.Name, "go/")

		// 跳过空的路径（根目录）
		if relativePath == "" {
			continue
		}

		target := filepath.Join(destPath, relativePath)

		switch header.Typeflag {
		case tar.TypeDir: // 目录
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg: // 文件
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			outFile, err := os.Create(target)
			if err != nil {
				return err
			}
			_, err = io.Copy(outFile, tarReader)
			outFile.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
