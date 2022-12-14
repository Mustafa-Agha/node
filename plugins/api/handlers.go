package api

import (
	"net/http"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	paramapi "github.com/cosmos/cosmos-sdk/x/paramHub/client/rest"

	hnd "github.com/Mustafa-Agha/node/plugins/api/handlers"
	"github.com/Mustafa-Agha/node/plugins/dex"
	dexapi "github.com/Mustafa-Agha/node/plugins/dex/client/rest"
	tksapi "github.com/Mustafa-Agha/node/plugins/tokens/client/rest"
	tkstore "github.com/Mustafa-Agha/node/plugins/tokens/store"
	"github.com/Mustafa-Agha/node/wire"
)

// middleware (limits, parsing, etc)

func (s *server) limitReqSize(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// reject suspiciously large form posts
		if r.ContentLength > s.maxPostSize {
			http.Error(w, "request too large", http.StatusExpectationFailed)
			return
		}
		// parse request body as multipart/form-data with maxPostSize in mind
		r.Body = http.MaxBytesReader(w, r.Body, s.maxPostSize)
		next(w, r)
	}
}

// withUrlEncForm parses application/x-www-form-urlencoded forms
func (s *server) withUrlEncForm(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
			http.Error(w, "application/x-www-form-urlencoded content-type expected", http.StatusExpectationFailed)
			return
		}
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		next(w, r)
	}
}

//nolint:unused
// withMultipartForm parses multipart/form-data forms
func (s *server) withMultipartForm(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "multipart/form-data") {
			http.Error(w, "multipart/form-data content-type expected", http.StatusExpectationFailed)
			return
		}
		err := r.ParseMultipartForm(1024)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		next(w, r)
	}
}

// withTextPlain parses text/plain forms
func (s *server) withTextPlainForm(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "text/plain") {
			http.Error(w, "text/plain content-type expected", http.StatusExpectationFailed)
			return
		}
		next(w, r)
	}
}

// -----

func (s *server) handleVersionReq() http.HandlerFunc {
	return hnd.CLIVersionReqHandler
}

func (s *server) handleNodeVersionReq() http.HandlerFunc {
	return hnd.NodeVersionReqHandler(s.ctx)
}

func (s *server) handleAccountReq(cdc *wire.Codec, ctx context.CLIContext) http.HandlerFunc {
	return hnd.AccountReqHandler(cdc, ctx)
}

func (s *server) handleSimulateReq(cdc *wire.Codec, ctx context.CLIContext) http.HandlerFunc {
	h := hnd.SimulateReqHandler(cdc, ctx)
	return s.withTextPlainForm(s.limitReqSize(h))
}

func (s *server) handleBEP2PairsReq(cdc *wire.Codec, ctx context.CLIContext) http.HandlerFunc {
	return dexapi.GetPairsReqHandler(cdc, ctx, dex.DexAbciQueryPrefix)
}

func (s *server) handleMiniPairsReq(cdc *wire.Codec, ctx context.CLIContext) http.HandlerFunc {
	return dexapi.GetPairsReqHandler(cdc, ctx, dex.DexMiniAbciQueryPrefix)
}

func (s *server) handleDexDepthReq(cdc *wire.Codec, ctx context.CLIContext) http.HandlerFunc {
	return dexapi.DepthReqHandler(cdc, ctx)
}

func (s *server) handleDexOrderReq(cdc *wire.Codec, ctx context.CLIContext, accStoreName string) http.HandlerFunc {
	h := dexapi.PutOrderReqHandler(cdc, ctx, accStoreName)
	return s.withUrlEncForm(s.limitReqSize(h))
}

func (s *server) handleDexOpenOrdersReq(cdc *wire.Codec, ctx context.CLIContext) http.HandlerFunc {
	return dexapi.OpenOrdersReqHandler(cdc, ctx)
}

func (s *server) handleTokenReq(cdc *wire.Codec, ctx context.CLIContext) http.HandlerFunc {
	return tksapi.GetTokenReqHandler(cdc, ctx, false)
}

func (s *server) handleTokensReq(cdc *wire.Codec, ctx context.CLIContext) http.HandlerFunc {
	return tksapi.GetTokensReqHandler(cdc, ctx, false)
}

func (s *server) handleMiniTokenReq(cdc *wire.Codec, ctx context.CLIContext) http.HandlerFunc {
	return tksapi.GetTokenReqHandler(cdc, ctx, true)
}

func (s *server) handleMiniTokensReq(cdc *wire.Codec, ctx context.CLIContext) http.HandlerFunc {
	return tksapi.GetTokensReqHandler(cdc, ctx, true)
}

func (s *server) handleBalancesReq(cdc *wire.Codec, ctx context.CLIContext, tokens tkstore.Mapper) http.HandlerFunc {
	return tksapi.BalancesReqHandler(cdc, ctx, tokens)
}

func (s *server) handleBalanceReq(cdc *wire.Codec, ctx context.CLIContext, tokens tkstore.Mapper) http.HandlerFunc {
	return tksapi.BalanceReqHandler(cdc, ctx, tokens)
}

func (s *server) handleFeesParamReq(cdc *wire.Codec, ctx context.CLIContext) http.HandlerFunc {
	return paramapi.GetFeesParamHandler(cdc, ctx)
}

func (s *server) handleValidatorsQueryReq(cdc *wire.Codec, ctx context.CLIContext) http.HandlerFunc {
	return hnd.ValidatorQueryReqHandler(cdc, ctx)
}

func (s *server) handleDelegatorUnbondingDelegationsQueryReq(cdc *wire.Codec, ctx context.CLIContext) http.HandlerFunc {
	return hnd.DelegatorUnbondindDelegationsQueryReqHandler(cdc, ctx)
}

func (s *server) handleTimeLocksReq(cdc *wire.Codec, ctx context.CLIContext) http.HandlerFunc {
	return tksapi.GetTimeLocksReqHandler(cdc, ctx)
}

func (s *server) handleTimeLockReq(cdc *wire.Codec, ctx context.CLIContext) http.HandlerFunc {
	return tksapi.GetTimeLockReqHandler(cdc, ctx)
}

func (s *server) handleQuerySwapReq(cdc *wire.Codec, ctx context.CLIContext) http.HandlerFunc {
	return tksapi.QuerySwapReqHandler(cdc, ctx)
}

func (s *server) handleQuerySwapIDsByCreatorReq(cdc *wire.Codec, ctx context.CLIContext) http.HandlerFunc {
	return tksapi.QuerySwapIDsByCreatorReqHandler(cdc, ctx)
}

func (s *server) handleQuerySwapIDsByRecipientReq(cdc *wire.Codec, ctx context.CLIContext) http.HandlerFunc {
	return tksapi.QuerySwapIDsByRecipientReqHandler(cdc, ctx)
}
