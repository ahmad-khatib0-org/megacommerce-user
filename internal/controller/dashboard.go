package controller

import (
	"context"
	"time"

	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/users/v1"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/models"
	intModels "github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"google.golang.org/grpc/codes"
)

func (c *Controller) GetSupplierDashboard(context context.Context, req *pb.DashboardRequest) (*pb.DashboardResponse, error) {
	start := time.Now()
	path := "user.controller.GetSupplierDashboard"
	errBuilder := func(e *models.AppError) (*pb.DashboardResponse, error) {
		return &pb.DashboardResponse{Response: &pb.DashboardResponse_Error{Error: models.AppErrorToProto(e)}}, nil
	}

	ctx, err := models.ContextGet(context)
	if err != nil {
		duration := time.Since(start).Seconds()
		c.metricsCollector.RecordDashboardGetRequest(false, duration)
		return errBuilder(err)
	}

	ar := models.AuditRecordNew(ctx, intModels.EventNameDashboardGet, models.EventStatusFail)
	defer c.ProcessAudit(ar)

	// Get user ID from context session
	userID := ctx.Session.UserID
	if userID == "" {
		duration := time.Since(start).Seconds()
		c.metricsCollector.RecordDashboardGetRequest(false, duration)
		return errBuilder(models.NewAppError(ctx, path, "error.unauthenticated", nil, "user not authenticated", int(codes.Unauthenticated), nil))
	}

	user, dbErr := c.store.UsersGetByID(ctx, userID)
	if dbErr != nil {
		duration := time.Since(start).Seconds()
		c.metricsCollector.RecordDashboardGetRequest(false, duration)
		if dbErr.ErrType == models.DBErrorTypeNoRows {
			return errBuilder(models.NewAppError(ctx, path, "error.not_found", nil, "supplier not found", int(codes.NotFound), nil))
		}
		return errBuilder(models.NewAppError(ctx, path, models.ErrMsgInternal, nil, dbErr.Details, int(codes.Internal), &models.AppErrorErrorsArgs{Err: dbErr}))
	}

	// Verify user is a supplier
	if !isSupplier(user) {
		duration := time.Since(start).Seconds()
		c.metricsCollector.RecordDashboardGetRequest(false, duration)
		return errBuilder(models.NewAppError(ctx, path, "error.permission_denied", nil, "user is not a supplier", int(codes.PermissionDenied), nil))
	}

	ar.Success()

	// Build dashboard stats
	// TODO: Integrate with products and inventory services to fetch real data
	_ = &pb.DashboardStats{
		TotalProducts:       0, // TODO: fetch from products service
		TotalInventoryItems: 0, // TODO: fetch from inventory service
		TotalReviews:        0, // TODO: fetch from products service
		ProductVisitsCount:  0, // TODO: fetch from analytics/products service
		VisitsByPeriod: &pb.VisitsByPeriod{
			Today:     0, // TODO: calculate from analytics
			Yesterday: 0, // TODO: calculate from analytics
			LastWeek:  0, // TODO: calculate from analytics
			LastMonth: 0, // TODO: calculate from analytics
			LastYear:  0, // TODO: calculate from analytics
		},
		PendingOrders: 0, // TODO: fetch from orders service
		TotalOrders:   0, // TODO: fetch from orders service
	}

	// Generate random data for now
	mockStats := generateMockDashboardStats()

	duration := time.Since(start).Seconds()
	c.metricsCollector.RecordDashboardGetRequest(true, duration)

	return &pb.DashboardResponse{Response: &pb.DashboardResponse_Data{Data: mockStats}}, nil
}

// Helper function to check if user is a supplier
func isSupplier(user *pb.User) bool {
	for _, role := range user.Roles {
		if role == string(models.RoleIDSupplierAdmin) || role == "supplier" {
			return true
		}
	}
	return false
}

// Generate mock dashboard stats with random data for now
func generateMockDashboardStats() *pb.DashboardStats {
	return &pb.DashboardStats{
		TotalProducts:       42,
		TotalInventoryItems: 156,
		TotalReviews:        28,
		ProductVisitsCount:  1847,
		VisitsByPeriod: &pb.VisitsByPeriod{
			Today:     245,
			Yesterday: 189,
			LastWeek:  1234,
			LastMonth: 4567,
			LastYear:  28945,
		},
		PendingOrders: 5,
		TotalOrders:   87,
	}
}
