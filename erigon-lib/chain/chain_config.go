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

package chain

import (
	"fmt"
	"math/big"
	"sort"
	"strconv"

	"github.com/ledgerwatch/erigon-lib/common"
	"github.com/ledgerwatch/erigon-lib/common/fixedgas"
)

// Boba chain config
var (
	// Mainnet
	BobaMainnetChainId = big.NewInt(288)
	// Boba Mainnet genesis gas limit
	BobaMainnetGenesisGasLimit = 11000000
	// Boba Mainnet genesis block coinbase
	BobaMainnetGenesisCoinbase = "0x0000000000000000000000000000000000000000"
	// Boba Mainnet genesis block extra data
	BobaMainnetGenesisExtraData = "000000000000000000000000000000000000000000000000000000000000000000000398232e2064f896018496b4b44b3d62751f0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	// Boba Mainnet genesis root
	BobaMainnetGenesisRoot = "0x7ec54492a4504ff1ef3491825cd55e01e5c75409e4287129170e98d4693848ce"

	// Boba Sepolia
	BobaSepoliaChainId = big.NewInt(28882)
	// Boba Sepolia genesis gas limit
	BobaSepoliaGenesisGasLimit = 11000000
	// Boba Sepolia genesis block coinbase
	BobaSepoliaGenesisCoinbase = "0x0000000000000000000000000000000000000000"
	// Boba Sepolia genesis block extra data
	BobaSepoliaGenesisExtraData = "000000000000000000000000000000000000000000000000000000000000000000000398232e2064f896018496b4b44b3d62751f0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	// Boba Sepolia genesis root
	BobaSepoliaGenesisRoot = "0x8c57d7486ebd810dc728748553b08919c81024f024651afdbd076780c48621b0"
)

// Config is the core config which determines the blockchain settings.
//
// Config is stored in the database on a per block basis. This means
// that any network, identified by its genesis block, can have its own
// set of configuration options.
type Config struct {
	ChainName string
	ChainID   *big.Int `json:"chainId"` // chainId identifies the current chain and is used for replay protection

	Consensus ConsensusName `json:"consensus,omitempty"` // aura, ethash or clique

	// *Block fields activate the corresponding hard fork at a certain block number,
	// while *Time fields do so based on the block's time stamp.
	// nil means that the hard-fork is not scheduled,
	// while 0 means that it's already activated from genesis.

	// ETH mainnet upgrades
	// See https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades
	HomesteadBlock        *big.Int `json:"homesteadBlock,omitempty"`
	DAOForkBlock          *big.Int `json:"daoForkBlock,omitempty"`
	TangerineWhistleBlock *big.Int `json:"eip150Block,omitempty"`
	SpuriousDragonBlock   *big.Int `json:"eip155Block,omitempty"`
	ByzantiumBlock        *big.Int `json:"byzantiumBlock,omitempty"`
	ConstantinopleBlock   *big.Int `json:"constantinopleBlock,omitempty"`
	PetersburgBlock       *big.Int `json:"petersburgBlock,omitempty"`
	IstanbulBlock         *big.Int `json:"istanbulBlock,omitempty"`
	MuirGlacierBlock      *big.Int `json:"muirGlacierBlock,omitempty"`
	BerlinBlock           *big.Int `json:"berlinBlock,omitempty"`
	LondonBlock           *big.Int `json:"londonBlock,omitempty"`
	ArrowGlacierBlock     *big.Int `json:"arrowGlacierBlock,omitempty"`
	GrayGlacierBlock      *big.Int `json:"grayGlacierBlock,omitempty"`

	// EIP-3675: Upgrade consensus to Proof-of-Stake (a.k.a. "Paris", "The Merge")
	TerminalTotalDifficulty       *big.Int `json:"terminalTotalDifficulty,omitempty"`       // The merge happens when terminal total difficulty is reached
	TerminalTotalDifficultyPassed bool     `json:"terminalTotalDifficultyPassed,omitempty"` // Disable PoW sync for networks that have already passed through the Merge
	MergeNetsplitBlock            *big.Int `json:"mergeNetsplitBlock,omitempty"`            // Virtual fork after The Merge to use as a network splitter; see FORK_NEXT_VALUE in EIP-3675

	// Mainnet fork scheduling switched from block numbers to timestamps after The Merge
	ShanghaiTime *big.Int `json:"shanghaiTime,omitempty"`
	CancunTime   *big.Int `json:"cancunTime,omitempty"`
	PragueTime   *big.Int `json:"pragueTime,omitempty"`

	// Optimism Forks
	BedrockBlock *big.Int `json:"bedrockBlock,omitempty"` // bedrockSwitch block (nil = no fork, 0 = already actived)
	RegolithTime *big.Int `json:"regolithTime,omitempty"` // Regolith switch time (nil = no fork, 0 = already on optimism regolith)
	CanyonTime   *big.Int `json:"canyonTime,omitempty"`   // Canyon switch time (nil = no fork, 0 = already on optimism canyon)

	// Optional EIP-4844 parameters
	MinBlobGasPrice            *uint64 `json:"minBlobGasPrice,omitempty"`
	MaxBlobGasPerBlock         *uint64 `json:"maxBlobGasPerBlock,omitempty"`
	TargetBlobGasPerBlock      *uint64 `json:"targetBlobGasPerBlock,omitempty"`
	BlobGasPriceUpdateFraction *uint64 `json:"blobGasPriceUpdateFraction,omitempty"`

	// (Optional) governance contract where EIP-1559 fees will be sent to that otherwise would be burnt since the London fork
	BurntContract map[string]common.Address `json:"burntContract,omitempty"`

	// Various consensus engines
	Ethash   *EthashConfig   `json:"ethash,omitempty"`
	Clique   *CliqueConfig   `json:"clique,omitempty"`
	Aura     *AuRaConfig     `json:"aura,omitempty"`
	Bor      *BorConfig      `json:"bor,omitempty"`
	Optimism *OptimismConfig `json:"optimism,omitempty"`
}

