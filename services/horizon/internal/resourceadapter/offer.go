package resourceadapter

import (
	"context"

	"github.com/danielnapierski/go-alt/amount"
	. "github.com/danielnapierski/go-alt/protocols/horizon"
	"github.com/danielnapierski/go-alt/services/horizon/internal/assets"
	"github.com/danielnapierski/go-alt/services/horizon/internal/db2/core"
	"github.com/danielnapierski/go-alt/services/horizon/internal/db2/history"
	"github.com/danielnapierski/go-alt/services/horizon/internal/httpx"
	"github.com/danielnapierski/go-alt/support/render/hal"
)

func PopulateOffer(ctx context.Context, dest *Offer, row core.Offer, ledger *history.Ledger) {
	dest.ID = row.OfferID
	dest.PT = row.PagingToken()
	dest.Seller = row.SellerID
	dest.Amount = amount.String(row.Amount)
	dest.PriceR.N = row.Pricen
	dest.PriceR.D = row.Priced
	dest.Price = row.PriceAsString()
	dest.Buying = Asset{
		Type:   assets.MustString(row.BuyingAssetType),
		Code:   row.BuyingAssetCode.String,
		Issuer: row.BuyingIssuer.String,
	}
	dest.Selling = Asset{
		Type:   assets.MustString(row.SellingAssetType),
		Code:   row.SellingAssetCode.String,
		Issuer: row.SellingIssuer.String,
	}
	dest.LastModifiedLedger = row.Lastmodified
	if ledger != nil {
		dest.LastModifiedTime = &ledger.ClosedAt
	}
	lb := hal.LinkBuilder{httpx.BaseURL(ctx)}
	dest.Links.Self = lb.Linkf("/offers/%d", row.OfferID)
	dest.Links.OfferMaker = lb.Linkf("/accounts/%s", row.SellerID)
	return
}
