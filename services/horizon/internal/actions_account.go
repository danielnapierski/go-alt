package horizon

import (
	"github.com/danielnapierski/go-alt/protocols/horizon"
	"github.com/danielnapierski/go-alt/services/horizon/internal/db2/core"
	"github.com/danielnapierski/go-alt/services/horizon/internal/db2/history"
	"github.com/danielnapierski/go-alt/services/horizon/internal/render/sse"
	"github.com/danielnapierski/go-alt/services/horizon/internal/resourceadapter"
	"github.com/danielnapierski/go-alt/support/render/hal"
)

// This file contains the actions:
//
// AccountShowAction: details for single account (including stellar-core state)

// AccountShowAction renders a account summary found by its address.
type AccountShowAction struct {
	Action
	Address        string
	HistoryRecord  history.Account
	CoreData       []core.AccountData
	CoreRecord     core.Account
	CoreSigners    []core.Signer
	CoreTrustlines []core.Trustline
	Resource       horizon.Account
}

// JSON is a method for actions.JSON
func (action *AccountShowAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadRecord,
		action.loadResource,
		func() {
			hal.Render(action.W, action.Resource)
		},
	)
}

// SSE is a method for actions.SSE
func (action *AccountShowAction) SSE(stream sse.Stream) {

	action.Do(
		action.loadParams,
		action.loadRecord,
		action.loadResource,
		func() {
			stream.SetLimit(10)
			stream.Send(sse.Event{Data: action.Resource})
		},
	)
}

func (action *AccountShowAction) loadParams() {
	action.Address = action.GetAddress("account_id")
}

func (action *AccountShowAction) loadRecord() {
	app := AppFromContext(action.R.Context())
	protocolVersion := app.protocolVersion

	action.Err = action.CoreQ().
		AccountByAddress(&action.CoreRecord, action.Address, protocolVersion)
	if action.Err != nil {
		return
	}

	action.Err = action.CoreQ().
		AllDataByAddress(&action.CoreData, action.Address)
	if action.Err != nil {
		return
	}

	action.Err = action.CoreQ().
		SignersByAddress(&action.CoreSigners, action.Address)
	if action.Err != nil {
		return
	}

	action.Err = action.CoreQ().
		TrustlinesByAddress(&action.CoreTrustlines, action.Address, protocolVersion)
	if action.Err != nil {
		return
	}

	action.Err = action.HistoryQ().
		AccountByAddress(&action.HistoryRecord, action.Address)

	// Do not fail when we cannot find the history record... it probably just
	// means that the account was created outside of our known history range.
	if action.HistoryQ().NoRows(action.Err) {
		action.Err = nil
	}

	if action.Err != nil {
		return
	}
}

func (action *AccountShowAction) loadResource() {
	action.Err = resourceadapter.PopulateAccount(
		action.R.Context(),
		&action.Resource,
		action.CoreRecord,
		action.CoreData,
		action.CoreSigners,
		action.CoreTrustlines,
		action.HistoryRecord,
	)
}
