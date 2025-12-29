package svc_googleads

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	"client-runaway-zenoti/packages/googleads"
	"context"
	"strconv"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

// ListAccountHierarchy returns the full account hierarchy accessible by the connection without using login-customer-id.
func ListAccountHierarchy(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	id, err := strconv.ParseUint(c.Param("accountId"), 10, 64)
	lvn.GinErr(c, 400, err, "invalid account id")
	if err != nil {
		return
	}

	connModel := models.GoogleAdsConnection{}
	err = db.DB.Where("profile_id = ? AND id = ?", user.ProfileID, id).First(&connModel).Error
	lvn.GinErr(c, 400, err, "connection not found")
	if err != nil {
		return
	}

	conn := fromModel(connModel)

	ctx := context.Background()
	ctx, _, err = Svc.WithHeaders(ctx, conn)
	lvn.GinErr(c, 400, err, "unable to refresh tokens")
	if err != nil {
		return
	}
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		mdCopy := md.Copy()
		delete(mdCopy, "login-customer-id")
		ctx = metadata.NewOutgoingContext(ctx, mdCopy)
	}

	grpcConn, err := grpc.Dial("googleads.googleapis.com:443", grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	lvn.GinErr(c, 400, err, "unable to connect to google ads")
	if err != nil {
		return
	}
	defer grpcConn.Close()

	accessibleIDs, err := googleads.ListAccessibleCustomers(ctx, grpcConn)
	lvn.GinErr(c, 400, err, "unable to list accessible customers")
	if err != nil {
		return
	}

	roots, err := googleads.GetAccountHierarchy(ctx, grpcConn, accessibleIDs)
	lvn.GinErr(c, 400, err, "unable to build account hierarchy")
	if err != nil {
		return
	}

	var data []gin.H
	var build func(node *googleads.AccountNode) gin.H
	build = func(node *googleads.AccountNode) gin.H {
		entry := gin.H{
			"name":       node.DescriptiveName,
			"customerId": node.CustomerID,
			"isManager":  node.Manager,
		}
		var clients []gin.H
		for _, ch := range node.DirectClients {
			clients = append(clients, build(ch))
		}
		entry["directClients"] = clients
		return entry
	}

	for _, r := range roots {
		data = append(data, build(r))
	}

	c.Data(lvn.Res(200, data, "success"))
}
