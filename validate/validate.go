// Package validate is a thin compatibility shim that re-exports validators
// from pkg/validator/<sub> under their pre-refactor flat names. Kept so older
// consumers (e.g. the auth service) keep building without import path changes.
// New code should import pkg/validator/<sub> directly.
package validate

import (
	addr "github.com/tentaculum-team/go-sdk/pkg/validator/addresses"
	company "github.com/tentaculum-team/go-sdk/pkg/validator/company"
	user "github.com/tentaculum-team/go-sdk/pkg/validator/user"
)

// user validators
func Mail(s string) error      { return user.Mail(s) }
func Username(s string) error  { return user.Username(s) }
func Password(s string) error  { return user.Password(s) }
func FirstName(s string) error { return user.FirstName(s) }
func LastName(s string) error  { return user.LastName(s) }
func Phone(s string) error     { return user.Phone(s) }

// address config aliases (auth passes zero-value configs explicitly)
type LabelConfig = addr.LabelConfig
type StreetConfig = addr.StreetConfig
type HouseNumberConfig = addr.HouseNumberConfig
type ComplementConfig = addr.ComplementConfig
type DistrictConfig = addr.DistrictConfig
type CityConfig = addr.CityConfig
type StateRegionConfig = addr.StateRegionConfig

// address validators
func Country(s string) error                               { return addr.Country(s) }
func Postal(s string) error                                { return addr.Postal(s) }
func Label(s string, cfg ...LabelConfig) error             { return addr.Label(s, cfg...) }
func Street(s string, cfg ...StreetConfig) error           { return addr.Street(s, cfg...) }
func HouseNumber(s string, cfg ...HouseNumberConfig) error { return addr.AddressNumber(s, cfg...) }
func Complement(s string, cfg ...ComplementConfig) error   { return addr.Complement(s, cfg...) }
func District(s string, cfg ...DistrictConfig) error       { return addr.District(s, cfg...) }
func City(s string, cfg ...CityConfig) error               { return addr.City(s, cfg...) }
func StateRegion(s string, cfg ...StateRegionConfig) error { return addr.StateRegion(s, cfg...) }

// company validators
func CNPJ(s string) error { return company.CNPJ(s) }
