package adminpanel

func InitButtonsAdmin() {
	AdminPanelInline.Inline(
		AdminPanelInline.Row(BtnAllowSubscription),
		AdminPanelInline.Row(BtnDenySubscription),
		AdminPanelInline.Row(BtnListOfTransactions),
		AdminPanelInline.Row(BtnGivePermissions),
		AdminPanelInline.Row(BtnBackPermissions),
		AdminPanelInline.Row(BtnRefundAdmin),
		AdminPanelInline.Row(BtnBroadcastMessage),
		AdminPanelInline.Row(BtnSetComission),
	)
}

func InitPermissionsButtons() {
	PermissionsInline.Inline(
		PermissionsInline.Row(BtnGrantAdminPanel),
		PermissionsInline.Row(BtnGrantRefundAccess),
		PermissionsInline.Row(BtnGrantTechSupport),
	)
}

func InitBackPermissionsButtons() {
	BackPermissionsInline.Inline(
		BackPermissionsInline.Row(BtnBackAdminPanel),
		BackPermissionsInline.Row(BtnBackRefundAccess),
		BackPermissionsInline.Row(BtnBackTechSupport),
	)
}