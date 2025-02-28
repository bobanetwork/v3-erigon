/*
   Copyright 2021 Erigon contributors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package dbg

import (
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strconv"
	"sync"
	"time"

	"github.com/erigontech/erigon-lib/log/v3"

	libcommon "github.com/erigontech/erigon-lib/common"
	"github.com/erigontech/erigon-lib/mmap"
)

var (
	// force skipping of any non-Erigon2 .torrent files
	DownloaderOnlyBlocks = EnvBool("DOWNLOADER_ONLY_BLOCKS", false)
	saveHeapProfile      = EnvBool("SAVE_HEAP_PROFILE", false)
	heapProfileFilePath  = EnvString("HEAP_PROFILE_FILE_PATH", "")
)

var StagesOnlyBlocks = EnvBool("STAGES_ONLY_BLOCKS", false)

var doMemstat = true

func init() {
	_, ok := os.LookupEnv("NO_MEMSTAT")
	if ok {
		doMemstat = false
	}
}

func DoMemStat() bool { return doMemstat }
func ReadMemStats(m *runtime.MemStats) {
	if doMemstat {
		runtime.ReadMemStats(m)
	}
}

var (
	writeMap     bool
	writeMapOnce sync.Once
)

func WriteMap() bool {
	writeMapOnce.Do(func() {
		v, _ := os.LookupEnv("WRITE_MAP")
		if v == "true" {
			writeMap = true
			log.Info("[Experiment]", "WRITE_MAP", writeMap)
		}
	})
	return writeMap
}

var (
	dirtySace     uint64
	dirtySaceOnce sync.Once
)

func DirtySpace() uint64 {
	dirtySaceOnce.Do(func() {
		v, _ := os.LookupEnv("MDBX_DIRTY_SPACE_MB")
		if v != "" {
			i, err := strconv.Atoi(v)
			if err != nil {
				panic(err)
			}
			dirtySace = uint64(i * 1024 * 1024)
			log.Info("[Experiment]", "MDBX_DIRTY_SPACE_MB", dirtySace)
		}
	})
	return dirtySace
}

var (
	noSync     bool
	noSyncOnce sync.Once
)

func NoSync() bool {
	noSyncOnce.Do(func() {
		v, _ := os.LookupEnv("NO_SYNC")
		if v == "true" {
			noSync = true
			log.Info("[Experiment]", "NO_SYNC", noSync)
		}
	})
	return noSync
}

var (
	mergeTr     int
	mergeTrOnce sync.Once
)

func MergeTr() int {
	mergeTrOnce.Do(func() {
		v, _ := os.LookupEnv("MERGE_THRESHOLD")
		if v != "" {
			i, err := strconv.Atoi(v)
			if err != nil {
				panic(err)
			}
			if i < 0 || i > 4 {
				panic(i)
			}
			mergeTr = i
			log.Info("[Experiment]", "MERGE_THRESHOLD", mergeTr)
		}
	})
	return mergeTr
}

var (
	mdbxReadahead     bool
	mdbxReadaheadOnce sync.Once
)

func MdbxReadAhead() bool {
	mdbxReadaheadOnce.Do(func() {
		v, _ := os.LookupEnv("MDBX_READAHEAD")
		if v == "true" {
			mdbxReadahead = true
			log.Info("[Experiment]", "MDBX_READAHEAD", mdbxReadahead)
		}
	})
	return mdbxReadahead
}

var (
	discardHistory     bool
	discardHistoryOnce sync.Once
)

func DiscardHistory() bool {
	discardHistoryOnce.Do(func() {
		v, _ := os.LookupEnv("DISCARD_HISTORY")
		if v == "true" {
			discardHistory = true
			log.Info("[Experiment]", "DISCARD_HISTORY", discardHistory)
		}
	})
	return discardHistory
}

var (
	bigRoTx    uint
	getBigRoTx sync.Once
)

// DEBUG_BIG_RO_TX_KB - print logs with info about large read-only transactions
// DEBUG_BIG_RW_TX_KB - print logs with info about large read-write transactions
// DEBUG_SLOW_COMMIT_MS - print logs with commit timing details if commit is slower than this threshold
func BigRoTxKb() uint {
	getBigRoTx.Do(func() {
		v, _ := os.LookupEnv("DEBUG_BIG_RO_TX_KB")
		if v != "" {
			i, err := strconv.Atoi(v)
			if err != nil {
				panic(err)
			}
			bigRoTx = uint(i)
			log.Info("[Experiment]", "DEBUG_BIG_RO_TX_KB", bigRoTx)
		}
	})
	return bigRoTx
}

var (
	bigRwTx    uint
	getBigRwTx sync.Once
)

func BigRwTxKb() uint {
	getBigRwTx.Do(func() {
		v, _ := os.LookupEnv("DEBUG_BIG_RW_TX_KB")
		if v != "" {
			i, err := strconv.Atoi(v)
			if err != nil {
				panic(err)
			}
			bigRwTx = uint(i)
			log.Info("[Experiment]", "DEBUG_BIG_RW_TX_KB", bigRwTx)
		}
	})
	return bigRwTx
}

var (
	slowCommit     time.Duration
	slowCommitOnce sync.Once
)

func SlowCommit() time.Duration {
	slowCommitOnce.Do(func() {
		v, _ := os.LookupEnv("SLOW_COMMIT")
		if v != "" {
			var err error
			slowCommit, err = time.ParseDuration(v)
			if err != nil {
				panic(err)
			}
			log.Info("[Experiment]", "SLOW_COMMIT", slowCommit.String())
		}
	})
	return slowCommit
}

var (
	slowTx     time.Duration
	slowTxOnce sync.Once
)

func SlowTx() time.Duration {
	slowTxOnce.Do(func() {
		v, _ := os.LookupEnv("SLOW_TX")
		if v != "" {
			var err error
			slowTx, err = time.ParseDuration(v)
			if err != nil {
				panic(err)
			}
			log.Info("[Experiment]", "SLOW_TX", slowTx.String())
		}
	})
	return slowTx
}

var (
	stopBeforeStage     string
	stopBeforeStageFlag sync.Once
	stopAfterStage      string
	stopAfterStageFlag  sync.Once
)

func StopBeforeStage() string {
	f := func() {
		v, _ := os.LookupEnv("STOP_BEFORE_STAGE") // see names in eth/stagedsync/stages/stages.go
		if v != "" {
			stopBeforeStage = v
			log.Info("[Experiment]", "STOP_BEFORE_STAGE", stopBeforeStage)
		}
	}
	stopBeforeStageFlag.Do(f)
	return stopBeforeStage
}

// TODO(allada) We should possibly consider removing `STOP_BEFORE_STAGE`, as `STOP_AFTER_STAGE` can
// perform all same the functionality, but due to reverse compatibility reasons we are going to
// leave it.
func StopAfterStage() string {
	f := func() {
		v, _ := os.LookupEnv("STOP_AFTER_STAGE") // see names in eth/stagedsync/stages/stages.go
		if v != "" {
			stopAfterStage = v
			log.Info("[Experiment]", "STOP_AFTER_STAGE", stopAfterStage)
		}
	}
	stopAfterStageFlag.Do(f)
	return stopAfterStage
}

var (
	stopAfterReconst     bool
	stopAfterReconstOnce sync.Once
)

func StopAfterReconst() bool {
	stopAfterReconstOnce.Do(func() {
		v, _ := os.LookupEnv("STOP_AFTER_RECONSTITUTE")
		if v == "true" {
			stopAfterReconst = true
			log.Info("[Experiment]", "STOP_AFTER_RECONSTITUTE", stopAfterReconst)
		}
	})
	return stopAfterReconst
}

var (
	snapshotVersion     uint8
	snapshotVersionOnce sync.Once
)

func SnapshotVersion() uint8 {
	snapshotVersionOnce.Do(func() {
		v, _ := os.LookupEnv("SNAPSHOT_VERSION")
		if i, _ := strconv.ParseUint(v, 10, 8); i > 0 {
			snapshotVersion = uint8(i)
			log.Info("[Experiment]", "SNAPSHOT_VERSION", snapshotVersion)
		}
	})
	return snapshotVersion
}

var (
	logHashMismatchReason     bool
	logHashMismatchReasonOnce sync.Once
)

func LogHashMismatchReason() bool {
	logHashMismatchReasonOnce.Do(func() {
		v, _ := os.LookupEnv("LOG_HASH_MISMATCH_REASON")
		if v == "true" {
			logHashMismatchReason = true
			log.Info("[Experiment]", "LOG_HASH_MISMATCH_REASON", logHashMismatchReason)
		}
	})
	return logHashMismatchReason
}

type saveHeapOptions struct {
	memStats *runtime.MemStats
	logger   *log.Logger
}

type SaveHeapOption func(options *saveHeapOptions)

func SaveHeapWithMemStats(memStats *runtime.MemStats) SaveHeapOption {
	return func(options *saveHeapOptions) {
		options.memStats = memStats
	}
}

func SaveHeapWithLogger(logger *log.Logger) SaveHeapOption {
	return func(options *saveHeapOptions) {
		options.logger = logger
	}
}

func SaveHeapProfileNearOOM(opts ...SaveHeapOption) {
	if !saveHeapProfile {
		return
	}

	var options saveHeapOptions
	for _, opt := range opts {
		opt(&options)
	}

	var logger log.Logger
	if options.logger != nil {
		logger = *options.logger
	}

	var memStats runtime.MemStats
	if options.memStats != nil {
		memStats = *options.memStats
	} else {
		ReadMemStats(&memStats)
	}

	totalMemory := mmap.TotalMemory()
	if logger != nil {
		logger.Info(
			"[Experiment] heap profile threshold check",
			"alloc", libcommon.ByteCount(memStats.Alloc),
			"total", libcommon.ByteCount(totalMemory),
		)
	}
	if memStats.Alloc < (totalMemory/100)*45 {
		return
	}

	// above 45%
	var filePath string
	if heapProfileFilePath == "" {
		filePath = filepath.Join(os.TempDir(), "erigon-mem.prof")
	} else {
		filePath = heapProfileFilePath
	}
	if logger != nil {
		logger.Info("[Experiment] saving heap profile as near OOM", "filePath", filePath)
	}

	f, err := os.Create(filePath)
	if err != nil && logger != nil {
		logger.Warn("[Experiment] could not create heap profile file", "err", err)
	}

	defer func() {
		err := f.Close()
		if err != nil && logger != nil {
			logger.Warn("[Experiment] could not close heap profile file", "err", err)
		}
	}()

	runtime.GC()
	err = pprof.WriteHeapProfile(f)
	if err != nil && logger != nil {
		logger.Warn("[Experiment] could not write heap profile file", "err", err)
	}
}
