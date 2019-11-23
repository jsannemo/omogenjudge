package models

import "fmt"

// A License is the license under which a problem is published.
type License string

const (
	// The Creative Commons by-sa 3 license.
	LicenseCcBySa3 License = "CC_BY_SA_3"
	// A license to reproduce the problem with specific permission.
	LicensePermission License = "BY_PERMISSION"
	// A license for problems that may freely be reproduced for any purpose.
	LicensePublicDomain License = "PUBLIC_DOMAIN"
	// A license for problems that you lack the rights to reproduce.
	LicensePrivate License = "PRIVATE"
)

// String returns a user-friendly name for a license.
func (l License) String() string {
	switch l {
	case LicenseCcBySa3:
		return "CC BY-SA 3.0"
	case LicensePublicDomain:
		return "Fri användning"
	case LicensePermission:
		return "Används med tillåtelse"
	case LicensePrivate:
		return "Endast för privat användning"
	}
	// cast back to string to avoid recursion
	panic(fmt.Errorf("unknown license: %s", string(l)))
}
