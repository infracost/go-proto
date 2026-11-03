// Package plugintest installs released Infracost plugin binaries for tests,
// using the same release layout the CLI installs from.
package plugintest

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	defaultBaseURL = "https://releases.infracost.io"

	// baseURLEnv overrides the release host, as it does for the CLI.
	baseURLEnv = "INFRACOST_CLI_PLUGIN_BASE_URL"

	maxBinarySize = 1 << 30
)

// Install downloads release version of the named plugin for the current OS and
// architecture into dir and returns the binary's path. version is either
// "latest" or a release version, with or without a leading "v".
//
// The archive is fetched from {base}/{name}/{goos}/{goarch}/{version}/, where
// base is $INFRACOST_CLI_PLUGIN_BASE_URL or https://releases.infracost.io, and
// checked against its published .sha256. The checksum is cached next to the
// binary, so a binary already installed at that checksum is not downloaded
// again. Install is safe to call concurrently for different plugins.
func Install(dir, name, version string) (string, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create plugin dir: %w", err)
	}

	binaryName, archiveName := name, "data.tar.gz"
	if runtime.GOOS == "windows" {
		binaryName, archiveName = name+".exe", "data.zip"
	}
	binaryPath := filepath.Join(dir, binaryName)
	shaPath := filepath.Join(dir, name+".sha256")

	baseURL := os.Getenv(baseURLEnv)
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	if version != "latest" {
		version = "v" + strings.TrimPrefix(version, "v")
	}
	archiveURL := strings.Join([]string{strings.TrimRight(baseURL, "/"), name, runtime.GOOS, runtime.GOARCH, version, archiveName}, "/")

	expectedSHA, err := fetchSHA(archiveURL + ".sha256")
	if err != nil {
		return "", fmt.Errorf("fetch sha for %s: %w", name, err)
	}

	if cached, err := os.ReadFile(shaPath); err == nil && strings.TrimSpace(string(cached)) == expectedSHA {
		if _, err := os.Stat(binaryPath); err == nil {
			return binaryPath, nil
		}
	}

	archivePath, err := download(archiveURL, expectedSHA, name)
	if err != nil {
		return "", fmt.Errorf("download %s: %w", name, err)
	}
	defer func() { _ = os.Remove(archivePath) }()

	extract := extractTarGz
	if archiveName == "data.zip" {
		extract = extractZip
	}
	if err := extract(archivePath, binaryPath, binaryName); err != nil {
		return "", fmt.Errorf("extract %s: %w", name, err)
	}
	if err := os.Chmod(binaryPath, 0o755); err != nil {
		return "", fmt.Errorf("chmod %s: %w", name, err)
	}
	if err := os.WriteFile(shaPath, []byte(expectedSHA), 0o644); err != nil {
		return "", fmt.Errorf("write sha marker for %s: %w", name, err)
	}
	return binaryPath, nil
}

// fetchSHA returns the hex digest from a .sha256 file, which holds either bare
// hex or sha256sum's "<hex>  <filename>".
func fetchSHA(rawURL string) (string, error) {
	resp, err := http.Get(rawURL)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s returned %s", rawURL, resp.Status)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024))
	if err != nil {
		return "", err
	}
	fields := strings.Fields(string(body))
	if len(fields) == 0 {
		return "", fmt.Errorf("%s returned empty body", rawURL)
	}
	return fields[0], nil
}

// download saves rawURL to a temporary file, returning its path once its
// SHA-256 matches expectedSHA.
func download(rawURL, expectedSHA, name string) (string, error) {
	resp, err := http.Get(rawURL)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s returned %s", rawURL, resp.Status)
	}

	tmp, err := os.CreateTemp("", name+"-*")
	if err != nil {
		return "", err
	}
	tmpPath := tmp.Name()

	hasher := sha256.New()
	_, err = io.Copy(io.MultiWriter(tmp, hasher), resp.Body)
	if closeErr := tmp.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		_ = os.Remove(tmpPath)
		return "", err
	}

	if actual := hex.EncodeToString(hasher.Sum(nil)); actual != expectedSHA {
		_ = os.Remove(tmpPath)
		return "", fmt.Errorf("sha mismatch for %s: expected %s, got %s", name, expectedSHA, actual)
	}
	return tmpPath, nil
}

func extractTarGz(archivePath, destPath, entryName string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer func() { _ = gz.Close() }()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return fmt.Errorf("entry %q not found", entryName)
		}
		if err != nil {
			return err
		}
		if filepath.Base(hdr.Name) == entryName {
			return writeBinary(destPath, tr)
		}
	}
}

func extractZip(archivePath, destPath, entryName string) error {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer func() { _ = r.Close() }()

	for _, zf := range r.File {
		if filepath.Base(zf.Name) != entryName {
			continue
		}
		zr, err := zf.Open()
		if err != nil {
			return err
		}
		defer func() { _ = zr.Close() }()
		return writeBinary(destPath, zr)
	}
	return fmt.Errorf("entry %q not found", entryName)
}

func writeBinary(destPath string, r io.Reader) error {
	out, err := os.OpenFile(destPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, io.LimitReader(r, maxBinarySize)); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}
