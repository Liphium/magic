package constants

// All constants used throughout Magic
const (
	// Everything written into locals for fiber
	LocalsAccountID       = "ses:acc"
	LocalsAccountName     = "ses:name"
	LocalsPermissionLevel = "ses:plvl"
	LocalsWizard          = "wiz"
	LocalsForgeBuild      = "fg:build"

	// Name for cookies used by Magic
	CookieMagicSession  = "mgc:session"
	CookieGitHubSession = "mgc:gh:session"

	// All permissions and their levels
	PermissionAdmin uint = 100
)
