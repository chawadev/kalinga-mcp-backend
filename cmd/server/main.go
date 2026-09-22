package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/chawadev/kalinga-backend/internal/api"
	"github.com/chawadev/kalinga-backend/internal/auth"
	"github.com/chawadev/kalinga-backend/internal/database"
	"github.com/joho/godotenv"
)

type JSONRPCRequest struct {
	Jsonrpc string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	} `json:"params"`
}

type JSONRPCResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *RPCError) Error() string {
	return e.Message
}

func handleCheckStatus(args map[string]interface{}) interface{} {
	return map[string]interface{}{
		"status":  "ok",
		"service": "kalinga",
		"version": "0.1.0",
		"health":  "healthy",
	}
}

func handleConnectMomoAccount(args map[string]interface{}) interface{} {
	phoneNumber, _ := args["phone_number"].(string)
	provider, _ := args["provider"].(string)

	return map[string]interface{}{
		"status":         "approved",
		"phone_number":   phoneNumber,
		"provider":       provider,
		"account_linked": true,
		"message":        "Mobile money account successfully connected",
	}
}

func handleBuildFinancialProfile(args map[string]interface{}) interface{} {
	userID, _ := args["user_id"].(string)

	return map[string]interface{}{
		"user_id":              userID,
		"credit_score":         75,
		"loan_readiness":       "high",
		"income_regularity":    85,
		"volatility":          12,
		"repayment_capacity":   92,
		"expense_to_income":    68,
		"projected_savings":    "+$420/mo",
		"financial_wellness":   78,
		"risk_assessment":      "low",
		"recommended_loan_amt": 5000,
	}
}

func handleVerifyClaim(args map[string]interface{}) interface{} {
	claimID, _ := args["claim_id"].(string)
	amount, _ := args["amount"].(float64)
	merchant, _ := args["merchant"].(string)

	isSuspicious := amount > 1000 && merchant == "unknown"

	result := map[string]interface{}{
		"claim_id":      claimID,
		"amount":        amount,
		"merchant":      merchant,
		"status":        "verified",
		"risk_level":    "low",
		"flagged":       isSuspicious,
		"pattern_match": "regular_merchant",
		"confidence":    0.95,
	}

	if isSuspicious {
		result["status"] = "flagged"
		result["risk_level"] = "high"
		result["pattern_match"] = "suspicious_pattern"
		result["confidence"] = 0.72
	}

	return result
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Disconnect()

	// Create database indexes (ignore if already exists)
	if err := database.CreateIndexes(database.Database); err != nil {
		log.Println("Warning: Failed to create indexes:", err)
	}

	// Initialize services
	authRepo := auth.NewRepository(database.Database)
	authService := auth.NewService(authRepo)

	// Setup router
	mux := http.NewServeMux()
	router := api.NewRouter(authService)
	router.SetupRoutes(mux)

	// MCP endpoint
	mux.HandleFunc("/mcp", mcpHandler)

	log.Println("Kalinga server listening on http://localhost:8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func mcpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var rpcReq JSONRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&rpcReq); err != nil {
		json.NewEncoder(w).Encode(JSONRPCResponse{
			Jsonrpc: "2.0",
			ID:      nil,
			Error:   &RPCError{Code: -32700, Message: "Parse error"},
		})
		return
	}

	// Handle tools/call method
	if rpcReq.Method == "tools/call" {
		var result interface{}
		var err error

		switch rpcReq.Params.Name {
		case "check_status":
			result = handleCheckStatus(rpcReq.Params.Arguments)
		case "connect_momo_account":
			result = handleConnectMomoAccount(rpcReq.Params.Arguments)
		case "build_financial_profile":
			result = handleBuildFinancialProfile(rpcReq.Params.Arguments)
		case "verify_claim":
			result = handleVerifyClaim(rpcReq.Params.Arguments)
		default:
			err = &RPCError{Code: -32601, Message: "Method not found"}
		}

		if err != nil {
			json.NewEncoder(w).Encode(JSONRPCResponse{
				Jsonrpc: "2.0",
				ID:      rpcReq.ID,
				Error:   err.(*RPCError),
			})
			return
		}

		json.NewEncoder(w).Encode(JSONRPCResponse{
			Jsonrpc: "2.0",
			ID:      rpcReq.ID,
			Result:  result,
		})
	} else {
		json.NewEncoder(w).Encode(JSONRPCResponse{
			Jsonrpc: "2.0",
			ID:      rpcReq.ID,
			Error:   &RPCError{Code: -32601, Message: "Method not found"},
		})
	}
}
