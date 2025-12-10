package googleads

import (
	"client-runaway-zenoti/internal/config"
	"context"
	"fmt"
	"time"

	"github.com/shenzhencenter/google-ads-pb/services"
	"golang.org/x/oauth2/google"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func getKeyData() []byte {
	return []byte(config.Confs.GAjson)
}

const (
	scopeAdwords = "https://www.googleapis.com/auth/adwords"
	baseURL      = "https://googleads.googleapis.com/v18"

	developerToken  = "gda-62oYP_OW7r3h4A4HJw"
	loginCustomerID = "6309578268"
)

func GetAdSpendGoPkg(
	ctx context.Context,
	customerID string,
	startDateT time.Time,
	endDateT time.Time,
) (float64, error) {

	startDate := startDateT.Format("2006-01-02")
	endDate := endDateT.Format("2006-01-02")
	// 1) Exchange service account JSON for an OAuth2 access token (scope: adwords)
	creds, err := google.CredentialsFromJSON(ctx, getKeyData(), "https://www.googleapis.com/auth/adwords")
	if err != nil {
		return 0, fmt.Errorf("creds from json: %w", err)
	}
	tok, err := creds.TokenSource.Token()
	if err != nil {
		return 0, fmt.Errorf("token: %w", err)
	}

	// 2) gRPC connection to Google Ads API
	conn, err := grpc.NewClient(
		"googleads.googleapis.com:443",
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")),
	)
	if err != nil {
		return 0, fmt.Errorf("dial: %w", err)
	}
	defer conn.Close()

	// 3) Required headers
	mdPairs := []string{
		"authorization", "Bearer " + tok.AccessToken,
		"developer-token", developerToken,
	}
	if loginCustomerID != "" {
		mdPairs = append(mdPairs, "login-customer-id", loginCustomerID)
	}
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(mdPairs...))

	// 4) Run GAQL (Search). We’ll sum cost_micros client-side.
	query := fmt.Sprintf(`
SELECT
  metrics.cost_micros
FROM customer
WHERE segments.date BETWEEN '%s' AND '%s'`, startDate, endDate)

	req := &services.SearchGoogleAdsRequest{
		CustomerId: customerID,
		Query:      query,
		// Optional: PageSize: 10000,
	}

	svc := services.NewGoogleAdsServiceClient(conn)
	resp, err := svc.Search(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("search: %w", err)
	}

	var totalMicros int64
	for _, row := range resp.Results {
		if row.GetMetrics() != nil {
			totalMicros += row.GetMetrics().GetCostMicros()
		}
	}

	return float64(totalMicros) / 1e6, nil // micros → currency units
}