func (c *Config) String() string {
	engine := c.getEngine()

	return fmt.Sprintf("{ChainID: %v, Homestead: %v, DAO: %v, Tangerine Whistle: %v, Spurious Dragon: %v, Byzantium: %v, Constantinople: %v, Petersburg: %v, Istanbul: %v, Muir Glacier: %v, Berlin: %v, London: %v, Arrow Glacier: %v, Gray Glacier: %v, Terminal Total Difficulty: %v, Merge Netsplit: %v, Shanghai: %v, Cancun: %v, Prague: %v, Engine: %v}",
		c.ChainID,
		c.HomesteadBlock,
		c.DAOForkBlock,
		c.TangerineWhistleBlock,
		c.SpuriousDragonBlock,
		c.ByzantiumBlock,
		c.ConstantinopleBlock,
		c.PetersburgBlock,
		c.IstanbulBlock,
		c.MuirGlacierBlock,
		c.BerlinBlock,
		c.LondonBlock,
		c.ArrowGlacierBlock,
		c.GrayGlacierBlock,
		c.TerminalTotalDifficulty,
		c.MergeNetsplitBlock,
		c.ShanghaiTime,
		c.CancunTime,
		c.PragueTime,
		engine,
	)
}

func (c *Config) getEngine() string {
	switch {
	case c.Ethash != nil:
		return c.Ethash.String()
	case c.Clique != nil:
		return c.Clique.String()
	case c.Bor != nil:
		return c.Bor.String()
	case c.Aura != nil:
		return c.Aura.String()
	case c.Optimism != nil:
		return c.Optimism.String()
	default:
		return "unknown"
	}
}

// IsHomestead returns whether num is either equal to the homestead block or greater.
func (c *Config) IsHomestead(num uint64) bool {
	return isForked(c.HomesteadBlock, num)
}

// IsDAOFork returns whether num is either equal to the DAO fork block or greater.
func (c *Config) IsDAOFork(num uint64) bool {
	return isForked(c.DAOForkBlock, num)
}

// IsTangerineWhistle returns whether num is either equal to the Tangerine Whistle (EIP150) fork block or greater.
func (c *Config) IsTangerineWhistle(num uint64) bool {
	return isForked(c.TangerineWhistleBlock, num)
}

// IsSpuriousDragon returns whether num is either equal to the Spurious Dragon fork block or greater.
func (c *Config) IsSpuriousDragon(num uint64) bool {
	return isForked(c.SpuriousDragonBlock, num)
}

