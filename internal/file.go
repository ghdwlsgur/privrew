package internal

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func ExtractTarGz(srcFile, destDir string, excludeFile string) error {
	file, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer file.Close()

	gzf, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzf.Close()

	tarReader := tar.NewReader(gzf)

	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		return err
	}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		targetPath := filepath.Join(destDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(targetPath, header.FileInfo().Mode().Perm())
			if err != nil {
				return err
			}
		case tar.TypeReg:
			if header.Name == excludeFile {
				targetPath = destDir + "/bin/" + header.Name
			}
			file, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY, header.FileInfo().Mode().Perm())
			if err != nil {
				return err
			}

			_, err = io.Copy(file, tarReader)
			file.Close()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func CreateSymbolicLink(sourcePath, targetPath string) error {
	cmd := exec.Command("ln", "-sf", sourcePath, targetPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
