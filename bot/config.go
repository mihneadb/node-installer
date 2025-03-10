package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const (
	VersionJSON     = "https://raw.githubusercontent.com/SwingbyProtocol/node-installer/master/.version.json"
	DataPath        = "./data"
	Network1        = "btc_eth"
	Network2        = "btc_bsc"
	Network3        = "tbtc_goerli"
	Network4        = "tbtc_bsc"
	GethLockVersion = "Geth/v1.10.1"
	BSCLockVersion  = "Geth/v1.1.0-beta"
	BTCLockVersion  = "210000"
)

var (
	Networks = map[string]string{
		"1": Network1,
		"2": Network2,
		"3": Network3,
		"4": Network4,
	}
	WalletContract = map[string]string{
		Network1: "0xbe83f11d3900F3a13d8D12fB62F5e85646cDA45e",
		Network2: "0xaD22900062e4cd766102A1f33E530F5303fe1aDF",
	}
	LPtokenContract = map[string]string{
		Network1: "0x22883a3db06737ece21f479a8009b8b9f22b6cc9",
		Network2: "0xdBa68BeF9b541999Fd9650FF72C19d5E1ceeCd10",
	}
	BTCTContract = map[string]string{
		Network1: "0x2260fac5e5542a773aa44fbcfedf7c193bc2c599",
		Network2: "0x7130d2a12b9bcbfae4f2634d864a1ee1ce3ead9c",
	}
	BootstrapNodeMain = map[string][]string{
		Network1: {
			"49.12.68.127:12131",  // https://moonfly-re-0078.yen.farm
			"49.12.7.120:12132",   // https://livemex-re-0079.yen.farm
			"116.203.56.22:12133", // https://motion-re-0080.yen.farm
		},
		Network2: {
			"163.172.141.211:12124", // https://ra-cailum.zoo.farm
			"51.158.68.138:12125",   // https://irish.zoo.farm
			"51.159.134.173:12126",  // https://gwaden.yen.farm
		},
	}
	stopTrigger = map[string]string{
		Network1: "https://btc-wbtc-mainnet.s3.eu-central-1.amazonaws.com/platform_status.json",
		Network2: "https://btc-bsc-mainnet.s3-ap-southeast-1.amazonaws.com/platform_status.json",
	}
	epochBlock = map[string]int{
		Network1: 3,
		Network2: 15,
	}
	threshold = map[string]int{
		Network1: 31,
		Network2: 31,
	}
	maxShare = map[string]int{
		Network1: 50,
		Network2: 50,
	}
	maxNode = map[string]int{
		Network1: 60,
		Network2: 60,
	}
	keygenPeer = map[string]int{
		Network1: 32,
		Network2: 35,
	}
	syncSnapshotBytes = map[string]int{
		Network1: 1175750002860,
		Network2: 971003535776,
	}
	minimumMountPathSizeMiB = map[string]int{
		Network1: 1430511,
		Network2: 965978,
	}
)

const (
	GethRPC        = "http://10.2.0.1:8545"
	BscRPC         = "http://10.2.0.1:8575"
	BlockBookBTC   = "http://10.2.0.1:9130"
	BlockBookBTCWS = "ws://10.2.0.1:9130/websocket"
	BlockBookETH   = "http://10.2.0.1:9131"
	BlockBookETHWS = "ws://10.2.0.1:9131/websocket"
	BlockBookBSC   = "http://10.2.0.1:9132"
	BlockBookBSCWS = "ws://10.2.0.1:9132/websocket"
)

var BnbSeedNodesMain = []string{
	"tcp://dataseed2.defibit.io:80",
}

const baseConfig = `
[p2p]
moniker = "**node_moniker_placeholder**"
listen = "0.0.0.0"
port = 12121

[general]
epoch_blocks = **epoch_block**

[rest]
listen = "0.0.0.0"
port = 8067
tls_enabled = false
preferred_uri = "**node_preferred_uri**"

[logger]
level = "INFO"
max_file_size_MB = 10
max_backup_files = 100
max_retain_files_days = 14
use_console_logger = true
use_file_logger = true
compress = true

[swaps]
testnet = **is_testnet**
coin_1 = "**coin_A**"
coin_2 = "**coin_B**"
fee_percent = 0.2
stop_trigger_uri = "**stop_trigger_uri**"
# (using defaults in code)
# max_amount = 1
# min_amount_swap = 0.0004
# min_amount_refund = 0.001

[tss]
threshold = **threshold_placeholder**
max_shares = **max_shares**
max_nodes = **max_nodes**
keygen_peers = **keygen_peers**
keygen_until = "2020-12-13T12:00:00Z"

[btc]
rest_uri = "**btc_blockbook_endpoint**"
ws_uri = "**btc_blockbook_ws_endpoint**"
miner_fee = 0.0002

[eth]
rpc_uri = "**eth_rpc_endpoint**"
rest_uri = "**eth_blockbook_endpoint**"
ws_uri = "**eth_blockbook_ws_endpoint**"
wallet_contract_addr = "**eth_wallet_contract**"
lp_token_contract_addr = "**eth_lpt_contract**"
btc_token_contract_addr = "**btc_token_contract_addr**"

[bsc_fees]
miner_fee = 0.000015

[bnb]
rpc_uri = "**rpc_uri_placeholder**"
http_uri = "https://explorer.binance.org"
stake_addr = "**stake_addr**"
`

