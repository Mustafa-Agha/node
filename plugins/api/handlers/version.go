package handlers

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"

	"github.com/Mustafa-Agha/node/version"
)

// CLIVersionReqHandler handles requests to the cli version REST handler endpoint
func CLIVersionReqHandler(w http.ResponseWriter, r *http.Request) {
	v := version.GetVersion()
	_, _ = w.Write([]byte(v))
}

// NodeVersionReqHandler handles requests to the connected node version REST handler endpoint
func NodeVersionReqHandler(ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		version, err := ctx.Query("/app/version", nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprintf("Could't query version. Error: %s", err.Error())))
			return
		}
		_, _ = w.Write(version)
	}
}
