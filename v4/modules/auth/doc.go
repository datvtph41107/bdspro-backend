// Package auth owns the v4 modular authentication surface.
// Login and permission/me flows are entered via controller → services;
// repository/entity packages are prepared for moving persistence off
// legacy AuthUsecase / RoleUsecase.
package auth
