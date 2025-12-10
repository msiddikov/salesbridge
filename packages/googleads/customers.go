package googleads

import (
	"context"
	"fmt"
	"strings"

	"github.com/shenzhencenter/google-ads-pb/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// CustomerSummary is a lightweight view of accessible customer accounts.
type CustomerSummary struct {
	CustomerID      string `json:"customerId"`
	DescriptiveName string `json:"descriptiveName"`
	Manager         bool   `json:"manager"`
	Hidden          bool   `json:"hidden"`
}

// AccountNode represents a customer and its direct clients in a hierarchy.
type AccountNode struct {
	CustomerSummary
	DirectClients []*AccountNode `json:"directClients,omitempty"`
}

// ListAccounts returns the specified account and, if it's a manager, its direct clients.
// ctx must already contain the required Google Ads headers (use Service.WithHeaders first).
func ListAccounts(ctx context.Context, conn *grpc.ClientConn, customerID string) ([]CustomerSummary, error) {
	if customerID == "" {
		return nil, fmt.Errorf("googleads: customer id required")
	}

	self, err := fetchSelfAccount(ctx, conn, customerID)
	if err != nil {
		return nil, err
	}

	accounts := []CustomerSummary{self}
	if !self.Manager {
		return accounts, nil
	}

	clients, err := fetchDirectClients(ctx, conn, customerID)
	if err != nil {
		return nil, err
	}
	accounts = append(accounts, clients...)
	return accounts, nil
}

// ListAccessibleCustomers returns all customer IDs the authenticated user can access.
// ctx must already include the Google Ads auth headers (use Service.WithHeaders).
func ListAccessibleCustomers(ctx context.Context, conn *grpc.ClientConn) ([]string, error) {
	// Removing login-customer-id ensures we enumerate every account the user can access
	// (per Google docs, the header scopes results; omit it to get the full list).
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		mdCopy := md.Copy()
		delete(mdCopy, "login-customer-id")
		ctx = metadata.NewOutgoingContext(ctx, mdCopy)
	}

	svc := services.NewCustomerServiceClient(conn)
	resp, err := svc.ListAccessibleCustomers(ctx, &services.ListAccessibleCustomersRequest{})
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(resp.ResourceNames))
	for _, rn := range resp.ResourceNames {
		if rn == "" {
			continue
		}
		if idx := strings.LastIndex(rn, "/"); idx >= 0 && idx < len(rn)-1 {
			ids = append(ids, rn[idx+1:])
			continue
		}
		ids = append(ids, rn)
	}
	return ids, nil
}

// GetCustomerSummary returns basic info for a single customer ID.
func GetCustomerSummary(ctx context.Context, conn *grpc.ClientConn, customerID string) (CustomerSummary, error) {
	return fetchSelfAccount(ctx, conn, customerID)
}

// GetAccountHierarchy builds the full hierarchy reachable from the provided root IDs.
// It does not rely on login-customer-id; instead, it iterates through accessible IDs.
func GetAccountHierarchy(ctx context.Context, conn *grpc.ClientConn, rootIDs []string) ([]*AccountNode, error) {
	type nodeRef struct {
		ptr *AccountNode
	}
	seen := map[string]bool{}
	parent := map[string]string{} // child -> parent
	nodes := map[string]*AccountNode{}
	queue := []string{}

	for _, id := range rootIDs {
		if id == "" || seen[id] {
			continue
		}
		seen[id] = true
		queue = append(queue, id)
	}

	for len(queue) > 0 {
		cid := queue[0]
		queue = queue[1:]

		self, err := fetchSelfAccount(ctx, conn, cid)
		if err != nil {
			continue
		}

		cur := nodes[cid]
		if cur == nil {
			cur = &AccountNode{}
			nodes[cid] = cur
		}
		cur.CustomerSummary = self

		if !self.Manager {
			continue
		}

		children, err := fetchDirectClients(ctx, conn, cid)
		if err != nil {
			return nil, err
		}
		for _, ch := range children {
			childNode := nodes[ch.CustomerID]
			if childNode == nil {
				childNode = &AccountNode{}
				nodes[ch.CustomerID] = childNode
			}
			childNode.CustomerSummary = ch
			childNodeID := ch.CustomerID
			if childNodeID != "" && parent[childNodeID] == "" {
				parent[childNodeID] = cid
			}
			cur.DirectClients = append(cur.DirectClients, childNode)
			if ch.Manager && !seen[childNodeID] {
				seen[childNodeID] = true
				queue = append(queue, childNodeID)
			}
		}
	}

	var roots []*AccountNode
	for _, id := range rootIDs {
		if id == "" {
			continue
		}
		n := nodes[id]
		if n == nil {
			continue
		}
		if parent[id] == "" {
			roots = append(roots, n)
		}
	}
	return roots, nil
}

