package controller

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type MetricsCollector struct {
	// Customer Profile metrics
	customerProfileGetTotal    metric.Int64Counter
	customerProfileGetErrors   metric.Int64Counter
	customerProfileGetDuration metric.Float64Histogram

	// Supplier Profile metrics
	supplierProfileGetTotal    metric.Int64Counter
	supplierProfileGetErrors   metric.Int64Counter
	supplierProfileGetDuration metric.Float64Histogram

	// Dashboard metrics
	dashboardGetTotal    metric.Int64Counter
	dashboardGetErrors   metric.Int64Counter
	dashboardGetDuration metric.Float64Histogram

	// Login metrics
	loginTotal    metric.Int64Counter
	loginErrors   metric.Int64Counter
	loginDuration metric.Float64Histogram

	// Email Confirmation metrics
	emailConfirmationTotal    metric.Int64Counter
	emailConfirmationErrors   metric.Int64Counter
	emailConfirmationDuration metric.Float64Histogram

	// Password Forgot metrics
	passwordForgotTotal    metric.Int64Counter
	passwordForgotErrors   metric.Int64Counter
	passwordForgotDuration metric.Float64Histogram

	// Customer Signup metrics
	customerCreateTotal    metric.Int64Counter
	customerCreateErrors   metric.Int64Counter
	customerCreateDuration metric.Float64Histogram

	// Supplier Signup metrics
	supplierCreateTotal    metric.Int64Counter
	supplierCreateErrors   metric.Int64Counter
	supplierCreateDuration metric.Float64Histogram

	// Database operation metrics
	dbOperationsTotal   metric.Int64Counter
	dbOperationErrors   metric.Int64Counter
	dbOperationDuration metric.Float64Histogram
}

func NewMetricsCollector() *MetricsCollector {
	meter := otel.GetMeterProvider().Meter("megacommerce-user", metric.WithInstrumentationVersion("0.1.0"))

	mc := &MetricsCollector{}

	// Customer Profile metrics
	mc.customerProfileGetTotal, _ = meter.Int64Counter("customer_profile_get_total",
		metric.WithDescription("Total customer profile get requests"))
	mc.customerProfileGetErrors, _ = meter.Int64Counter("customer_profile_get_errors_total",
		metric.WithDescription("Total customer profile get errors"))
	mc.customerProfileGetDuration, _ = meter.Float64Histogram("customer_profile_get_duration_seconds",
		metric.WithDescription("Customer profile get request duration in seconds"))

	// Supplier Profile metrics
	mc.supplierProfileGetTotal, _ = meter.Int64Counter("supplier_profile_get_total",
		metric.WithDescription("Total supplier profile get requests"))
	mc.supplierProfileGetErrors, _ = meter.Int64Counter("supplier_profile_get_errors_total",
		metric.WithDescription("Total supplier profile get errors"))
	mc.supplierProfileGetDuration, _ = meter.Float64Histogram("supplier_profile_get_duration_seconds",
		metric.WithDescription("Supplier profile get request duration in seconds"))

	// Dashboard metrics
	mc.dashboardGetTotal, _ = meter.Int64Counter("dashboard_get_total",
		metric.WithDescription("Total dashboard get requests"))
	mc.dashboardGetErrors, _ = meter.Int64Counter("dashboard_get_errors_total",
		metric.WithDescription("Total dashboard get errors"))
	mc.dashboardGetDuration, _ = meter.Float64Histogram("dashboard_get_duration_seconds",
		metric.WithDescription("Dashboard get request duration in seconds"))

	// Login metrics
	mc.loginTotal, _ = meter.Int64Counter("login_total",
		metric.WithDescription("Total login requests"))
	mc.loginErrors, _ = meter.Int64Counter("login_errors_total",
		metric.WithDescription("Total login errors"))
	mc.loginDuration, _ = meter.Float64Histogram("login_duration_seconds",
		metric.WithDescription("Login request duration in seconds"))

	// Email Confirmation metrics
	mc.emailConfirmationTotal, _ = meter.Int64Counter("email_confirmation_total",
		metric.WithDescription("Total email confirmation requests"))
	mc.emailConfirmationErrors, _ = meter.Int64Counter("email_confirmation_errors_total",
		metric.WithDescription("Total email confirmation errors"))
	mc.emailConfirmationDuration, _ = meter.Float64Histogram("email_confirmation_duration_seconds",
		metric.WithDescription("Email confirmation request duration in seconds"))

	// Password Forgot metrics
	mc.passwordForgotTotal, _ = meter.Int64Counter("password_forgot_total",
		metric.WithDescription("Total password forgot requests"))
	mc.passwordForgotErrors, _ = meter.Int64Counter("password_forgot_errors_total",
		metric.WithDescription("Total password forgot errors"))
	mc.passwordForgotDuration, _ = meter.Float64Histogram("password_forgot_duration_seconds",
		metric.WithDescription("Password forgot request duration in seconds"))

	// Customer Create metrics
	mc.customerCreateTotal, _ = meter.Int64Counter("customer_create_total",
		metric.WithDescription("Total customer create requests"))
	mc.customerCreateErrors, _ = meter.Int64Counter("customer_create_errors_total",
		metric.WithDescription("Total customer create errors"))
	mc.customerCreateDuration, _ = meter.Float64Histogram("customer_create_duration_seconds",
		metric.WithDescription("Customer create request duration in seconds"))

	// Supplier Create metrics
	mc.supplierCreateTotal, _ = meter.Int64Counter("supplier_create_total",
		metric.WithDescription("Total supplier create requests"))
	mc.supplierCreateErrors, _ = meter.Int64Counter("supplier_create_errors_total",
		metric.WithDescription("Total supplier create errors"))
	mc.supplierCreateDuration, _ = meter.Float64Histogram("supplier_create_duration_seconds",
		metric.WithDescription("Supplier create request duration in seconds"))

	// Database operation metrics
	mc.dbOperationsTotal, _ = meter.Int64Counter("db_operations_total",
		metric.WithDescription("Total database operations"))
	mc.dbOperationErrors, _ = meter.Int64Counter("db_operation_errors_total",
		metric.WithDescription("Total database operation errors"))
	mc.dbOperationDuration, _ = meter.Float64Histogram("db_operation_duration_seconds",
		metric.WithDescription("Database operation duration in seconds"))

	return mc
}