// IsByzantium returns whether num is either equal to the Byzantium fork block or greater.
func (c *Config) IsByzantium(num uint64) bool {
	return isForked(c.ByzantiumBlock, num)
}

// IsConstantinople returns whether num is either equal to the Constantinople fork block or greater.
func (c *Config) IsConstantinople(num uint64) bool {
	return isForked(c.ConstantinopleBlock, num)
}

// IsMuirGlacier returns whether num is either equal to the Muir Glacier (EIP-2384) fork block or greater.
func (c *Config) IsMuirGlacier(num uint64) bool {
	return isForked(c.MuirGlacierBlock, num)
}

// IsPetersburg returns whether num is either
// - equal to or greater than the PetersburgBlock fork block,
// - OR is nil, and Constantinople is active
func (c *Config) IsPetersburg(num uint64) bool {
	return isForked(c.PetersburgBlock, num) || c.PetersburgBlock == nil && isForked(c.ConstantinopleBlock, num)
}

// IsIstanbul returns whether num is either equal to the Istanbul fork block or greater.
func (c *Config) IsIstanbul(num uint64) bool {
	return isForked(c.IstanbulBlock, num)
}

// IsBerlin returns whether num is either equal to the Berlin fork block or greater.
func (c *Config) IsBerlin(num uint64) bool {
	return isForked(c.BerlinBlock, num)
}

// IsLondon returns whether num is either equal to the London fork block or greater.
func (c *Config) IsLondon(num uint64) bool {
	return isForked(c.LondonBlock, num)
}

// IsArrowGlacier returns whether num is either equal to the Arrow Glacier (EIP-4345) fork block or greater.
func (c *Config) IsArrowGlacier(num uint64) bool {
	return isForked(c.ArrowGlacierBlock, num)
}

// IsGrayGlacier returns whether num is either equal to the Gray Glacier (EIP-5133) fork block or greater.
func (c *Config) IsGrayGlacier(num uint64) bool {
	return isForked(c.GrayGlacierBlock, num)
}

// IsShanghai returns whether time is either equal to the Shanghai fork time or greater.
func (c *Config) IsShanghai(time uint64) bool {
	return isForked(c.ShanghaiTime, time)
}

// IsAgra returns whether num is either equal to the Agra fork block or greater.
// The Agra hard fork is based on the Shanghai hard fork, but it doesn't include withdrawals.
// Also Agra is activated based on the block number rather than the timestamp.
// Refer to https://forum.polygon.technology/t/pip-28-agra-hardfork
func (c *Config) IsAgra(num uint64) bool {
	if c == nil || c.Bor == nil {
		return false
	}
	return isForked(c.Bor.AgraBlock, num)
}

// IsCancun returns whether time is either equal to the Cancun fork time or greater.
func (c *Config) IsCancun(time uint64) bool {
	return isForked(c.CancunTime, time)
}

// IsPrague returns whether time is either equal to the Prague fork time or greater.
func (c *Config) IsPrague(time uint64) bool {
	return isForked(c.PragueTime, time)
}

func (c *Config) IsBedrock(num uint64) bool {
	return isForked(c.BedrockBlock, num)
}

func (c *Config) IsRegolith(time uint64) bool {
	return isForked(c.RegolithTime, time)
}

func (c *Config) IsCanyon(time uint64) bool {
	return isForked(c.CanyonTime, time)
}

// IsOptimism returns whether the node is an optimism node or not.
func (c *Config) IsOptimism() bool {
	return c.Optimism != nil
}

func (c *Config) IsOptimismBedrock(num uint64) bool {
	return c.IsOptimism() && c.IsBedrock(num)
}

func (c *Config) IsOptimismRegolith(time uint64) bool {
	// Optimism op-geth has additional complexity which is not yet ported here.
	return /* c.IsOptimism() && */ c.IsRegolith(time)
}

func (c *Config) IsOptimismCanyon(time uint64) bool {
	return c.IsOptimism() && c.IsCanyon(time)
}

// IsOptimismPreBedrock returns true iff this is an optimism node & bedrock is not yet active
func (c *Config) IsOptimismPreBedrock(num uint64) bool {
	return c.IsOptimism() && !c.IsBedrock(num)
}

