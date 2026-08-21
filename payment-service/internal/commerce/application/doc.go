// Package application owns Payment Commerce use-case orchestration.
//
// It deliberately contains no transport, environment, broker, or PostgreSQL
// mechanics. Order, Settlement, Fulfillment, and Recovery are separate
// responsibilities because they have different invariants and failure windows.
package application
