// Package requestactor contains the temporary transport context key used by
// the legacy Wallet surface. New business usecases should receive explicit
// actor/service capabilities rather than decode transport context themselves.
package requestactor

const UserContextKey = "user-context"