type NodeConfig struct {
	Network          string
	Moniker          string
	BootstrapNode    []string
	Domain           string
	PreferredURI     string
	BNBSeed          string
	CoinA            string
	CoinB            string
	RewardAddressETH string
	RewardAddressBNB string
	GethRPC          string
	BlockBookBTC     string
	BlockBookBTCWS   string
	BlockBookETH     string
	BlockBookETHWS   string
	StakeAddr        string
	StakeTx          string
	WalletContract   string
	LPtoken          string
	BTCTContract     string
	StopTrigger      string
	Memo             string
	KeygenUntil      string
	IsTestnet        bool
	Threshold        int
	EpochBlock       int
	MaxShares        int
	MaxNodes         int
	KeygenPeers      int
}

func NewNodeConfig() *NodeConfig {
	initTime := time.Date(2014, time.December, 31, 12, 13, 24, 0, time.UTC)
	nConf := &NodeConfig{
		CoinA:          "WBTC",
		CoinB:          "BTC",
		GethRPC:        GethRPC,
		BNBSeed:        BnbSeedNodesMain[0],
		BlockBookBTC:   BlockBookBTC,
		BlockBookBTCWS: BlockBookBTCWS,
		BlockBookETH:   BlockBookETH,
		BlockBookETHWS: BlockBookETHWS,
		KeygenUntil:    initTime.Format(time.RFC3339),
		Network:        Network1,
		BootstrapNode:  BootstrapNodeMain[Network1],
		Moniker:        "Default Node",
		WalletContract: WalletContract[Network1],
		LPtoken:        LPtokenContract[Network1],
		BTCTContract:   BTCTContract[Network1],
		StopTrigger:    stopTrigger[Network1],
	}
	return nConf
}

func (n *NodeConfig) SetNetwork(network string) {
	n.Network = network
	n.IsTestnet = false
	n.WalletContract = WalletContract[network]
	n.LPtoken = LPtokenContract[network]
	n.BootstrapNode = BootstrapNodeMain[network]
	n.BTCTContract = BTCTContract[network]
	n.BootstrapNode = BootstrapNodeMain[network]
	n.StopTrigger = stopTrigger[network]
	n.EpochBlock = epochBlock[network]
	n.Threshold = threshold[network]
	n.MaxShares = maxShare[network]
	n.MaxNodes = maxNode[network]
	n.KeygenPeers = keygenPeer[network]

	switch n.Network {
	case Network1:
		n.CoinA = "WBTC"
		n.CoinB = "BTC"
	case Network2:
		n.CoinA = "BTCB"
		n.CoinB = "BTC"
	}
}

func (n *NodeConfig) SetGlobalNode() {
	switch n.Network {
	case Network1:
		n.BlockBookBTC = "http://btc-eth-indexer.swingby.network:9130"
		n.BlockBookBTCWS = "wss://btc-eth-indexer.swingby.network/btc-websocket"
		n.GethRPC = "http://51.159.56.104:8545" // foundation geth_1
		n.BlockBookETH = "http://btc-eth-indexer.swingby.network:9131"
		n.BlockBookETHWS = "wss://btc-eth-indexer.swingby.network/eth-websocket"
	case Network2:
		n.BlockBookBTC = "https://btc-eth-indexer-2.swingby.network/bb-btc"
		n.BlockBookBTCWS = "wss://btc-eth-indexer-2.swingby.network/btc-websocket"
		n.GethRPC = "http://btc-bsc-indexer.swingby.network:8575" // foundation bsc_2
		n.BlockBookETH = "https://btc-bsc-indexer.swingby.network/bb-bsc"
		n.BlockBookETHWS = "wss://btc-bsc-indexer.swingby.network/bsc-websocket"
	}
}

func (n *NodeConfig) SetLocalNode() {
	n.BlockBookBTC = BlockBookBTC
	n.BlockBookBTCWS = BlockBookBTCWS
	switch n.Network {
	case Network1:
		n.GethRPC = GethRPC
		n.BlockBookETH = BlockBookETH
		n.BlockBookETHWS = BlockBookETHWS
	case Network2:
		n.GethRPC = BscRPC
		n.BlockBookETH = BlockBookBSC
		n.BlockBookETHWS = BlockBookBSCWS
	}
}

