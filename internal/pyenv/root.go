package pyenv

import (
	"archive/tar"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/zstd"
)

func Main() {
	downloadPath := "cpython-3.12.3+20240415-aarch64-apple-darwin-pgo+lto-full.tar.zst"
	downloadUrl := fmt.Sprintf("https://github.com/indygreg/python-build-standalone/releases/download/20240415/%s", downloadPath)

	r, err := http.Get(downloadUrl)
	if err != nil {
		log.Fatalf("download failed: %v", err)
		os.Exit(1)
	}
	if r.StatusCode == http.StatusNotFound {
		log.Fatal("404 not found")
		os.Exit(1)
	}
	defer r.Body.Close()

	fileData, err := io.ReadAll(r.Body)

	err = os.WriteFile(downloadPath, fileData, 0o640)
	if err != nil {
		log.Fatalf("writing file failed: %v", err)
		os.Remove(downloadPath)
		os.Exit(1)
	}

	if _, err := os.Stat(downloadPath); err == nil {

		extractPath := "apple-darwin.extracted"
		err := os.RemoveAll(extractPath)
		if err != nil {
			log.Panic(err)
		}

		extract(downloadPath, extractPath)

		// installPath := filepath.Join(extractPath, "python", "install")

		// var libPath string
		// libPath = filepath.Join(installPath, "lib", fmt.Sprintf("python%s", pythonVersionBase))
		//
		// for _, lib := range removeLibs {
		// 	_ = os.RemoveAll(filepath.Join(libPath, lib))
		// }
		//
		// err = internal.CleanupPythonDir(installPath, keepPatterns)
		// if err != nil {
		// 	panic(err)
		// }
	}
}

func extract(archivePath string, targetPath string) string {
	f, err := os.Open(archivePath)
	if err != nil {
		log.Fatalf("opening file failed: %v", err)
		os.Exit(1)
	}
	defer f.Close()

	z, err := zstd.NewReader(f)
	if err != nil {
		log.Fatalf("decompression failed: %v", err)
		os.Exit(1)
	}
	defer z.Close()

	log.Printf("decompressing %s\n", archivePath)
	err = ExtractTarStream(z, targetPath)
	if err != nil {
		log.Fatalf("decompression failed: %v", err)
		os.Exit(1)
	}

	return targetPath
}

func ExtractTarStream(r io.Reader, targetPath string) error {
	tarReader := tar.NewReader(r)
	for true {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("ExtractTarStream: Next() failed: %w", err)
		}

		if !validRelPath(header.Name) {
			return fmt.Errorf("tar contained invalid name error %q", header.Name)
		}

		p := filepath.FromSlash(header.Name)
		p = filepath.Join(targetPath, p)

		err = os.MkdirAll(filepath.Dir(p), 0755)
		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(p, 0755); err != nil {
				return fmt.Errorf("ExtractTarStream: Mkdir() failed: %w", err)
			}
		case tar.TypeReg:
			_ = os.Remove(p) // we allow overwriting, which easily happens on case insensitive filesystems
			outFile, err := os.Create(p)
			if err != nil {
				return fmt.Errorf("ExtractTarStream: Create() failed: %w", err)
			}
			_, err = io.Copy(outFile, tarReader)
			_ = outFile.Close()
			if err != nil {
				return fmt.Errorf("ExtractTarStream: Copy() failed: %w", err)
			}
			err = os.Chmod(p, header.FileInfo().Mode())
			if err != nil {
				return fmt.Errorf("ExtractTarStream: Chmod() failed: %w", err)
			}
			err = os.Chtimes(p, header.AccessTime, header.ModTime)
			if err != nil {
				return err
			}
		case tar.TypeSymlink:
			_ = os.Remove(p) // we allow overwriting, which easily happens on case insensitive filesystems
			if err := os.Symlink(header.Linkname, p); err != nil {
				return fmt.Errorf("ExtractTarStream: Symlink() failed: %w", err)
			}
		default:
			return fmt.Errorf("ExtractTarStream: uknown type %v in %v", header.Typeflag, header.Name)
		}
	}
	return nil
}

func validRelPath(p string) bool {
	if p == "" || strings.Contains(p, `\`) || strings.HasPrefix(p, "/") || strings.Contains(p, "../") {
		return false
	}
	return true
}
