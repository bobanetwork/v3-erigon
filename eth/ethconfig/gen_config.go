// Code generated by github.com/fjl/gencodec. DO NOT EDIT.

package ethconfig

import (
	"math/big"
	"time"

	"github.com/c2h5oh/datasize"
	"github.com/erigontech/erigon-lib/chain"
	"github.com/erigontech/erigon/consensus/ethash/ethashcfg"
	"github.com/erigontech/erigon/core/types"
	"github.com/erigontech/erigon/eth/gasprice/gaspricecfg"

	"github.com/erigontech/erigon/ethdb/prune"
	"github.com/erigontech/erigon/params"
	"github.com/erigontech/erigon-lib/common"
	"github.com/erigontech/erigon-lib/common/datadir"
	"github.com/erigontech/erigon-lib/downloader/downloadercfg"
	"github.com/erigontech/erigon-lib/txpool/txpoolcfg"
	"github.com/erigontech/erigon/cl/beacon/beacon_router_configuration"
	"github.com/erigontech/erigon/cl/clparams"
)

// MarshalTOML marshals as TOML.
func (c Config) MarshalTOML() (interface{}, error) {
	type Config struct {
		Sync                                    Sync
		Genesis                                 *types.Genesis `toml:",omitempty"`
		NetworkID                               uint64
		EthDiscoveryURLs                        []string
		Prune                                   prune.Mode
		BatchSize                               datasize.ByteSize
		ImportMode                              bool
		BadBlockHash                            common.Hash
		Snapshot                                BlocksFreezing
		Downloader                              *downloadercfg.Cfg
		BeaconRouter                            beacon_router_configuration.RouterConfiguration
		CaplinConfig                            clparams.CaplinConfig
		Dirs                                    datadir.Dirs
		ExternalSnapshotDownloaderAddr          string
		Whitelist                               map[uint64]common.Hash `toml:"-"`
		Miner                                   params.MiningConfig
		Ethash                                  ethashcfg.Config
		Clique                                  params.ConsensusSnapshotConfig
		Aura                                    chain.AuRaConfig
		DeprecatedTxPool                        DeprecatedTxPoolConfig
		TxPool                                  txpoolcfg.Config
		GPO                                     gaspricecfg.Config
		RPCGasCap                               uint64  `toml:",omitempty"`
		RPCTxFeeCap                             float64 `toml:",omitempty"`
		StateStream                             bool
		HistoryV3                               bool
		HeimdallURL                             string
		WithoutHeimdall                         bool
		WithHeimdallMilestones                  bool
		WithHeimdallWaypointRecording           bool
		PolygonSync                             bool
		Ethstats                                string
		InternalCL                              bool
		LightClientDiscoveryAddr                string
		LightClientDiscoveryPort                uint64
		LightClientDiscoveryTCPPort             uint64
		SentinelAddr                            string
		SentinelPort                            uint64
		ForcePartialCommit                      bool
		OverrideCancunTime                      *big.Int `toml:",omitempty"`
		OverrideShanghaiTime                    *big.Int `toml:",omitempty"`
		OverridePragueTime                      *big.Int `toml:",omitempty"`
		OverrideOptimismCanyonTime              *big.Int `toml:",omitempty"`
		OverrideOptimismEcotoneTime             *big.Int `toml:",omitempty"`
		OverrideOptimismFjordTime               *big.Int `toml:",omitempty"`
		OverrideOptimismGraniteTime             *big.Int `toml:",omitempty"`
		OverrideOptimismHoloceneTime            *big.Int `toml:",omitempty"`
		SilkwormExecution                       bool
		SilkwormRpcDaemon                       bool
		SilkwormSentry                          bool
		SilkwormVerbosity                       string
		SilkwormNumContexts                     uint32
		SilkwormRpcLogEnabled                   bool
		SilkwormRpcLogDirPath                   string
		SilkwormRpcLogMaxFileSize               uint16
		SilkwormRpcLogMaxFiles                  uint16
		SilkwormRpcLogDumpResponse              bool
		SilkwormRpcNumWorkers                   uint32
		SilkwormRpcJsonCompatibility            bool
		DisableTxPoolGossip                     bool
		RollupSequencerHTTP                     string
		RollupHistoricalRPC                     string
		RollupHistoricalRPCTimeout              time.Duration
		RollupHaltOnIncompatibleProtocolVersion string
	}
	var enc Config
	enc.Sync = c.Sync
	enc.Genesis = c.Genesis
	enc.NetworkID = c.NetworkID
	enc.EthDiscoveryURLs = c.EthDiscoveryURLs
	enc.Prune = c.Prune
	enc.BatchSize = c.BatchSize
	enc.ImportMode = c.ImportMode
	enc.BadBlockHash = c.BadBlockHash
	enc.Snapshot = c.Snapshot
	enc.Downloader = c.Downloader
	enc.BeaconRouter = c.BeaconRouter
	enc.CaplinConfig = c.CaplinConfig
	enc.Dirs = c.Dirs
	enc.ExternalSnapshotDownloaderAddr = c.ExternalSnapshotDownloaderAddr
	enc.Whitelist = c.Whitelist
	enc.Miner = c.Miner
	enc.Ethash = c.Ethash
	enc.Clique = c.Clique
	enc.Aura = c.Aura
	enc.DeprecatedTxPool = c.DeprecatedTxPool
	enc.TxPool = c.TxPool
	enc.GPO = c.GPO
	enc.RPCGasCap = c.RPCGasCap
	enc.RPCTxFeeCap = c.RPCTxFeeCap
	enc.StateStream = c.StateStream
	enc.HistoryV3 = c.HistoryV3
	enc.HeimdallURL = c.HeimdallURL
	enc.WithoutHeimdall = c.WithoutHeimdall
	enc.WithHeimdallMilestones = c.WithHeimdallMilestones
	enc.WithHeimdallWaypointRecording = c.WithHeimdallWaypointRecording
	enc.PolygonSync = c.PolygonSync
	enc.Ethstats = c.Ethstats
	enc.InternalCL = c.InternalCL
	enc.LightClientDiscoveryAddr = c.LightClientDiscoveryAddr
	enc.LightClientDiscoveryPort = c.LightClientDiscoveryPort
	enc.LightClientDiscoveryTCPPort = c.LightClientDiscoveryTCPPort
	enc.SentinelAddr = c.SentinelAddr
	enc.SentinelPort = c.SentinelPort
	enc.ForcePartialCommit = c.ForcePartialCommit
	enc.OverrideCancunTime = c.OverrideCancunTime
	enc.OverrideShanghaiTime = c.OverrideShanghaiTime
	enc.OverridePragueTime = c.OverridePragueTime
	enc.OverrideOptimismCanyonTime = c.OverrideOptimismCanyonTime
	enc.OverrideOptimismEcotoneTime = c.OverrideOptimismEcotoneTime
	enc.OverrideOptimismFjordTime = c.OverrideOptimismFjordTime
	enc.OverrideOptimismGraniteTime = c.OverrideOptimismGraniteTime
	enc.OverrideOptimismHoloceneTime = c.OverrideOptimismHoloceneTime
	enc.SilkwormExecution = c.SilkwormExecution
	enc.SilkwormRpcDaemon = c.SilkwormRpcDaemon
	enc.SilkwormSentry = c.SilkwormSentry
	enc.SilkwormVerbosity = c.SilkwormVerbosity
	enc.SilkwormNumContexts = c.SilkwormNumContexts
	enc.SilkwormRpcLogEnabled = c.SilkwormRpcLogEnabled
	enc.SilkwormRpcLogDirPath = c.SilkwormRpcLogDirPath
	enc.SilkwormRpcLogMaxFileSize = c.SilkwormRpcLogMaxFileSize
	enc.SilkwormRpcLogMaxFiles = c.SilkwormRpcLogMaxFiles
	enc.SilkwormRpcLogDumpResponse = c.SilkwormRpcLogDumpResponse
	enc.SilkwormRpcNumWorkers = c.SilkwormRpcNumWorkers
	enc.SilkwormRpcJsonCompatibility = c.SilkwormRpcJsonCompatibility
	enc.DisableTxPoolGossip = c.DisableTxPoolGossip
	enc.RollupSequencerHTTP = c.RollupSequencerHTTP
	enc.RollupHistoricalRPC = c.RollupHistoricalRPC
	enc.RollupHistoricalRPCTimeout = c.RollupHistoricalRPCTimeout
	enc.RollupHaltOnIncompatibleProtocolVersion = c.RollupHaltOnIncompatibleProtocolVersion
	return &enc, nil
}

