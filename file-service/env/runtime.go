package env

// ENV_RUNTIME is assigned by config.LoadConfig at process composition time.
// Package import must not read process environment implicitly.
var ENV_RUNTIME string
