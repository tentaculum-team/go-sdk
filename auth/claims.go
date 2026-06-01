package auth

import (
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// UserType mirrors the JWT/middleware `user_type` claim verbatim.
// Do NOT confuse with the register DTO `plan` field (user|enterprise).
type UserType string

const (
	UserTypePersonal   UserType = "user"
	UserTypeEnterprise UserType = "enterprise_user"
)

// Role mirrors the JWT `role` claim.
type Role string

const (
	RoleUser  Role = "USER"
	RoleAdmin Role = "ADMIN"
)

// Identity is the authenticated subject handed to consumers.
//
// Email/Username are only populated by remote /auth/validate; offline
// validation cannot see them (not present in the access token).
type Identity struct {
	UserID     uuid.UUID
	OrgID      uuid.UUID
	UserType   UserType
	IsOwner    bool
	Role       Role
	Email      string // remote validation only
	Username   string // remote validation only
	AvatarUUID *uuid.UUID
	ExpiresAt  time.Time // from JWT exp (offline) — zero when remote
}

// claims is the wire shape of access/refresh tokens. Matches
// auth-api pkg/jwt.Claims exactly.
type claims struct {
	UserID   uuid.UUID `json:"user_id"`
	OrgID    uuid.UUID `json:"org_id"`
	UserType string    `json:"user_type"`
	IsOwner  bool      `json:"is_owner"`
	Role     string    `json:"role"`
	gojwt.RegisteredClaims
}

func (cl *claims) toIdentity() *Identity {
	id := &Identity{
		UserID:   cl.UserID,
		OrgID:    cl.OrgID,
		UserType: UserType(cl.UserType),
		IsOwner:  cl.IsOwner,
		Role:     Role(cl.Role),
	}
	if cl.ExpiresAt != nil {
		id.ExpiresAt = cl.ExpiresAt.Time
	}
	return id
}