func (c *Config) GetBurntContract(num uint64) *common.Address {
	if len(c.BurntContract) == 0 {
		return nil
	}
	addr := borKeyValueConfigHelper(c.BurntContract, num)
	return &addr
}

// BaseFeeChangeDenominator bounds the amount the base fee can change between blocks.
func (c *Config) BaseFeeChangeDenominator(defaultParam, time uint64) uint64 {
	if c.IsOptimism() {
		if c.IsCanyon(time) {
			return c.Optimism.EIP1559DenominatorCanyon
		}
		return c.Optimism.EIP1559Denominator
	}
	return defaultParam
}

// ElasticityMultiplier bounds the maximum gas limit an EIP-1559 block may have.
func (c *Config) ElasticityMultiplier(defaultParam int) uint64 {
	if c.IsOptimism() {
		return c.Optimism.EIP1559Elasticity
	}
	return uint64(defaultParam)
}

func (c *Config) GetMinBlobGasPrice() uint64 {
	if c.MinBlobGasPrice != nil {
		return *c.MinBlobGasPrice
	}
	return 1 // MIN_BLOB_GASPRICE (EIP-4844)
}

func (c *Config) GetMaxBlobGasPerBlock() uint64 {
	if c.MaxBlobGasPerBlock != nil {
		return *c.MaxBlobGasPerBlock
	}
	return 786432 // MAX_BLOB_GAS_PER_BLOCK (EIP-4844)
}

func (c *Config) GetTargetBlobGasPerBlock() uint64 {
	if c.TargetBlobGasPerBlock != nil {
		return *c.TargetBlobGasPerBlock
	}
	return 393216 // TARGET_BLOB_GAS_PER_BLOCK (EIP-4844)
}

func (c *Config) GetBlobGasPriceUpdateFraction() uint64 {
	if c.BlobGasPriceUpdateFraction != nil {
		return *c.BlobGasPriceUpdateFraction
	}
	return 3338477 // BLOB_GASPRICE_UPDATE_FRACTION (EIP-4844)
}

func (c *Config) GetMaxBlobsPerBlock() uint64 {
	return c.GetMaxBlobGasPerBlock() / fixedgas.BlobGasPerBlob
}

func (c *Config) IsBobaLegacyBlock(num uint64) bool {
	// Boba Mainnet
	if BobaMainnetChainId.Cmp(c.ChainID) == 0 {
		return c.BedrockBlock.Uint64() > num
	}
	// Boba Sepolia
	if BobaSepoliaChainId.Cmp(c.ChainID) == 0 {
		return c.BedrockBlock.Uint64() > num
	}
	return false
}

func (c *Config) GetBobaGenesisGasLimit() int {
	// Boba Mainnet
	if BobaMainnetChainId.Cmp(c.ChainID) == 0 {
		return BobaMainnetGenesisGasLimit
	}
	// Boba Sepolia
	if BobaSepoliaChainId.Cmp(c.ChainID) == 0 {
		return BobaSepoliaGenesisGasLimit
	}
	return 11000000
}

func (c *Config) GetBobaGenesisCoinbase() string {
	// Boba Mainnet
	if BobaMainnetChainId.Cmp(c.ChainID) == 0 {
		return BobaMainnetGenesisCoinbase
	}
	// Boba Sepolia
	if BobaSepoliaChainId.Cmp(c.ChainID) == 0 {
		return BobaSepoliaGenesisCoinbase
	}
	return "0x0000000000000000000000000000000000000000"
}

func (c *Config) GetBobaGenesisExtraData() string {
	// Boba Mainnet
	if BobaMainnetChainId.Cmp(c.ChainID) == 0 {
		return BobaMainnetGenesisExtraData
	}
	// Boba Sepolia
	if BobaSepoliaChainId.Cmp(c.ChainID) == 0 {
		return BobaSepoliaGenesisExtraData
	}
	return ""
}

func (c *Config) GetBobaGenesisRoot() string {
	// Boba Mainnet
	if BobaMainnetChainId.Cmp(c.ChainID) == 0 {
		return BobaMainnetGenesisRoot
	}
	// Boba Sepolia
	if BobaSepoliaChainId.Cmp(c.ChainID) == 0 {
		return BobaSepoliaGenesisRoot
	}
	return ""
}

