package models

type RoleID string

const (
	RoleIDSystemAdmin           RoleID = "system_admin"
	RoleIDSystemUser            RoleID = "system_user"
	RoleIDSupplierAdmin         RoleID = "supplier_admin"
	RoleIDSupplierVendorManager RoleID = "supplier_vendor_manager"
	RoleIDSupplierModerator     RoleID = "supplier_moderator"
	RoleIDCustomer              RoleID = "customer"
)