func (m *MetricsCollector) RecordCustomerProfileGetRequest(success bool, duration float64) {
	ctx := context.Background()
	m.customerProfileGetTotal.Add(ctx, 1)
	m.customerProfileGetDuration.Record(ctx, duration)
	if !success {
		m.customerProfileGetErrors.Add(ctx, 1)
	}
}

func (m *MetricsCollector) RecordSupplierProfileGetRequest(success bool, duration float64) {
	ctx := context.Background()
	m.supplierProfileGetTotal.Add(ctx, 1)
	m.supplierProfileGetDuration.Record(ctx, duration)
	if !success {
		m.supplierProfileGetErrors.Add(ctx, 1)
	}
}

func (m *MetricsCollector) RecordDashboardGetRequest(success bool, duration float64) {
	ctx := context.Background()
	m.dashboardGetTotal.Add(ctx, 1)
	m.dashboardGetDuration.Record(ctx, duration)
	if !success {
		m.dashboardGetErrors.Add(ctx, 1)
	}
}

func (m *MetricsCollector) RecordLoginRequest(success bool, duration float64) {
	ctx := context.Background()
	m.loginTotal.Add(ctx, 1)
	m.loginDuration.Record(ctx, duration)
	if !success {
		m.loginErrors.Add(ctx, 1)
	}
}

func (m *MetricsCollector) RecordEmailConfirmationRequest(success bool, duration float64) {
	ctx := context.Background()
	m.emailConfirmationTotal.Add(ctx, 1)
	m.emailConfirmationDuration.Record(ctx, duration)
	if !success {
		m.emailConfirmationErrors.Add(ctx, 1)
	}
}

func (m *MetricsCollector) RecordPasswordForgotRequest(success bool, duration float64) {
	ctx := context.Background()
	m.passwordForgotTotal.Add(ctx, 1)
	m.passwordForgotDuration.Record(ctx, duration)
	if !success {
		m.passwordForgotErrors.Add(ctx, 1)
	}
}

func (m *MetricsCollector) RecordCustomerCreateRequest(success bool, duration float64) {
	ctx := context.Background()
	m.customerCreateTotal.Add(ctx, 1)
	m.customerCreateDuration.Record(ctx, duration)
	if !success {
		m.customerCreateErrors.Add(ctx, 1)
	}
}

func (m *MetricsCollector) RecordSupplierCreateRequest(success bool, duration float64) {
	ctx := context.Background()
	m.supplierCreateTotal.Add(ctx, 1)
	m.supplierCreateDuration.Record(ctx, duration)
	if !success {
		m.supplierCreateErrors.Add(ctx, 1)
	}
}

func (m *MetricsCollector) RecordDBOperation(success bool, duration float64) {
	ctx := context.Background()
	m.dbOperationsTotal.Add(ctx, 1)
	m.dbOperationDuration.Record(ctx, duration)
	if !success {
		m.dbOperationErrors.Add(ctx, 1)
	}
}