// UnmarshalTOML unmarshals from TOML.
func (c *Config) UnmarshalTOML(unmarshal func(interface{}) error) error {
	type Config struct {
		Sync                                    *Sync
		Genesis                                 *types.Genesis `toml:",omitempty"`
		NetworkID                               *uint64
		EthDiscoveryURLs                        []string
		Prune                                   *prune.Mode
		BatchSize                               *datasize.ByteSize
		ImportMode                              *bool
		BadBlockHash                            *common.Hash
		Snapshot                                *BlocksFreezing
		Downloader                              *downloadercfg.Cfg
		BeaconRouter                            *beacon_router_configuration.RouterConfiguration
		CaplinConfig                            *clparams.CaplinConfig
		Dirs                                    *datadir.Dirs
		ExternalSnapshotDownloaderAddr          *string
		Whitelist                               map[uint64]common.Hash `toml:"-"`
		Miner                                   *params.MiningConfig
		Ethash                                  *ethashcfg.Config
		Clique                                  *params.ConsensusSnapshotConfig
		Aura                                    *chain.AuRaConfig
		DeprecatedTxPool                        *DeprecatedTxPoolConfig
		TxPool                                  *txpoolcfg.Config
		GPO                                     *gaspricecfg.Config
		RPCGasCap                               *uint64  `toml:",omitempty"`
		RPCTxFeeCap                             *float64 `toml:",omitempty"`
		StateStream                             *bool
		HistoryV3                               *bool
		HeimdallURL                             *string
		WithoutHeimdall                         *bool
		WithHeimdallMilestones                  *bool
		WithHeimdallWaypointRecording           *bool
		PolygonSync                             *bool
		Ethstats                                *string
		InternalCL                              *bool
		LightClientDiscoveryAddr                *string
		LightClientDiscoveryPort                *uint64
		LightClientDiscoveryTCPPort             *uint64
		SentinelAddr                            *string
		SentinelPort                            *uint64
		ForcePartialCommit                      *bool
		OverrideCancunTime                      *big.Int `toml:",omitempty"`
		OverrideShanghaiTime                    *big.Int `toml:",omitempty"`
		OverridePragueTime                      *big.Int `toml:",omitempty"`
		OverrideOptimismCanyonTime              *big.Int `toml:",omitempty"`
		OverrideOptimismEcotoneTime             *big.Int `toml:",omitempty"`
		OverrideOptimismFjordTime               *big.Int `toml:",omitempty"`
		OverrideOptimismGraniteTime             *big.Int `toml:",omitempty"`
		OverrideOptimismHoloceneTime            *big.Int `toml:",omitempty"`
		SilkwormExecution                       *bool
		SilkwormRpcDaemon                       *bool
		SilkwormSentry                          *bool
		SilkwormVerbosity                       *string
		SilkwormNumContexts                     *uint32
		SilkwormRpcLogEnabled                   *bool
		SilkwormRpcLogDirPath                   *string
		SilkwormRpcLogMaxFileSize               *uint16
		SilkwormRpcLogMaxFiles                  *uint16
		SilkwormRpcLogDumpResponse              *bool
		SilkwormRpcNumWorkers                   *uint32
		SilkwormRpcJsonCompatibility            *bool
		DisableTxPoolGossip                     *bool
		RollupSequencerHTTP                     *string
		RollupHistoricalRPC                     *string
		RollupHistoricalRPCTimeout              *time.Duration
		RollupHaltOnIncompatibleProtocolVersion *string
	}
	var dec Config
	if err := unmarshal(&dec); err != nil {
		return err
	}
	if dec.Sync != nil {
		c.Sync = *dec.Sync
	}
	if dec.Genesis != nil {
		c.Genesis = dec.Genesis
	}
	if dec.NetworkID != nil {
		c.NetworkID = *dec.NetworkID
	}
	if dec.EthDiscoveryURLs != nil {
		c.EthDiscoveryURLs = dec.EthDiscoveryURLs
	}
	if dec.Prune != nil {
		c.Prune = *dec.Prune
	}
	if dec.BatchSize != nil {
		c.BatchSize = *dec.BatchSize
	}
	if dec.ImportMode != nil {
		c.ImportMode = *dec.ImportMode
	}
	if dec.BadBlockHash != nil {
		c.BadBlockHash = *dec.BadBlockHash
	}
	if dec.Snapshot != nil {
		c.Snapshot = *dec.Snapshot
	}
	if dec.Downloader != nil {
		c.Downloader = dec.Downloader
	}
	if dec.BeaconRouter != nil {
		c.BeaconRouter = *dec.BeaconRouter
	}
	if dec.CaplinConfig != nil {
		c.CaplinConfig = *dec.CaplinConfig
	}
	if dec.Dirs != nil {
		c.Dirs = *dec.Dirs
	}
	if dec.ExternalSnapshotDownloaderAddr != nil {
		c.ExternalSnapshotDownloaderAddr = *dec.ExternalSnapshotDownloaderAddr
	}
	if dec.Whitelist != nil {
		c.Whitelist = dec.Whitelist
	}
	if dec.Miner != nil {
		c.Miner = *dec.Miner
	}
	if dec.Ethash != nil {
		c.Ethash = *dec.Ethash
	}
	if dec.Clique != nil {
		c.Clique = *dec.Clique
	}
	if dec.Aura != nil {
		c.Aura = *dec.Aura
	}
	if dec.DeprecatedTxPool != nil {
		c.DeprecatedTxPool = *dec.DeprecatedTxPool
	}
	if dec.TxPool != nil {
		c.TxPool = *dec.TxPool
	}
	if dec.GPO != nil {
		c.GPO = *dec.GPO
	}
	if dec.RPCGasCap != nil {
		c.RPCGasCap = *dec.RPCGasCap
	}
	if dec.RPCTxFeeCap != nil {
		c.RPCTxFeeCap = *dec.RPCTxFeeCap
	}
	if dec.StateStream != nil {
		c.StateStream = *dec.StateStream
	}
	if dec.HistoryV3 != nil {
		c.HistoryV3 = *dec.HistoryV3
	}
	if dec.HeimdallURL != nil {
		c.HeimdallURL = *dec.HeimdallURL
	}
	if dec.WithoutHeimdall != nil {
		c.WithoutHeimdall = *dec.WithoutHeimdall
	}
	if dec.WithHeimdallMilestones != nil {
		c.WithHeimdallMilestones = *dec.WithHeimdallMilestones
	}
	if dec.WithHeimdallWaypointRecording != nil {
		c.WithHeimdallWaypointRecording = *dec.WithHeimdallWaypointRecording
	}
	if dec.PolygonSync != nil {
		c.PolygonSync = *dec.PolygonSync
	}
	if dec.Ethstats != nil {
		c.Ethstats = *dec.Ethstats
	}
	if dec.InternalCL != nil {
		c.InternalCL = *dec.InternalCL
	}
	if dec.LightClientDiscoveryAddr != nil {
		c.LightClientDiscoveryAddr = *dec.LightClientDiscoveryAddr
	}
	if dec.LightClientDiscoveryPort != nil {
		c.LightClientDiscoveryPort = *dec.LightClientDiscoveryPort
	}
	if dec.LightClientDiscoveryTCPPort != nil {
		c.LightClientDiscoveryTCPPort = *dec.LightClientDiscoveryTCPPort
	}
	if dec.SentinelAddr != nil {
		c.SentinelAddr = *dec.SentinelAddr
	}
	if dec.SentinelPort != nil {
		c.SentinelPort = *dec.SentinelPort
	}
	if dec.ForcePartialCommit != nil {
		c.ForcePartialCommit = *dec.ForcePartialCommit
	}
	if dec.OverrideCancunTime != nil {
		c.OverrideCancunTime = dec.OverrideCancunTime
	}
	if dec.OverrideShanghaiTime != nil {
		c.OverrideShanghaiTime = dec.OverrideShanghaiTime
	}
	if dec.OverridePragueTime != nil {
		c.OverridePragueTime = dec.OverridePragueTime
	}
	if dec.OverrideOptimismCanyonTime != nil {
		c.OverrideOptimismCanyonTime = dec.OverrideOptimismCanyonTime
	}
	if dec.OverrideOptimismEcotoneTime != nil {
		c.OverrideOptimismEcotoneTime = dec.OverrideOptimismEcotoneTime
	}
	if dec.OverrideOptimismFjordTime != nil {
		c.OverrideOptimismFjordTime = dec.OverrideOptimismFjordTime
	}
	if dec.OverrideOptimismGraniteTime != nil {
		c.OverrideOptimismGraniteTime = dec.OverrideOptimismGraniteTime
	}
	if dec.OverrideOptimismHoloceneTime != nil {
		c.OverrideOptimismHoloceneTime = dec.OverrideOptimismHoloceneTime
	}
	if dec.SilkwormExecution != nil {
		c.SilkwormExecution = *dec.SilkwormExecution
	}
	if dec.SilkwormRpcDaemon != nil {
		c.SilkwormRpcDaemon = *dec.SilkwormRpcDaemon
	}
	if dec.SilkwormSentry != nil {
		c.SilkwormSentry = *dec.SilkwormSentry
	}
	if dec.SilkwormVerbosity != nil {
		c.SilkwormVerbosity = *dec.SilkwormVerbosity
	}
	if dec.SilkwormNumContexts != nil {
		c.SilkwormNumContexts = *dec.SilkwormNumContexts
	}
	if dec.SilkwormRpcLogEnabled != nil {
		c.SilkwormRpcLogEnabled = *dec.SilkwormRpcLogEnabled
	}
	if dec.SilkwormRpcLogDirPath != nil {
		c.SilkwormRpcLogDirPath = *dec.SilkwormRpcLogDirPath
	}
	if dec.SilkwormRpcLogMaxFileSize != nil {
		c.SilkwormRpcLogMaxFileSize = *dec.SilkwormRpcLogMaxFileSize
	}
	if dec.SilkwormRpcLogMaxFiles != nil {
		c.SilkwormRpcLogMaxFiles = *dec.SilkwormRpcLogMaxFiles
	}
	if dec.SilkwormRpcLogDumpResponse != nil {
		c.SilkwormRpcLogDumpResponse = *dec.SilkwormRpcLogDumpResponse
	}
	if dec.SilkwormRpcNumWorkers != nil {
		c.SilkwormRpcNumWorkers = *dec.SilkwormRpcNumWorkers
	}
	if dec.SilkwormRpcJsonCompatibility != nil {
		c.SilkwormRpcJsonCompatibility = *dec.SilkwormRpcJsonCompatibility
	}
	if dec.DisableTxPoolGossip != nil {
		c.DisableTxPoolGossip = *dec.DisableTxPoolGossip
	}
	if dec.RollupSequencerHTTP != nil {
		c.RollupSequencerHTTP = *dec.RollupSequencerHTTP
	}
	if dec.RollupHistoricalRPC != nil {
		c.RollupHistoricalRPC = *dec.RollupHistoricalRPC
	}
	if dec.RollupHistoricalRPCTimeout != nil {
		c.RollupHistoricalRPCTimeout = *dec.RollupHistoricalRPCTimeout
	}
	if dec.RollupHaltOnIncompatibleProtocolVersion != nil {
		c.RollupHaltOnIncompatibleProtocolVersion = *dec.RollupHaltOnIncompatibleProtocolVersion
	}
	return nil
}