// CheckCompatible checks whether scheduled fork transitions have been imported
// with a mismatching chain configuration.
func (c *Config) CheckCompatible(newcfg *Config, height uint64) *ConfigCompatError {
	bhead := height

	// Iterate checkCompatible to find the lowest conflict.
	var lasterr *ConfigCompatError
	for {
		err := c.checkCompatible(newcfg, bhead)
		if err == nil || (lasterr != nil && err.RewindTo == lasterr.RewindTo) {
			break
		}
		lasterr = err
		bhead = err.RewindTo
	}
	return lasterr
}

type forkBlockNumber struct {
	name        string
	blockNumber *big.Int
	optional    bool // if true, the fork may be nil and next fork is still allowed
}

func (c *Config) forkBlockNumbers() []forkBlockNumber {
	return []forkBlockNumber{
		{name: "homesteadBlock", blockNumber: c.HomesteadBlock},
		{name: "daoForkBlock", blockNumber: c.DAOForkBlock, optional: true},
		{name: "eip150Block", blockNumber: c.TangerineWhistleBlock},
		{name: "eip155Block", blockNumber: c.SpuriousDragonBlock},
		{name: "byzantiumBlock", blockNumber: c.ByzantiumBlock},
		{name: "constantinopleBlock", blockNumber: c.ConstantinopleBlock},
		{name: "petersburgBlock", blockNumber: c.PetersburgBlock},
		{name: "istanbulBlock", blockNumber: c.IstanbulBlock},
		{name: "muirGlacierBlock", blockNumber: c.MuirGlacierBlock, optional: true},
		{name: "berlinBlock", blockNumber: c.BerlinBlock},
		{name: "londonBlock", blockNumber: c.LondonBlock},
		{name: "arrowGlacierBlock", blockNumber: c.ArrowGlacierBlock, optional: true},
		{name: "grayGlacierBlock", blockNumber: c.GrayGlacierBlock, optional: true},
		{name: "mergeNetsplitBlock", blockNumber: c.MergeNetsplitBlock, optional: true},
	}
}

// CheckConfigForkOrder checks that we don't "skip" any forks
func (c *Config) CheckConfigForkOrder() error {
	if c != nil && c.ChainID != nil && c.ChainID.Uint64() == 77 {
		return nil
	}

	var lastFork forkBlockNumber

	for _, fork := range c.forkBlockNumbers() {
		if lastFork.name != "" {
			// Next one must be higher number
			if lastFork.blockNumber == nil && fork.blockNumber != nil {
				return fmt.Errorf("unsupported fork ordering: %v not enabled, but %v enabled at %v",
					lastFork.name, fork.name, fork.blockNumber)
			}
			if lastFork.blockNumber != nil && fork.blockNumber != nil {
				if lastFork.blockNumber.Cmp(fork.blockNumber) > 0 {
					return fmt.Errorf("unsupported fork ordering: %v enabled at %v, but %v enabled at %v",
						lastFork.name, lastFork.blockNumber, fork.name, fork.blockNumber)
				}
			}
			// If it was optional and not set, then ignore it
		}
		if !fork.optional || fork.blockNumber != nil {
			lastFork = fork
		}
	}
	return nil
}

