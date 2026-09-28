// Package user owns the v4 modular user-profile surface.
// Profile info/me and public profile v3 are entered via controller → services;
// legacy ProfileUsecase remains the lookup backend during migration (v4 database.dsn).
package user
