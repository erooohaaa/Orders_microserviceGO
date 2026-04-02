package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type HTTPPaymentGateway struct {
	client  *http.Client
	baseURL string
}

func NewHTTPPaymentGateway(client *http.Client, baseURL string) *HTTPPaymentGateway {
	return &HTTPPaymentGateway{client: client, baseURL: baseURL}
}

func (g *HTTPPaymentGateway) Authorize(ctx context.Context, orderID string, amount int64) (string, string, error) {
	body, _ := json.Marshal(map[string]any{
		"order_id": orderID,
		"amount":   amount,
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.baseURL+"/payments", bytes.NewReader(body))
	if err != nil {
		return "", "", ErrGatewayUnavailable
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return "", "", ErrGatewayUnavailable
	}
	defer resp.Body.Close()

	var res struct {
		TransactionID string `json:"transaction_id"`
		Status        string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", "", fmt.Errorf("decode payment response: %w", err)
	}
	return res.TransactionID, res.Status, nil
}
