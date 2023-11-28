package api

import (
	"strings"
)

func (r *RegisterRequest) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	r.Email = strings.TrimSpace(r.Email)
	r.Password = strings.TrimSpace(r.Password)

	if r.Name == "" || r.Email == "" || r.Password == "" {
		return ErrMissingField
	}

	if len(r.Password) < 8 {
		return ErrWeakPassword
	}

	return nil
}
