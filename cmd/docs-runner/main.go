package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/chain/definitions"
	"github.com/ChefBingbong/viem-go/client"
	"github.com/ChefBingbong/viem-go/client/transport"
	"github.com/ChefBingbong/viem-go/contracts/erc20"
	"github.com/ChefBingbong/viem-go/utils/unit"
)

type runReadContractRequest struct {
	RPCURL       string `json:"rpcUrl"`
	TokenAddress string `json:"tokenAddress"`
	UserAddress  string `json:"userAddress"`
}

type runReadContractResponse struct {
	RawBalance       string `json:"rawBalance"`
	FormattedBalance string `json:"formattedBalance"`
	Decimals         int    `json:"decimals"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func main() {
	// Railway (and many PaaS) provide a PORT env var.
	port := strings.TrimSpace(os.Getenv("PORT"))
	defaultAddr := ":8787"
	if port != "" {
		defaultAddr = ":" + port
	}

	addr := env("DOCS_RUNNER_ADDR", defaultAddr)
	allowOrigin := env("DOCS_RUNNER_ALLOW_ORIGIN", "*") // set to your docs origin in production

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("POST /run/read-contract", func(w http.ResponseWriter, r *http.Request) {
		// CORS
		applyCORS(w, r, allowOrigin)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		var req runReadContractRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid json body"})
			return
		}

		rpcURL := strings.TrimSpace(req.RPCURL)
		if rpcURL == "" {
			rpcURL = "https://eth.merkle.io"
		}
		if err := validateRPCURL(rpcURL); err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
			return
		}

		if !common.IsHexAddress(req.TokenAddress) {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid tokenAddress"})
			return
		}
		if !common.IsHexAddress(req.UserAddress) {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid userAddress"})
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 12*time.Second)
		defer cancel()

		publicClient, err := client.CreatePublicClient(client.PublicClientConfig{
			Chain:     &definitions.Mainnet,
			Transport: transport.HTTP(rpcURL),
		})

		if err != nil {
			writeJSON(w, http.StatusBadGateway, errorResponse{Error: err.Error()})
			return
		}

		balanceAny, err := publicClient.ReadContract(ctx, client.ReadContractOptions{
			Address:      common.HexToAddress(req.TokenAddress),
			ABI:          erc20.ContractABI,
			FunctionName: "balanceOf",
			Args:         []any{common.HexToAddress(req.UserAddress)},
		})
		if err != nil {
			writeJSON(w, http.StatusBadGateway, errorResponse{Error: err.Error()})
			return
		}

		balance := balanceAny.(*big.Int) // ERC20.balanceOf always returns uint256
		decimals := 6                    // USDC default; you can extend this endpoint to call decimals()
		formatted := unit.FormatUnits(balance, decimals)

		writeJSON(w, http.StatusOK, runReadContractResponse{
			RawBalance:       balance.String(),
			FormattedBalance: formatted,
			Decimals:         decimals,
		})
	})

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	fmt.Printf("docs-runner listening on %s\n", addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func env(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func applyCORS(w http.ResponseWriter, r *http.Request, allowOrigin string) {
	// allowOrigin can be:
	// - "*" (allow all)
	// - "https://a.com" (single origin)
	// - "https://a.com,https://b.com" (list)
	origin := strings.TrimSpace(r.Header.Get("Origin"))

	w.Header().Set("access-control-allow-methods", "POST, OPTIONS")
	w.Header().Set("access-control-allow-headers", "content-type")

	if allowOrigin == "*" {
		w.Header().Set("access-control-allow-origin", "*")
		return
	}

	// If no Origin header, don't set allow-origin.
	if origin == "" {
		return
	}

	for _, o := range strings.Split(allowOrigin, ",") {
		if strings.TrimSpace(o) == origin {
			w.Header().Set("access-control-allow-origin", origin)
			w.Header().Set("vary", "Origin")
			return
		}
	}
}

func validateRPCURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("invalid rpcUrl")
	}
	if u.Scheme != "https" && u.Scheme != "http" {
		return fmt.Errorf("rpcUrl must be http(s)")
	}
	host := u.Hostname()
	if host == "" {
		return fmt.Errorf("invalid rpcUrl host")
	}

	// Block obvious SSRF targets.
	if host == "localhost" {
		return fmt.Errorf("rpcUrl host not allowed")
	}
	if ip := net.ParseIP(host); ip != nil {
		if ip.IsLoopback() || ip.IsPrivate() {
			return fmt.Errorf("rpcUrl host not allowed")
		}
	}
	// If it resolves to private/loopback, also block.
	ips, _ := net.LookupIP(host)
	for _, ip := range ips {
		if ip.IsLoopback() || ip.IsPrivate() {
			return fmt.Errorf("rpcUrl host not allowed")
		}
	}
	return nil
}
