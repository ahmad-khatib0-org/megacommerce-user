package models

// Permission represents a single permission with metadata
type Permission struct {
	ID          string
	Name        string
	Description string // Optional extended description
	Category    string // e.g., "account", "orders", etc.
}

// String implements fmt.Stringer for easy printing
func (p *Permission) String() string {
	return p.Name
}

// Permission registry

var (
	// Normal role permissions
	PermissionProfileView = &Permission{
		ID:          "profile_view",
		Name:        "authentication.permissions.profile_view.name",
		Description: "authentication.permissions.profile_view.description",
		Category:    "account",
	}
	PermissionProfileEdit = &Permission{
		ID:          "profile_edit",
		Name:        "authentication.permissions.profile_edit.name",
		Description: "authentication.permissions.profile_edit.description",
		Category:    "account",
	}
	PermissionPasswordUpdate = &Permission{
		ID:          "password_update",
		Name:        "authentication.permissions.password_update.name",
		Description: "authentication.permissions.password_update.description",
		Category:    "account",
	}
	PermissionAccountDelete = &Permission{
		ID:          "account_delete",
		Name:        "authentication.permissions.account_delete.name",
		Description: "authentication.permissions.account_delete.description",
		Category:    "account",
	}
	PermissionPreferencesSet = &Permission{
		ID:          "preferences_set",
		Name:        "authentication.permissions.preferences_set.name",
		Description: "authentication.permissions.preferences_set.description",
		Category:    "account",
	}

	PermissionOrderPlace = &Permission{
		ID:          "order_place",
		Name:        "authentication.permissions.order_place.name",
		Description: "authentication.permissions.order_place.description",
		Category:    "orders",
	}
	PermissionOrderCancel = &Permission{
		ID:          "order_cancel",
		Name:        "authentication.permissions.order_cancel.name",
		Description: "authentication.permissions.order_cancel.description",
		Category:    "orders",
	}
	PermissionOrderTrack = &Permission{
		ID:          "order_track",
		Name:        "authentication.permissions.order_track.name",
		Description: "authentication.permissions.order_track.description",
		Category:    "orders",
	}
	PermissionOrderHistoryView = &Permission{
		ID:          "order_history_view",
		Name:        "authentication.permissions.order_history_view.name",
		Description: "authentication.permissions.order_history_view.description",
		Category:    "orders",
	}

	PermissionCardAdd = &Permission{
		ID:          "card_add",
		Name:        "authentication.permissions.card_add.name",
		Description: "authentication.permissions.card_add.description",
		Category:    "payments",
	}
	PermissionCardRemove = &Permission{
		ID:          "card_remove",
		Name:        "authentication.permissions.card_remove.name",
		Description: "authentication.permissions.card_remove.description",
		Category:    "payments",
	}
	PermissionTransactionsView = &Permission{
		ID:          "transactions_view",
		Name:        "authentication.permissions.transactions_view.name",
		Description: "authentication.permissions.transactions_view.description",
		Category:    "payments",
	}
	PermissionCouponsApply = &Permission{
		ID:          "coupons_apply",
		Name:        "authentication.permissions.coupons_apply.name",
		Description: "authentication.permissions.coupons_apply.description",
		Category:    "payments",
	}
	PermissionWalletSave = &Permission{
		ID:          "wallet_save",
		Name:        "authentication.permissions.wallet_save.name",
		Description: "authentication.permissions.wallet_save.description",
		Category:    "payments",
	}

	PermissionReviewWrite = &Permission{
		ID:          "review_write",
		Name:        "authentication.permissions.review_write.name",
		Description: "authentication.permissions.review_write.description",
		Category:    "reviews",
	}
	PermissionProductRate = &Permission{
		ID:          "product_rate",
		Name:        "authentication.permissions.product_rate.name",
		Description: "authentication.permissions.product_rate.description",
		Category:    "reviews",
	}

	PermissionWishlistAdd = &Permission{
		ID:          "wishlist_add",
		Name:        "authentication.permissions.wishlist_add.name",
		Description: "authentication.permissions.wishlist_add.description",
		Category:    "wishlist",
	}
	PermissionWishlistRemove = &Permission{
		ID:          "wishlist_remove",
		Name:        "authentication.permissions.wishlist_remove.name",
		Description: "authentication.permissions.wishlist_remove.description",
		Category:    "wishlist",
	}
	PermissionWishlistView = &Permission{
		ID:          "wishlist_view",
		Name:        "authentication.permissions.wishlist_view.name",
		Description: "authentication.permissions.wishlist_view.description",
		Category:    "wishlist",
	}
	PermissionSuppliersFollow = &Permission{
		ID:          "suppliers_follow",
		Name:        "authentication.permissions.suppliers_follow.name",
		Description: "authentication.permissions.suppliers_follow.description",
		Category:    "wishlist",
	}

	PermissionTicketCreate = &Permission{
		ID:          "ticket_create",
		Name:        "authentication.permissions.ticket_create.name",
		Description: "authentication.permissions.ticket_create.description",
		Category:    "support",
	}
	PermissionTicketHistoryView = &Permission{
		ID:          "ticket_history_view",
		Name:        "authentication.permissions.ticket_history_view.name",
		Description: "authentication.permissions.ticket_history_view.description",
		Category:    "support",
	}
	PermissionTicketClose = &Permission{
		ID:          "ticket_close",
		Name:        "authentication.permissions.ticket_close.name",
		Description: "authentication.permissions.ticket_close.description",
		Category:    "support",
	}

	// Pro role permissions
	PermissionOrderReturn = &Permission{
		ID:          "order_return",
		Name:        "authentication.permissions.order_return.name",
		Description: "authentication.permissions.order_return.description",
		Category:    "orders",
	}
	PermissionReturnPriority = &Permission{
		ID:          "return_priority",
		Name:        "authentication.permissions.return_priority.name",
		Description: "authentication.permissions.return_priority.description",
		Category:    "orders",
	}
	PermissionDeliveriesScheduled = &Permission{
		ID:          "deliveries_scheduled",
		Name:        "authentication.permissions.deliveries_scheduled.name",
		Description: "authentication.permissions.deliveries_scheduled.description",
		Category:    "orders",
	}

	PermissionRewardsCashback = &Permission{
		ID:          "rewards_cashback",
		Name:        "authentication.permissions.rewards_cashback.name",
		Description: "authentication.permissions.rewards_cashback.description",
		Category:    "payments",
	}
	PermissionCheckoutOneClick = &Permission{
		ID:          "checkout_one_click",
		Name:        "authentication.permissions.checkout_one_click.name",
		Description: "authentication.permissions.checkout_one_click.description",
		Category:    "payments",
	}

	PermissionReviewEdit = &Permission{
		ID:          "review_edit",
		Name:        "authentication.permissions.review_edit.name",
		Description: "authentication.permissions.review_edit.description",
		Category:    "reviews",
	}
	PermissionReviewDelete = &Permission{
		ID:          "review_delete",
		Name:        "authentication.permissions.review_delete.name",
		Description: "authentication.permissions.review_delete.description",
		Category:    "reviews",
	}
	PermissionReviewReport = &Permission{
		ID:          "review_report",
		Name:        "authentication.permissions.review_report.name",
		Description: "authentication.permissions.review_report.description",
		Category:    "reviews",
	}
	PermissionReviewerPowerBadge = &Permission{
		ID:          "reviewer_power_badge",
		Name:        "authentication.permissions.reviewer_power_badge.name",
		Description: "authentication.permissions.reviewer_power_badge.description",
		Category:    "reviews",
	}

	PermissionWishlistShare = &Permission{
		ID:          "wishlist_share",
		Name:        "authentication.permissions.wishlist_share.name",
		Description: "authentication.permissions.wishlist_share.description",
		Category:    "wishlist",
	}
	PermissionTagsFollow = &Permission{
		ID:          "tags_follow",
		Name:        "authentication.permissions.tags_follow.name",
		Description: "authentication.permissions.tags_follow.description",
		Category:    "wishlist",
	}
	PermissionNotificationsRestock = &Permission{
		ID:          "notifications_restock",
		Name:        "authentication.permissions.notifications_restock.name",
		Description: "authentication.permissions.notifications_restock.description",
		Category:    "wishlist",
	}

	PermissionAgentChat = &Permission{
		ID:          "agent_chat",
		Name:        "authentication.permissions.agent_chat.name",
		Description: "authentication.permissions.agent_chat.description",
		Category:    "support",
	}
	PermissionSupportPriority = &Permission{
		ID:          "support_priority",
		Name:        "authentication.permissions.support_priority.name",
		Description: "authentication.permissions.support_priority.description",
		Category:    "support",
	}

	PermissionAlertsPriceDrop = &Permission{
		ID:          "alerts_price_drop",
		Name:        "authentication.permissions.alerts_price_drop.name",
		Description: "authentication.permissions.alerts_price_drop.description",
		Category:    "market_awareness",
	}
	PermissionAlertsRestock = &Permission{
		ID:          "alerts_restock",
		Name:        "authentication.permissions.alerts_restock.name",
		Description: "authentication.permissions.alerts_restock.description",
		Category:    "market_awareness",
	}

	// Org role permissions
	PermissionOrderingBulk = &Permission{
		ID:          "ordering_bulk",
		Name:        "authentication.permissions.ordering_bulk.name",
		Description: "authentication.permissions.ordering_bulk.description",
		Category:    "orders",
	}
	PermissionOrderApprovalWorkflows = &Permission{
		ID:          "order_approval_workflows",
		Name:        "authentication.permissions.order_approval_workflows.name",
		Description: "authentication.permissions.order_approval_workflows.description",
		Category:    "orders",
	}

	PermissionPaymentsInvoice = &Permission{
		ID:          "payments_invoice",
		Name:        "authentication.permissions.payments_invoice.name",
		Description: "authentication.permissions.payments_invoice.description",
		Category:    "payments",
	}
	PermissionTermsNet = &Permission{
		ID:          "terms_net",
		Name:        "authentication.permissions.terms_net.name",
		Description: "authentication.permissions.terms_net.description",
		Category:    "payments",
	}
	PermissionPricingCustom = &Permission{
		ID:          "pricing_custom",
		Name:        "authentication.permissions.pricing_custom.name",
		Description: "authentication.permissions.pricing_custom.description",
		Category:    "payments",
	}

	PermissionTrendsMarketRealtime = &Permission{
		ID:          "trends_market_realtime",
		Name:        "authentication.permissions.trends_market_realtime.name",
		Description: "authentication.permissions.trends_market_realtime.description",
		Category:    "market_awareness",
	}
	PermissionReportsCategoryPerformance = &Permission{
		ID:          "reports_category_performance",
		Name:        "authentication.permissions.reports_category_performance.name",
		Description: "authentication.permissions.reports_category_performance.description",
		Category:    "market_awareness",
	}
	PermissionTrackingProductLifecycle = &Permission{
		ID:          "tracking_product_lifecycle",
		Name:        "authentication.permissions.tracking_product_lifecycle.name",
		Description: "authentication.permissions.tracking_product_lifecycle.description",
		Category:    "market_awareness",
	}
	PermissionTrackingCompetitorProduct = &Permission{
		ID:          "tracking_competitor_product",
		Name:        "authentication.permissions.tracking_competitor_product.name",
		Description: "authentication.permissions.tracking_competitor_product.description",
		Category:    "market_awareness",
	}

	PermissionAccountsMultiUser = &Permission{
		ID:          "accounts_multi_user",
		Name:        "authentication.permissions.accounts_multi_user.name",
		Description: "authentication.permissions.accounts_multi_user.description",
		Category:    "org_tools",
	}
	PermissionAccessRoleBasedControl = &Permission{
		ID:          "access_role_based_control",
		Name:        "authentication.permissions.access_role_based_control.name",
		Description: "authentication.permissions.access_role_based_control.description",
		Category:    "org_tools",
	}
	PermissionDashboardOrg = &Permission{
		ID:          "dashboard_org",
		Name:        "authentication.permissions.dashboard_org.name",
		Description: "authentication.permissions.dashboard_org.description",
		Category:    "org_tools",
	}
	PermissionLogsAudit = &Permission{
		ID:          "logs_audit",
		Name:        "authentication.permissions.logs_audit.name",
		Description: "authentication.permissions.logs_audit.description",
		Category:    "org_tools",
	}
	PermissionAuthSso = &Permission{
		ID:          "auth_sso",
		Name:        "authentication.permissions.auth_sso.name",
		Description: "authentication.permissions.auth_sso.description",
		Category:    "org_tools",
	}
	PermissionAccessSecureAPI = &Permission{
		ID:          "access_secure_api",
		Name:        "authentication.permissions.access_secure_api.name",
		Description: "authentication.permissions.access_secure_api.description",
		Category:    "org_tools",
	}

	PermissionAccountManagerDedicated = &Permission{
		ID:          "account_manager_dedicated",
		Name:        "authentication.permissions.account_manager_dedicated.name",
		Description: "authentication.permissions.account_manager_dedicated.description",
		Category:    "support",
	}
)