func (c *Config) checkCompatible(newcfg *Config, head uint64) *ConfigCompatError {
	// returns true if a fork scheduled at s1 cannot be rescheduled to block s2 because head is already past the fork.
	incompatible := func(s1, s2 *big.Int, head uint64) bool {
		return (isForked(s1, head) || isForked(s2, head)) && !numEqual(s1, s2)
	}

	// Ethereum mainnet forks
	if incompatible(c.HomesteadBlock, newcfg.HomesteadBlock, head) {
		return newCompatError("Homestead fork block", c.HomesteadBlock, newcfg.HomesteadBlock)
	}
	if incompatible(c.DAOForkBlock, newcfg.DAOForkBlock, head) {
		return newCompatError("DAO fork block", c.DAOForkBlock, newcfg.DAOForkBlock)
	}
	if incompatible(c.TangerineWhistleBlock, newcfg.TangerineWhistleBlock, head) {
		return newCompatError("Tangerine Whistle fork block", c.TangerineWhistleBlock, newcfg.TangerineWhistleBlock)
	}
	if incompatible(c.SpuriousDragonBlock, newcfg.SpuriousDragonBlock, head) {
		return newCompatError("Spurious Dragon fork block", c.SpuriousDragonBlock, newcfg.SpuriousDragonBlock)
	}
	if c.IsSpuriousDragon(head) && !numEqual(c.ChainID, newcfg.ChainID) {
		return newCompatError("EIP155 chain ID", c.SpuriousDragonBlock, newcfg.SpuriousDragonBlock)
	}
	if incompatible(c.ByzantiumBlock, newcfg.ByzantiumBlock, head) {
		return newCompatError("Byzantium fork block", c.ByzantiumBlock, newcfg.ByzantiumBlock)
	}
	if incompatible(c.ConstantinopleBlock, newcfg.ConstantinopleBlock, head) {
		return newCompatError("Constantinople fork block", c.ConstantinopleBlock, newcfg.ConstantinopleBlock)
	}
	if incompatible(c.PetersburgBlock, newcfg.PetersburgBlock, head) {
		// the only case where we allow Petersburg to be set in the past is if it is equal to Constantinople
		// mainly to satisfy fork ordering requirements which state that Petersburg fork be set if Constantinople fork is set
		if incompatible(c.ConstantinopleBlock, newcfg.PetersburgBlock, head) {
			return newCompatError("Petersburg fork block", c.PetersburgBlock, newcfg.PetersburgBlock)
		}
	}
	if incompatible(c.IstanbulBlock, newcfg.IstanbulBlock, head) {
		return newCompatError("Istanbul fork block", c.IstanbulBlock, newcfg.IstanbulBlock)
	}
	if incompatible(c.MuirGlacierBlock, newcfg.MuirGlacierBlock, head) {
		return newCompatError("Muir Glacier fork block", c.MuirGlacierBlock, newcfg.MuirGlacierBlock)
	}
	if incompatible(c.BerlinBlock, newcfg.BerlinBlock, head) {
		return newCompatError("Berlin fork block", c.BerlinBlock, newcfg.BerlinBlock)
	}
	if incompatible(c.LondonBlock, newcfg.LondonBlock, head) {
		return newCompatError("London fork block", c.LondonBlock, newcfg.LondonBlock)
	}
	if incompatible(c.ArrowGlacierBlock, newcfg.ArrowGlacierBlock, head) {
		return newCompatError("Arrow Glacier fork block", c.ArrowGlacierBlock, newcfg.ArrowGlacierBlock)
	}
	if incompatible(c.GrayGlacierBlock, newcfg.GrayGlacierBlock, head) {
		return newCompatError("Gray Glacier fork block", c.GrayGlacierBlock, newcfg.GrayGlacierBlock)
	}
	if incompatible(c.MergeNetsplitBlock, newcfg.MergeNetsplitBlock, head) {
		return newCompatError("Merge netsplit block", c.MergeNetsplitBlock, newcfg.MergeNetsplitBlock)
	}

	return nil
}

func numEqual(x, y *big.Int) bool {
	if x == nil {
		return y == nil
	}
	if y == nil {
		return x == nil
	}
	return x.Cmp(y) == 0
}

// ConfigCompatError is raised if the locally-stored blockchain is initialised with a
// ChainConfig that would alter the past.
type ConfigCompatError struct {
	What string
	// block numbers of the stored and new configurations
	StoredConfig, NewConfig *big.Int
	// the block number to which the local chain must be rewound to correct the error
	RewindTo uint64
}

func newCompatError(what string, storedblock, newblock *big.Int) *ConfigCompatError {
	var rew *big.Int
	switch {
	case storedblock == nil:
		rew = newblock
	case newblock == nil || storedblock.Cmp(newblock) < 0:
		rew = storedblock
	default:
		rew = newblock
	}
	err := &ConfigCompatError{what, storedblock, newblock, 0}
	if rew != nil && rew.Sign() > 0 {
		err.RewindTo = rew.Uint64() - 1
	}
	return err
}

func (err *ConfigCompatError) Error() string {
	return fmt.Sprintf("mismatching %s in database (have %d, want %d, rewindto %d)", err.What, err.StoredConfig, err.NewConfig, err.RewindTo)
}

// EthashConfig is the consensus engine configs for proof-of-work based sealing.
type EthashConfig struct{}

// String implements the stringer interface, returning the consensus engine details.
func (c *EthashConfig) String() string {
	return "ethash"
}

