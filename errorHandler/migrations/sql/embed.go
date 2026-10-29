// Package migrationsql embeds the SQL migrations so they ship with the module.
package migrationsql

import "embed"

//go:embed *.sql
var FS embed.FS