// ListAccountHierarchy returns the manager account and its direct clients using the documented GAQL query.
// It aligns with https://developers.google.com/google-ads/api/docs/account-management/listing-accounts
// by querying customer_client for level-1 relationships under the given manager.
func ListAccountHierarchy(ctx context.Context, conn *grpc.ClientConn, managerID string) ([]CustomerSummary, error) {
	if managerID == "" {
		return nil, fmt.Errorf("googleads: manager id required")
	}

	// Ensure login-customer-id is set to the manager for this query.
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		mdCopy := md.Copy()
		mdCopy.Set("login-customer-id", managerID)
		ctx = metadata.NewOutgoingContext(ctx, mdCopy)
	} else {
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("login-customer-id", managerID))
	}

	// Include the manager itself.
	self, err := fetchSelfAccount(ctx, conn, managerID)
	if err != nil {
		return nil, err
	}

	query := `
SELECT
  customer_client.client_customer,
  customer_client.descriptive_name,
  customer_client.manager,
  customer_client.hidden
FROM customer_client
WHERE customer_client.level = 1
`
	req := &services.SearchGoogleAdsRequest{
		CustomerId: managerID,
		Query:      query,
	}

	svc := services.NewGoogleAdsServiceClient(conn)
	resp, err := svc.Search(ctx, req)
	if err != nil {
		return nil, err
	}

	out := []CustomerSummary{self}
	for _, row := range resp.Results {
		cc := row.GetCustomerClient()
		if cc == nil {
			continue
		}
		out = append(out, CustomerSummary{
			CustomerID:      strings.ReplaceAll(cc.GetClientCustomer(), "-", ""),
			DescriptiveName: cc.GetDescriptiveName(),
			Manager:         cc.GetManager(),
			Hidden:          cc.GetHidden(),
		})
	}
	return out, nil
}

func fetchSelfAccount(ctx context.Context, conn *grpc.ClientConn, customerID string) (CustomerSummary, error) {
	query := `
SELECT
  customer.id,
  customer.descriptive_name,
  customer.manager
FROM customer
LIMIT 1
`
	req := &services.SearchGoogleAdsRequest{
		Query: query,
	}

	if customerID != "" {
		req.CustomerId = customerID
	}

	svc := services.NewGoogleAdsServiceClient(conn)
	resp, err := svc.Search(ctx, req)
	if err != nil {
		return CustomerSummary{}, err
	}
	if len(resp.Results) == 0 || resp.Results[0].GetCustomer() == nil {
		return CustomerSummary{}, fmt.Errorf("googleads: unable to fetch customer %s", customerID)
	}
	cust := resp.Results[0].GetCustomer()
	return CustomerSummary{
		CustomerID:      strings.ReplaceAll(fmt.Sprint(cust.GetId()), "-", ""),
		DescriptiveName: cust.GetDescriptiveName(),
		Manager:         cust.GetManager(),
	}, nil
}

func fetchDirectClients(ctx context.Context, conn *grpc.ClientConn, customerID string) ([]CustomerSummary, error) {
	query := `
SELECT
  customer_client.client_customer,
  customer_client.descriptive_name,
  customer_client.manager,
  customer_client.hidden
FROM customer_client
WHERE customer_client.level = 1
`
	req := &services.SearchGoogleAdsRequest{
		CustomerId: customerID,
		Query:      query,
	}

	svc := services.NewGoogleAdsServiceClient(conn)
	resp, err := svc.Search(ctx, req)
	if err != nil {
		return nil, err
	}

	var out []CustomerSummary
	for _, row := range resp.Results {
		cc := row.GetCustomerClient()
		if cc == nil {
			continue
		}
		out = append(out, CustomerSummary{
			CustomerID:      strings.ReplaceAll(cc.GetClientCustomer(), "-", ""),
			DescriptiveName: cc.GetDescriptiveName(),
			Manager:         cc.GetManager(),
			Hidden:          cc.GetHidden(),
		})
	}
	return out, nil
}