// CliqueConfig is the consensus engine configs for proof-of-authority based sealing.
type CliqueConfig struct {
	Period uint64 `json:"period"` // Number of seconds between blocks to enforce
	Epoch  uint64 `json:"epoch"`  // Epoch length to reset votes and checkpoint
}

// String implements the stringer interface, returning the consensus engine details.
func (c *CliqueConfig) String() string {
	return "clique"
}

// OptimismConfig is the optimism config.
type OptimismConfig struct {
	EIP1559Elasticity        uint64 `json:"eip1559Elasticity"`
	EIP1559Denominator       uint64 `json:"eip1559Denominator"`
	EIP1559DenominatorCanyon uint64 `json:"eip1559DenominatorCanyon"`
}

// String implements the stringer interface, returning the optimism fee config details.
func (o *OptimismConfig) String() string {
	return "optimism"
}

// BorConfig is the consensus engine configs for Matic bor based sealing.
type BorConfig struct {
	Period                map[string]uint64 `json:"period"`                // Number of seconds between blocks to enforce
	ProducerDelay         map[string]uint64 `json:"producerDelay"`         // Number of seconds delay between two producer interval
	Sprint                map[string]uint64 `json:"sprint"`                // Epoch length to proposer
	BackupMultiplier      map[string]uint64 `json:"backupMultiplier"`      // Backup multiplier to determine the wiggle time
	ValidatorContract     string            `json:"validatorContract"`     // Validator set contract
	StateReceiverContract string            `json:"stateReceiverContract"` // State receiver contract

	OverrideStateSyncRecords map[string]int         `json:"overrideStateSyncRecords"` // override state records count
	BlockAlloc               map[string]interface{} `json:"blockAlloc"`

	JaipurBlock                *big.Int          `json:"jaipurBlock"`                // Jaipur switch block (nil = no fork, 0 = already on jaipur)
	DelhiBlock                 *big.Int          `json:"delhiBlock"`                 // Delhi switch block (nil = no fork, 0 = already on delhi)
	IndoreBlock                *big.Int          `json:"indoreBlock"`                // Indore switch block (nil = no fork, 0 = already on indore)
	AgraBlock                  *big.Int          `json:"agraBlock"`                  // Agra switch block (nil = no fork, 0 = already in agra)
	StateSyncConfirmationDelay map[string]uint64 `json:"stateSyncConfirmationDelay"` // StateSync Confirmation Delay, in seconds, to calculate `to`

	sprints sprints
}

// String implements the stringer interface, returning the consensus engine details.
func (b *BorConfig) String() string {
	return "bor"
}

func (c *BorConfig) CalculateProducerDelay(number uint64) uint64 {
	return borKeyValueConfigHelper(c.ProducerDelay, number)
}

func (c *BorConfig) CalculateSprint(number uint64) uint64 {
	if c.sprints == nil {
		c.sprints = asSprints(c.Sprint)
	}

	for i := 0; i < len(c.sprints)-1; i++ {
		if number >= c.sprints[i].from && number < c.sprints[i+1].from {
			return c.sprints[i].size
		}
	}

	return c.sprints[len(c.sprints)-1].size
}

func (c *BorConfig) CalculateSprintCount(from, to uint64) int {
	switch {
	case from > to:
		return 0
	case from < to:
		to--
	}

	if c.sprints == nil {
		c.sprints = asSprints(c.Sprint)
	}

	count := uint64(0)
	startCalc := from

	zeroth := func(boundary uint64, size uint64) uint64 {
		if boundary%size == 0 {
			return 1
		}

		return 0
	}

	for i := 0; i < len(c.sprints)-1; i++ {
		if startCalc >= c.sprints[i].from && startCalc < c.sprints[i+1].from {
			if to >= c.sprints[i].from && to < c.sprints[i+1].from {
				if startCalc == to {
					return int(count + zeroth(startCalc, c.sprints[i].size))
				}
				return int(count + zeroth(startCalc, c.sprints[i].size) + (to-startCalc)/c.sprints[i].size)
			} else {
				endCalc := c.sprints[i+1].from - 1
				count += zeroth(startCalc, c.sprints[i].size) + (endCalc-startCalc)/c.sprints[i].size
				startCalc = endCalc + 1
			}
		}
	}

	if startCalc == to {
		return int(count + zeroth(startCalc, c.sprints[len(c.sprints)-1].size))
	}

	return int(count + zeroth(startCalc, c.sprints[len(c.sprints)-1].size) + (to-startCalc)/c.sprints[len(c.sprints)-1].size)
}

