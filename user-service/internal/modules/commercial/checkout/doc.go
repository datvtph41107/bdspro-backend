// Package checkout owns the User-side commercial checkout decision.
// It resolves authoritative Catalog terms, classifies current Subscription
// state, durably freezes the resolved command, then issues one semantic Payment
// CreateOrder command. Durable replay uses the frozen snapshot before reading
// mutable Catalog/Subscription state again.
package checkout
