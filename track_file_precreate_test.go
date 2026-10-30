package logbus

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNextTrackDirPreparationDelay(t *testing.T) {
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name  string
		now   time.Time
		retry bool
		want  time.Duration
	}{
		{"midday", time.Date(2026, 10, 27, 12, 0, 0, 0, location), false, 11*time.Hour + 55*time.Minute},
		{"before_schedule", time.Date(2026, 10, 27, 23, 54, 0, 0, location), false, time.Minute},
		{"at_schedule", time.Date(2026, 10, 27, 23, 55, 0, 0, location), false, 24 * time.Hour},
		{"after_schedule", time.Date(2026, 10, 27, 23, 59, 0, 0, location), false, 23*time.Hour + 56*time.Minute},
		{"year_end", time.Date(2026, 12, 31, 23, 55, 0, 0, location), false, 24 * time.Hour},
		{"dst_spring", time.Date(2026, 3, 7, 23, 55, 0, 0, location), false, 23 * time.Hour},
		{"dst_fall", time.Date(2026, 10, 31, 23, 55, 0, 0, location), false, 25 * time.Hour},
		{"failure_before_schedule", time.Date(2026, 10, 27, 12, 0, 0, 0, location), true, time.Minute},
		{"failure_after_schedule", time.Date(2026, 10, 27, 23, 59, 0, 0, location), true, time.Minute},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := nextTrackDirPreparationDelay(tc.now, tc.retry); got != tc.want {
				t.Fatalf("delay = %s, want %s", got, tc.want)
			}
		})
	}
}

func TestTrackWorkerPrecreatesNextDayWithoutTraffic(t *testing.T) {
	base := t.TempDir()
	before := time.Now()
	w := newTrackWorker("tga", base, "pod-a", HourlyRotation, 8, 0, nil)
	defer w.stop()
	if err := w.doSync(); err != nil {
		t.Fatal(err)
	}
	// If this test crosses midnight, startup may have used either local date.
	for _, now := range []time.Time{before, time.Now()} {
		dir := filepath.Join(base, "tga", dateFromSlot(now.AddDate(0, 0, 1)))
		entries, err := os.ReadDir(dir)
		if err == nil {
			if len(entries) != 0 {
				t.Fatal("precreation must not create log files")
			}
			return
		}
		if !os.IsNotExist(err) {
			t.Fatal(err)
		}
	}
	t.Fatal("next-day directory was not created without traffic or keepalive")
}

func TestTrackWorkerPrepareNextDayCalendar(t *testing.T) {
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name string
		now  time.Time
		want string
	}{
		{"month_end", time.Date(2026, 10, 31, 23, 59, 0, 0, location), "20261101"},
		{"year_end", time.Date(2026, 12, 31, 23, 59, 0, 0, location), "20270101"},
		{"dst_spring", time.Date(2026, 3, 8, 0, 30, 0, 0, location), "20260309"},
		{"dst_fall", time.Date(2026, 11, 1, 0, 30, 0, 0, location), "20261102"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w := &trackWorker{baseDir: t.TempDir(), channel: "tga"}
			for _, now := range []time.Time{tc.now, tc.now, tc.now.AddDate(0, 0, 1)} {
				if err := w.prepareNextDayDir(now); err != nil {
					t.Fatal(err)
				}
			}
			for _, date := range []string{tc.want, dateFromSlot(tc.now.AddDate(0, 0, 2))} {
				entries, err := os.ReadDir(filepath.Join(w.baseDir, "tga", date))
				if err != nil {
					t.Fatal(err)
				}
				if len(entries) != 0 {
					t.Fatal("precreation must only create directories")
				}
			}
		})
	}
}

func TestTrackWorkerPrepareNextDayRetriesFailure(t *testing.T) {
	w := &trackWorker{baseDir: t.TempDir(), channel: "tga"}
	now := time.Date(2026, 10, 27, 12, 0, 0, 0, time.Local)
	path := filepath.Join(w.baseDir, "tga")
	if err := os.WriteFile(path, []byte("blocked"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := w.prepareNextDayDir(now); err == nil {
		t.Fatal("expected failure when channel path is a file")
	}
	if w.nextDirDate != "" {
		t.Fatal("failed preparation must not be cached")
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := w.prepareNextDayDir(now); err != nil {
		t.Fatal(err)
	}
	if w.nextDirDate != "20261028" {
		t.Fatalf("prepared date = %q, want 20261028", w.nextDirDate)
	}
}

func TestTrackWorkerPrecreateFailureDoesNotDropCurrentBatch(t *testing.T) {
	base := t.TempDir()
	now := time.Now()
	channelDir := filepath.Join(base, "tga")
	if err := mkdirTrackDirs(channelDir); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(channelDir, dateFromSlot(now.AddDate(0, 0, 1))), nil, 0o600); err != nil {
		t.Fatal(err)
	}
	w := &trackWorker{baseDir: base, channel: "tga", hostName: "pod-a", rotation: HourlyRotation}
	if err := w.prepareNextDayDir(now); err == nil {
		t.Fatal("expected preparation failure")
	}
	if err := w.writeBatch([]byte("current batch\n")); err != nil {
		t.Fatal(err)
	}
	path := w.file.Name()
	if err := w.file.Close(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "current batch\n" {
		t.Fatalf("current batch = %q", data)
	}
}
