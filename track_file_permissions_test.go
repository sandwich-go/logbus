package logbus

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestOpenTrackFileDirectoryPermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix directory permissions")
	}
	// Keep the restrictive umask isolated from other tests and background workers.
	const childEnv = "LOGBUS_TEST_TRACK_PERMISSIONS_CHILD"
	if os.Getenv(childEnv) != "1" {
		cmd := exec.Command("sh", "-c", `umask 077; exec "$@"`, "sh", os.Args[0], "-test.run=^TestOpenTrackFileDirectoryPermissions$", "-test.count=1")
		cmd.Env = append(os.Environ(), childEnv+"=1")
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("permission tests with umask 077: %v\n%s", err, output)
		}
		return
	}

	for _, existing := range []bool{false, true} {
		name := "new_directories"
		if existing {
			name = "existing_restrictive_directories"
		}
		t.Run(name, func(t *testing.T) {
			parent := t.TempDir()
			base := filepath.Join(parent, "logs", "tracklogs")
			channelDir := filepath.Join(base, "tga")
			dateDir := filepath.Join(channelDir, "20261015")
			if existing {
				if err := os.MkdirAll(dateDir, 0o700); err != nil {
					t.Fatal(err)
				}
			}
			for _, date := range []string{"20261015", "20261016"} {
				for _, pod := range []string{"pod-a", "pod-b"} {
					f, _, _, err := openTrackFile(base, "tga", pod, date+"00", date)
					if err != nil {
						t.Fatal(err)
					}
					if err := f.Close(); err != nil {
						t.Fatal(err)
					}
				}
				for _, dir := range []string{filepath.Dir(base), base, channelDir, filepath.Join(channelDir, date)} {
					info, err := os.Stat(dir)
					if err != nil {
						t.Fatal(err)
					}
					want := os.FileMode(0o777)
					if existing && dir != filepath.Join(channelDir, "20261016") {
						want = 0o700
					}
					if got := info.Mode().Perm(); got != want {
						t.Errorf("%s permissions = %04o, want %04o", dir, got, want)
					}
				}
			}
			info, err := os.Stat(parent)
			if err != nil {
				t.Fatal(err)
			}
			if got := info.Mode().Perm(); got != 0o700 {
				t.Errorf("parent permissions = %04o, want unchanged 0700", got)
			}
		})
	}
}
