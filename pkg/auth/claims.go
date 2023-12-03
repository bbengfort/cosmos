package auth

import (
	"context"
	"strconv"
	"time"

	"github.com/bbengfort/cosmos/pkg/db/models"
	jwt "github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	Name        string   `json:"name,omitempty"`
	Email       string   `json:"email,omitempty"`
	Role        string   `json:"role,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

func NewClaimsForUser(ctx context.Context, u *models.User) (claims *Claims, err error) {
	claims = &Claims{
		Name:  u.Name.String,
		Email: u.Email,
	}

	claims.SetSubjectID(u.ID)

	var role *models.Role
	if role, err = u.Role(ctx); err != nil {
		return nil, err
	}
	claims.Role = role.Title

	var perms []*models.Permission
	if perms, err = u.Permissions(ctx); err != nil {
		return nil, err
	}

	claims.Permissions = make([]string, 0, len(perms))
	for _, perm := range perms {
		claims.Permissions = append(claims.Permissions, perm.Title)
	}

	return claims, nil
}

func (c *Claims) SetSubjectID(uid int64) {
	c.Subject = strconv.FormatInt(uid, 36)
}

func (c Claims) SubjectID() (int64, error) {
	return strconv.ParseInt(c.Subject, 36, 64)
}

func (c Claims) HasPermission(required string) bool {
	for _, permisison := range c.Permissions {
		if permisison == required {
			return true
		}
	}
	return false
}

func (c Claims) HasAllPermissions(required ...string) bool {
	for _, perm := range required {
		if !c.HasPermission(perm) {
			return false
		}
	}
	return true
}

// Used to extract expiration and not before timestamps without having to use public keys
var tsparser = &jwt.Parser{SkipClaimsValidation: true}

func ParseUnverified(tks string) (claims *jwt.RegisteredClaims, err error) {
	claims = &jwt.RegisteredClaims{}
	if _, _, err = tsparser.ParseUnverified(tks, claims); err != nil {
		return nil, err
	}
	return claims, nil
}

func ExpiresAt(tks string) (_ time.Time, err error) {
	var claims *jwt.RegisteredClaims
	if claims, err = ParseUnverified(tks); err != nil {
		return time.Time{}, err
	}
	return claims.ExpiresAt.Time, nil
}

func NotBefore(tks string) (_ time.Time, err error) {
	var claims *jwt.RegisteredClaims
	if claims, err = ParseUnverified(tks); err != nil {
		return time.Time{}, err
	}
	return claims.NotBefore.Time, nil
}