func (c *BorConfig) CalculateBackupMultiplier(number uint64) uint64 {
	return borKeyValueConfigHelper(c.BackupMultiplier, number)
}

func (c *BorConfig) CalculatePeriod(number uint64) uint64 {
	return borKeyValueConfigHelper(c.Period, number)
}

func (c *BorConfig) IsJaipur(number uint64) bool {
	return isForked(c.JaipurBlock, number)
}

func (c *BorConfig) IsDelhi(number uint64) bool {
	return isForked(c.DelhiBlock, number)
}

func (c *BorConfig) IsIndore(number uint64) bool {
	return isForked(c.IndoreBlock, number)
}

func (c *BorConfig) CalculateStateSyncDelay(number uint64) uint64 {
	return borKeyValueConfigHelper(c.StateSyncConfirmationDelay, number)
}

func borKeyValueConfigHelper[T uint64 | common.Address](field map[string]T, number uint64) T {
	fieldUint := make(map[uint64]T)
	for k, v := range field {
		keyUint, err := strconv.ParseUint(k, 10, 64)
		if err != nil {
			panic(err)
		}
		fieldUint[keyUint] = v
	}

	keys := common.SortedKeys(fieldUint)

	for i := 0; i < len(keys)-1; i++ {
		if number >= keys[i] && number < keys[i+1] {
			return fieldUint[keys[i]]
		}
	}

	return fieldUint[keys[len(keys)-1]]
}

type sprint struct {
	from, size uint64
}

type sprints []sprint

func (s sprints) Len() int {
	return len(s)
}

func (s sprints) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sprints) Less(i, j int) bool {
	return s[i].from < s[j].from
}

func asSprints(configSprints map[string]uint64) sprints {
	sprints := make(sprints, len(configSprints))

	i := 0
	for key, value := range configSprints {
		sprints[i].from, _ = strconv.ParseUint(key, 10, 64)
		sprints[i].size = value
		i++
	}

	sort.Sort(sprints)

	return sprints
}

// Rules is syntactic sugar over Config. It can be used for functions
// that do not have or require information about the block.
//
// Rules is a one time interface meaning that it shouldn't be used in between transition
// phases.
type Rules struct {
	ChainID                                                 *big.Int
	IsHomestead, IsTangerineWhistle, IsSpuriousDragon       bool
	IsByzantium, IsConstantinople, IsPetersburg, IsIstanbul bool
	IsBerlin, IsLondon, IsShanghai, IsCancun, IsPrague      bool
	IsAura                                                  bool
	IsBedrock, IsOptimismRegolith                           bool
}

// Rules ensures c's ChainID is not nil and returns a new Rules instance
func (c *Config) Rules(num uint64, time uint64) *Rules {
	chainID := c.ChainID
	if chainID == nil {
		chainID = new(big.Int)
	}

	return &Rules{
		ChainID:            new(big.Int).Set(chainID),
		IsHomestead:        c.IsHomestead(num),
		IsTangerineWhistle: c.IsTangerineWhistle(num),
		IsSpuriousDragon:   c.IsSpuriousDragon(num),
		IsByzantium:        c.IsByzantium(num),
		IsConstantinople:   c.IsConstantinople(num),
		IsPetersburg:       c.IsPetersburg(num),
		IsIstanbul:         c.IsIstanbul(num),
		IsBerlin:           c.IsBerlin(num),
		IsLondon:           c.IsLondon(num),
		IsShanghai:         c.IsShanghai(time) || c.IsAgra(num),
		IsCancun:           c.IsCancun(time),
		IsPrague:           c.IsPrague(time),
		IsBedrock:          c.IsBedrock(num),
		IsOptimismRegolith: c.IsOptimismRegolith(time),
		IsAura:             c.Aura != nil,
	}
}

// isForked returns whether a fork scheduled at block s is active at the given head block.
func isForked(s *big.Int, head uint64) bool {
	if s == nil {
		return false
	}
	return s.Uint64() <= head
}
