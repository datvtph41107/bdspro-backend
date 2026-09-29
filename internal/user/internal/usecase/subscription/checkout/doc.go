// Package checkout owns the User-side subscription checkout decision.
// It resolves authoritative Plan terms, classifies current Subscription
// state, durably freezes the resolved command, then issues one semantic Payment
// CreateOrder command. Durable replay uses the frozen snapshot before reading
// mutable Plan/Subscription state again.
package checkout
