package server

import (
	"math/big"
	"net/http"

	"github.com/danielnapierski/go-alt/services/bifrost/bitcoin"
	"github.com/danielnapierski/go-alt/services/bifrost/config"
	"github.com/danielnapierski/go-alt/services/bifrost/database"
	"github.com/danielnapierski/go-alt/services/bifrost/ethereum"
	"github.com/danielnapierski/go-alt/services/bifrost/queue"
	"github.com/danielnapierski/go-alt/services/bifrost/sse"
	"github.com/danielnapierski/go-alt/services/bifrost/stellar"
	"github.com/danielnapierski/go-alt/support/log"
)

// ProtocolVersion is the version of the protocol that Bifrost server and
// JS SDK use to communicate.
const ProtocolVersion int = 2

type Server struct {
	BitcoinListener            *bitcoin.Listener            `inject:""`
	BitcoinAddressGenerator    *bitcoin.AddressGenerator    `inject:""`
	Config                     *config.Config               `inject:""`
	Database                   database.Database            `inject:""`
	EthereumListener           *ethereum.Listener           `inject:""`
	EthereumAddressGenerator   *ethereum.AddressGenerator   `inject:""`
	StellarAccountConfigurator *stellar.AccountConfigurator `inject:""`
	TransactionsQueue          queue.Queue                  `inject:""`
	SSEServer                  sse.ServerInterface          `inject:""`

	MinimumValueBtc string
	MinimumValueEth string
	SignerPublicKey string

	minimumValueSat int64
	minimumValueWei *big.Int
	httpServer      *http.Server
	log             *log.Entry
}

type GenerateAddressResponse struct {
	ProtocolVersion int    `json:"protocol_version"`
	Chain           string `json:"chain"`
	Address         string `json:"address"`
	Signer          string `json:"signer"`
}
