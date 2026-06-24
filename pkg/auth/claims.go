package auth

import "time"

// System roles, mirroring the auth-api `sys_role` claim.
const (
	RoleUser  = "USER"
	RoleAdmin = "ADMIN"
)

// Identity is the authenticated subject handed to consumers.
//
// Email/Username/ImgURL are populated only by remote /auth/validate; offline
// validation sees only what the access token carries (UserUUID, SysRole).
type Identity struct {
	UserUUID  string
	SysRole   string
	Email     string    // remote validation only
	Username  string    // remote validation only
	ImgURL    *string   // remote validation only
	ExpiresAt time.Time // from the token's exp (offline) — zero when remote
}

// IsAdmin reports whether the subject holds the ADMIN system role.
func (i *Identity) IsAdmin() bool { return i.SysRole == RoleAdmin }
