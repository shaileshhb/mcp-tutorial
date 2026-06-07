package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ConvertCurrencyInput is the input for the ConvertCurrencyTool
type ConvertCurrencyInput struct {
	FromCurrency string  `json:"from_currency" jsonschema:"currency codes are ISO 4217 format, e.g. USD, EUR, GBP, etc. the currency to convert from"`
	ToCurrency   string  `json:"to_currency" jsonschema:"currency codes are ISO 4217 format, e.g. USD, EUR, GBP, etc. the currency to convert to"`
	Amount       float64 `json:"amount" jsonschema:"the amount to convert"`
}

// ConvertCurrencyOutput is the output for the ConvertCurrencyTool
type ConvertCurrencyOutput struct {
	ConvertedAmount float64 `json:"converted_amount" jsonschema:"the converted amount"`
}

type CurrencyResponse struct {
	Date string                        `json:"date" jsonschema:"the date of the currency rates"`
	Data map[string]map[string]float64 `json:"-" jsonschema:"the currency rates"`
}

const BASE_CURRENCY_API_URL = `https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies`

// ConvertCurrencyTool converts a currency to another currency
func ConvertCurrencyTool(ctx context.Context, req *mcp.CallToolRequest, input ConvertCurrencyInput) (*mcp.CallToolResult, any, error) {
	input.FromCurrency = strings.ToLower(input.FromCurrency)
	input.ToCurrency = strings.ToLower(input.ToCurrency)
	fmt.Println("inside convert currency tool input", input)

	apiURL := fmt.Sprintf("%s/%s.json", BASE_CURRENCY_API_URL, input.FromCurrency)
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	response, err := UnmarshalJSON(body)
	if err != nil {
		return nil, nil, err
	}

	if _, ok := response.Data[input.FromCurrency]; !ok {
		return nil, nil, fmt.Errorf("currency %s not found", input.FromCurrency)
	}

	if _, ok := response.Data[input.FromCurrency][input.ToCurrency]; !ok {
		return nil, nil, fmt.Errorf("currency %s not found", input.ToCurrency)
	}

	convertedAmount := response.Data[input.FromCurrency][input.ToCurrency] * input.Amount
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("Converted amount: %f", convertedAmount),
			},
			// ToolResult: &mcp.ToolResult{
			// 	ToolName:   "currency-conversion",
			// 	ToolInput:  input,
			// 	ToolOutput: &ConvertCurrencyOutput{ConvertedAmount: convertedAmount},
			// },
		},
	}, nil, nil
}

func UnmarshalJSON(data []byte) (*CurrencyResponse, error) {
	var r CurrencyResponse
	var raw map[string]json.RawMessage

	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	// Extract date
	if v, ok := raw["date"]; ok {
		if err := json.Unmarshal(v, &r.Date); err != nil {
			return nil, err
		}
		delete(raw, "date")
	}

	// Remaining keys are currency codes
	r.Data = make(map[string]map[string]float64)

	for currency, v := range raw {
		var rates map[string]float64
		if err := json.Unmarshal(v, &rates); err != nil {
			return nil, err
		}
		r.Data[currency] = rates
	}

	return &r, nil
}