func (n *NodeConfig) SetDomain(domain string) {
	n.Domain = domain
	n.PreferredURI = fmt.Sprintf("https://%s", domain)
}

func (n *NodeConfig) checkConfig() error {
	pConfigFileName := fmt.Sprintf("%s/%s/config.toml", DataPath, n.Network)
	_, err := ioutil.ReadFile(pConfigFileName)
	if err != nil {
		return err
	}
	return nil
}

func (n *NodeConfig) storeConfigToml() error {
	pConfigFileName := fmt.Sprintf("%s/%s/config.toml", DataPath, n.Network)
	newBaseConfig := strings.ReplaceAll(baseConfig, "**node_moniker_placeholder**", n.Moniker)

	newBaseConfig = strings.ReplaceAll(newBaseConfig, "**epoch_block**", fmt.Sprintf("%d", n.EpochBlock))

	newBaseConfig = strings.ReplaceAll(newBaseConfig, "**node_preferred_uri**", n.PreferredURI)
	newBaseConfig = strings.ReplaceAll(newBaseConfig, "**coin_A**", n.CoinA)
	newBaseConfig = strings.ReplaceAll(newBaseConfig, "**coin_B**", n.CoinB)
	newBaseConfig = strings.ReplaceAll(newBaseConfig, "**stop_trigger_uri**", n.StopTrigger)
	newBaseConfig = strings.ReplaceAll(newBaseConfig, "**is_testnet**", fmt.Sprintf("%t", n.IsTestnet))

	newBaseConfig = strings.ReplaceAll(newBaseConfig, "**threshold_placeholder**", fmt.Sprintf("%d", n.Threshold))
	newBaseConfig = strings.ReplaceAll(newBaseConfig, "**max_shares**", fmt.Sprintf("%d", n.MaxShares))
	newBaseConfig = strings.ReplaceAll(newBaseConfig, "**max_nodes**", fmt.Sprintf("%d", n.MaxNodes))
	newBaseConfig = strings.ReplaceAll(newBaseConfig, "**keygen_peers**", fmt.Sprintf("%d", n.KeygenPeers))

	newBaseConfig = strings.ReplaceAll(newBaseConfig, "**btc_blockbook_endpoint**", n.BlockBookBTC)
	newBaseConfig = strings.ReplaceAll(newBaseConfig, "**btc_blockbook_ws_endpoint**", n.BlockBookBTCWS)

	newBaseConfig = strings.ReplaceAll(newBaseConfig, "**eth_rpc_endpoint**", n.GethRPC)
	newBaseConfig = strings.ReplaceAll(newBaseConfig, "**eth_blockbook_endpoint**", n.BlockBookETH)
	newBaseConfig = strings.ReplaceAll(newBaseConfig, "**eth_blockbook_ws_endpoint**", n.BlockBookETHWS)

	newBaseConfig = strings.ReplaceAll(newBaseConfig, "**eth_wallet_contract**", n.WalletContract)
	newBaseConfig = strings.ReplaceAll(newBaseConfig, "**eth_lpt_contract**", n.LPtoken)
	newBaseConfig = strings.ReplaceAll(newBaseConfig, "**btc_token_contract_addr**", n.BTCTContract)
	//newBaseConfig = strings.ReplaceAll(newBaseConfig, "**reward_address_eth**", n.RewardAddressETH)

	newBaseConfig = strings.ReplaceAll(newBaseConfig, "**rpc_uri_placeholder**", n.BNBSeed)
	newBaseConfig = strings.ReplaceAll(newBaseConfig, "**stake_tx**", n.StakeTx)
	newBaseConfig = strings.ReplaceAll(newBaseConfig, "**stake_addr**", n.StakeAddr)
	newBaseConfig = strings.ReplaceAll(newBaseConfig, "**reward_addr_bnb**", n.RewardAddressBNB)

	newConfigToml := fmt.Sprintf("%s\n", newBaseConfig)
	if err := ioutil.WriteFile(pConfigFileName, []byte(newConfigToml), os.ModePerm); err != nil {
		return err
	}
	return nil
}

func (n *NodeConfig) saveConfig() error {
	data, err := json.Marshal(n)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s/node_config.json", DataPath), []byte(data), 0600)
	if err != nil {
		return err
	}
	return nil
}

func (n *NodeConfig) loadConfig() error {
	str, err := ioutil.ReadFile(fmt.Sprintf("%s/node_config.json", DataPath))
	if err != nil {
		return err
	}
	err = json.Unmarshal(str, &n)
	if err != nil {
		return err
	}
	return nil
}
