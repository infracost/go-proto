package plugintest

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const pluginName = "infracost-parser-test"

type release struct {
	archive   []byte
	sha       string
	downloads atomic.Int32
}

// serveRelease serves r at the given version path and points Install at it.
func serveRelease(t *testing.T, version string, r *release) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("test archives are tar.gz only")
	}
	prefix := "/" + pluginName + "/" + runtime.GOOS + "/" + runtime.GOARCH + "/" + version + "/data.tar.gz"
	mux := http.NewServeMux()
	mux.HandleFunc(prefix+".sha256", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(r.sha + "  data.tar.gz\n"))
	})
	mux.HandleFunc(prefix, func(w http.ResponseWriter, _ *http.Request) {
		r.downloads.Add(1)
		_, _ = w.Write(r.archive)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	t.Setenv(baseURLEnv, srv.URL+"/")
}

func newRelease(t *testing.T, entries map[string]string) *release {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	for name, body := range entries {
		require.NoError(t, tw.WriteHeader(&tar.Header{Name: name, Mode: 0o755, Size: int64(len(body))}))
		_, err := tw.Write([]byte(body))
		require.NoError(t, err)
	}
	require.NoError(t, tw.Close())
	require.NoError(t, gz.Close())
	sum := sha256.Sum256(buf.Bytes())
	return &release{archive: buf.Bytes(), sha: hex.EncodeToString(sum[:])}
}

func TestInstall(t *testing.T) {
	for _, tc := range []struct{ version, path string }{
		{"latest", "latest"},
		{"1.2.3", "v1.2.3"},
		{"v1.2.3", "v1.2.3"},
	} {
		t.Run(tc.version, func(t *testing.T) {
			r := newRelease(t, map[string]string{"README": "docs", "bin/" + pluginName: "binary"})
			serveRelease(t, tc.path, r)
			dir := filepath.Join(t.TempDir(), "plugins")

			path, err := Install(dir, pluginName, tc.version)
			require.NoError(t, err)
			assert.Equal(t, filepath.Join(dir, pluginName), path)

			got, err := os.ReadFile(path)
			require.NoError(t, err)
			assert.Equal(t, "binary", string(got))
			info, err := os.Stat(path)
			require.NoError(t, err)
			assert.Equal(t, os.FileMode(0o755), info.Mode().Perm())
		})
	}
}

func TestInstallSkipsDownloadAtSameSHA(t *testing.T) {
	r := newRelease(t, map[string]string{pluginName: "binary"})
	serveRelease(t, "latest", r)
	dir := t.TempDir()

	_, err := Install(dir, pluginName, "latest")
	require.NoError(t, err)
	_, err = Install(dir, pluginName, "latest")
	require.NoError(t, err)
	assert.Equal(t, int32(1), r.downloads.Load())

	require.NoError(t, os.Remove(filepath.Join(dir, pluginName)))
	_, err = Install(dir, pluginName, "latest")
	require.NoError(t, err)
	assert.Equal(t, int32(2), r.downloads.Load(), "a missing binary is downloaded again")
}

func TestInstallErrors(t *testing.T) {
	t.Run("sha mismatch", func(t *testing.T) {
		r := newRelease(t, map[string]string{pluginName: "binary"})
		r.sha = "deadbeef"
		serveRelease(t, "latest", r)
		_, err := Install(t.TempDir(), pluginName, "latest")
		require.ErrorContains(t, err, "sha mismatch")
	})

	t.Run("binary missing from archive", func(t *testing.T) {
		serveRelease(t, "latest", newRelease(t, map[string]string{"other": "binary"}))
		_, err := Install(t.TempDir(), pluginName, "latest")
		require.ErrorContains(t, err, "not found")
	})

	t.Run("unknown version", func(t *testing.T) {
		serveRelease(t, "latest", newRelease(t, map[string]string{pluginName: "binary"}))
		_, err := Install(t.TempDir(), pluginName, "9.9.9")
		require.ErrorContains(t, err, "404")
	})
}
