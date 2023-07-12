package rest

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/gorilla/mux"

	"github.com/mineplexio/mineplex-2-node/x/gravity/types"
)

func getValsetRequestHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nonce := vars[nonce]

		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/valsetRequest/%s", storeName, nonce))
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			WriteErrorResponse(w, http.StatusNotFound, "valset not found")
			return
		}

		var out types.Valset
		cliCtx.Codec.MustUnmarshalJSON(res, &out)
		PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

// USED BY RUST
func batchByNonceHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nonce := vars[nonce]
		denom := vars[tokenAddress]

		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/batch/%s/%s", storeName, nonce, denom))
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			WriteErrorResponse(w, http.StatusNotFound, "valset not found")
			return
		}

		var out types.OutgoingTxBatch
		cliCtx.Codec.MustUnmarshalJSON(res, &out)
		PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

// USED BY RUST
func lastBatchesHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/lastBatches", storeName))
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			WriteErrorResponse(w, http.StatusNotFound, "valset not found")
			return
		}

		PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

// gets all the confirm messages for a given validator set nonce
func allValsetConfirmsHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nonce := vars[nonce]

		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/valsetConfirms/%s", storeName, nonce))
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			WriteErrorResponse(w, http.StatusNotFound, "valset confirms not found")
			return
		}

		PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

// gets all the confirm messages for a given transaction batch
func allBatchConfirmsHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nonce := vars[nonce]
		denom := vars[tokenAddress]

		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/batchConfirms/%s/%s", storeName, nonce, denom))
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			WriteErrorResponse(w, http.StatusNotFound, "valset confirms not found")
			return
		}

		PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

func lastValsetRequestsHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/lastValsetRequests", storeName))
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			WriteErrorResponse(w, http.StatusNotFound, "valset requests not found")
			return
		}

		PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

func lastValsetRequestsByAddressHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		operatorAddr := vars[bech32ValidatorAddress]

		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/lastPendingValsetRequest/%s", storeName, operatorAddr))
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			WriteErrorResponse(w, http.StatusNotFound, "no pending valset requests found")
			return
		}

		var out types.Valset
		cliCtx.Codec.MustUnmarshalJSON(res, &out)
		PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

func lastBatchesByAddressHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		operatorAddr := vars[bech32ValidatorAddress]

		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/lastPendingBatchRequest/%s", storeName, operatorAddr))
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			WriteErrorResponse(w, http.StatusNotFound, "no pending valset requests found")
			return
		}

		var out types.OutgoingTxBatch
		cliCtx.Codec.MustUnmarshalJSON(res, &out)
		PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

func currentValsetHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/currentValset", storeName))
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		var out types.Valset
		cliCtx.Codec.MustUnmarshalJSON(res, &out)
		PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

func denomToERC20Handler(cliCtx client.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		denom := vars[denom]

		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/DenomToERC20/%s", storeName, denom))
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

func ERC20ToDenomHandler(cliCtx client.Context, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ERC20 := vars[tokenAddress]

		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/ERC20ToDenom/%s", storeName, ERC20))
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

func WriteErrorResponse(w http.ResponseWriter, status int, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(legacy.Cdc.MustMarshalJSON(NewErrorResponse(0, err)))
}

type ErrorResponse struct {
	Code  int    `json:"code,omitempty"`
	Error string `json:"error"`
}

// NewErrorResponse creates a new ErrorResponse instance.
func NewErrorResponse(code int, err string) ErrorResponse {
	return ErrorResponse{Code: code, Error: err}
}

// PostProcessResponse performs post processing for a REST response. The result
// returned to clients will contain two fields, the height at which the resource
// was queried at and the original result.
func PostProcessResponse(w http.ResponseWriter, ctx client.Context, resp interface{}) {
	var (
		result []byte
		err    error
	)

	if ctx.Height < 0 {
		WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("negative height in response").Error())
		return
	}

	// LegacyAmino used intentionally for REST
	marshaler := ctx.LegacyAmino

	switch res := resp.(type) {
	case []byte:
		result = res

	default:
		result, err = marshaler.MarshalJSON(resp)
		if CheckInternalServerError(w, err) {
			return
		}
	}

	wrappedResp := NewResponseWithHeight(ctx.Height, result)

	output, err := marshaler.MarshalJSON(wrappedResp)
	if CheckInternalServerError(w, err) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(output)
}

// CheckInternalServerError attaches an error message to an HTTP 500 INTERNAL SERVER ERROR response.
// Returns false when err is nil; it returns true otherwise.
func CheckInternalServerError(w http.ResponseWriter, err error) bool {
	return CheckError(w, http.StatusInternalServerError, err)
}

func CheckError(w http.ResponseWriter, status int, err error) bool {
	if err != nil {
		WriteErrorResponse(w, status, err.Error())
		return true
	}

	return false
}

// ResponseWithHeight defines a response object type that wraps an original
// response with a height.
type ResponseWithHeight struct {
	Height int64           `json:"height"`
	Result json.RawMessage `json:"result"`
}

// NewResponseWithHeight creates a new ResponseWithHeight instance
func NewResponseWithHeight(height int64, result json.RawMessage) ResponseWithHeight {
	return ResponseWithHeight{
		Height: height,
		Result: result,
	}
}
