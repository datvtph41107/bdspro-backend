package paymentcompleted

import (
	"testing"
	"time"
)

func validEvent() V1 {
	now := time.Date(2026, 8, 28, 0, 0, 0, 0, time.UTC)
	return V1{
		EventID: "payment.completed.order.1", OrderID: 1,
		SubjectKind: "profile", SubjectID: "42", ProductCode: "qhpro", PlanCode: "qhpro.pro",
		PlanVersionID: 3, PlanVersion: "2.0.0", Currency: "VND", AmountMinor: 19900000,
		FundsConfirmedAt: now, CompletedAt: now.Add(time.Minute),
	}
}

func TestV1RejectsIncompleteFact(t *testing.T) {
	if !validEvent().IsValid() {
		t.Fatal("valid payment event rejected")
	}
	cases := []func(*V1){
		func(e *V1) { e.EventID = "" },
		func(e *V1) { e.SubjectKind = "unknown" },
		func(e *V1) { e.AmountMinor = 0 },
		func(e *V1) { e.FundsConfirmedAt = time.Time{} },
	}
	for i, mutate := range cases {
		e := validEvent()
		mutate(&e)
		if e.IsValid() {
			t.Fatalf("case %d accepted incomplete event: %+v", i, e)
		}
	}
}
