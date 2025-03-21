//go:build tools
// +build tools

package tools

import (
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"
)

/**
 * This file is here to manage required dev dependancies.
 * The sqlc-dev package is never used in the code base to is removed from the go.mod file.
 * However it is a requried dependancy for the generate-sql make command used to generate
*  the sql query code.
*/
