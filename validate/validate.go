// Package validate is a thin compatibility shim that re-exports the user
// validators from pkg/validator/user under their pre-refactor names. Kept so
// older consumers (e.g. the auth service) keep building without import path
// changes. New code should import pkg/validator/<sub> directly.
package validate

import user "github.com/tentaculum-team/go-sdk/pkg/validator/user"

func Mail(s string) error      { return user.Mail(s) }
func Username(s string) error  { return user.Username(s) }
func Password(s string) error  { return user.Password(s) }
func FirstName(s string) error { return user.FirstName(s) }
func LastName(s string) error  { return user.LastName(s) }
func Phone(s string) error     { return user.Phone(s) }
