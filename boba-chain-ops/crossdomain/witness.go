package crossdomain

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ledgerwatch/erigon-lib/common"
	"github.com/ledgerwatch/erigon-lib/common/hexutility"
	"github.com/ledgerwatch/erigon/boba-bindings/bindings"
	"github.com/ledgerwatch/erigon/common/hexutil"

	"github.com/ethereum/go-ethereum/log"
)

// SentMessage represents an entry in the JSON file that is created by
// the `migration-data` package. Each entry represents a call to the
// `LegacyMessagePasser`. The `who` should always be the
// `L2CrossDomainMessenger` and the `msg` should be an abi encoded
// `relayMessage(address,address,bytes,uint256)`
type SentMessage struct {
	Who common.Address   `json:"who"`
	Msg hexutility.Bytes `json:"msg"`
}

// NewSentMessageFromJSON will read a JSON file from disk given a path to the JSON
// file. The JSON file this function reads from disk is an output from the
// `migration-data` package.
func NewSentMessageFromJSON(path string) ([]*SentMessage, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot find sent message json at %s: %w", path, err)
	}

	if len(file) == 0 {
		log.Warn("sent message json file is empty")
		return nil, nil
	}

	var j []*SentMessage
	if err := json.Unmarshal(file, &j); err != nil {
		return nil, err
	}

	return j, nil
}

// ReadWitnessData will read messages and addresses from a raw l2geth state
// dump file.
func ReadWitnessData(path string) ([]*SentMessage, OVMETHAddresses, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot open witness data file: %w", err)
	}
	defer f.Close()

	scan := bufio.NewScanner(f)
	var witnesses []*SentMessage
	addresses := make(map[common.Address]bool)
	for scan.Scan() {
		line := scan.Text()
		splits := strings.Split(line, "|")
		if len(splits) < 2 {
			return nil, nil, fmt.Errorf("invalid line: %s", line)
		}

		switch splits[0] {
		case "MSG":
			if len(splits) != 3 {
				return nil, nil, fmt.Errorf("invalid line: %s", line)
			}

			msg := splits[2]
			// Make sure that the witness data has a 0x prefix
			if !strings.HasPrefix(msg, "0x") {
				msg = "0x" + msg
			}

			abi, err := bindings.LegacyMessagePasserMetaData.GetAbi()
			if err != nil {
				return nil, nil, fmt.Errorf("failed to get abi: %w", err)
			}

			msgB := hexutil.MustDecode(msg)
			method, err := abi.MethodById(msgB[:4])
			if err != nil {
				return nil, nil, fmt.Errorf("failed to get method: %w", err)
			}

			out, err := method.Inputs.Unpack(msgB[4:])
			if err != nil {
				return nil, nil, fmt.Errorf("failed to unpack: %w", err)
			}

			cast, ok := out[0].([]byte)
			if !ok {
				return nil, nil, fmt.Errorf("failed to cast to bytes")
			}

			witnesses = append(witnesses, &SentMessage{
				Who: common.HexToAddress(splits[1]),
				Msg: cast,
			})
		case "ETH":
			addresses[common.HexToAddress(splits[1])] = true
		default:
			return nil, nil, fmt.Errorf("invalid line: %s", line)
		}
	}

	return witnesses, addresses, nil
}

// Allowance represents the allowances that were set in the
// legacy ERC20 representation of ether
type Allowance struct {
	From common.Address `json:"fr"`
	To   common.Address `json:"to"`
}

// NewAllowances will read the ovm-allowances.json from the file system.
func NewAllowances(path string) ([]*Allowance, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot find allowances json at %s: %w", path, err)
	}

	if len(file) == 0 {
		log.Warn("allowances json file is empty")
		return nil, nil
	}

	var allowances []*Allowance
	if err := json.Unmarshal(file, &allowances); err != nil {
		return nil, err
	}

	return allowances, nil
}

// NewAddresses represents the middle layer that saves the addresses
type EthAddresses struct {
	BlockNumber int64             `json:"blockNumber"`
	Addresses   []*common.Address `json:"addresses"`
}

// OVMETHAddresses represents a list of addresses that interacted with
// the ERC20 representation of ether in the pre-bedrock system.
type OVMETHAddresses map[common.Address]bool

// NewAddresses will read an addresses.json file from the filesystem.
func NewAddresses(path string) (OVMETHAddresses, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot find addresses json at %s: %w", path, err)
	}

	var ethAddresses EthAddresses
	if err := json.Unmarshal(file, &ethAddresses); err != nil {
		return nil, err
	}

	ovmeth := make(OVMETHAddresses)
	for _, addr := range ethAddresses.Addresses {
		ovmeth[*addr] = true
	}

	return ovmeth, nil
}

// MigrationData represents all of the data required to do a migration
type MigrationData struct {
	// OvmAddresses represents the set of addresses that interacted with the
	// LegacyERC20ETH contract before the evm equivalence upgrade
	OvmAddresses OVMETHAddresses
	// EvmAddresses represents the set of addresses that interacted with the
	// LegacyERC20ETH contract after the evm equivalence upgrade
	EvmAddresses OVMETHAddresses
	// OvmAllowances represents the set of allowances in the LegacyERC20ETH from
	// before the evm equivalence upgrade
	OvmAllowances []*Allowance
	// OvmMessages represents the set of withdrawals through the
	// L2CrossDomainMessenger from before the evm equivalence upgrade
	OvmMessages []*SentMessage
	// OvmMessages represents the set of withdrawals through the
	// L2CrossDomainMessenger from after the evm equivalence upgrade
	EvmMessages []*SentMessage
}
