package domain

import "testing"

func TestPlanningNewsVisibleOnValidation(t *testing.T) {
	t.Parallel()

	valid := []PlanningNewsVisibleOn{
		PlanningNewsVisibleOnTop,
		PlanningNewsVisibleOnNotable,
		PlanningNewsVisibleOnLatest,
		PlanningNewsVisibleOnMostRead,
		PlanningNewsVisibleOnQuickNews,
	}

	for _, value := range valid {
		value := value

		t.Run(value.String(), func(t *testing.T) {
			t.Parallel()

			if !value.IsValid() {
				t.Fatalf("expected %d to be valid", value)
			}

			if value.String() == "" {
				t.Fatalf("expected %d to have code", value)
			}

			if value.Label() == "" {
				t.Fatalf("expected %d to have label", value)
			}
		})
	}

	invalid := []PlanningNewsVisibleOn{
		0,
		1,
		9,
		11,
		29,
		31,
		49,
		51,
		99,
	}

	for _, value := range invalid {
		value := value

		t.Run("invalid", func(t *testing.T) {
			t.Parallel()

			if value.IsValid() {
				t.Fatalf("expected %d to be invalid", value)
			}
		})
	}
}

func TestParsePlanningNewsVisibleOn(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected PlanningNewsVisibleOn
	}{
		{
			input:    "top",
			expected: PlanningNewsVisibleOnTop,
		},
		{
			input:    " notable ",
			expected: PlanningNewsVisibleOnNotable,
		},
		{
			input:    "LATEST",
			expected: PlanningNewsVisibleOnLatest,
		},
		{
			input:    " MOST_READ ",
			expected: PlanningNewsVisibleOnMostRead,
		},
		{
			input:    "quick_news",
			expected: PlanningNewsVisibleOnQuickNews,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.input, func(t *testing.T) {
			t.Parallel()

			value, ok := ParsePlanningNewsVisibleOn(test.input)
			if !ok {
				t.Fatalf("expected %q to be accepted", test.input)
			}

			if value != test.expected {
				t.Fatalf(
					"unexpected value for %q: got=%d expected=%d",
					test.input,
					value,
					test.expected,
				)
			}
		})
	}

	if _, ok := ParsePlanningNewsVisibleOn("unknown"); ok {
		t.Fatal("unknown code must not be accepted")
	}
}

func TestResolvePlanningNewsFeed(t *testing.T) {
	t.Parallel()

	tests := []struct {
		feed     string
		expected PlanningNewsVisibleOn
	}{
		// visible_on = 10
		{"top", PlanningNewsVisibleOnTop},
		{"pinned", PlanningNewsVisibleOnTop},
		{"featured", PlanningNewsVisibleOnTop},

		// visible_on = 20
		{"notable", PlanningNewsVisibleOnNotable},
		{"attention", PlanningNewsVisibleOnNotable},

		// visible_on = 30
		{"latest", PlanningNewsVisibleOnLatest},

		// visible_on = 40
		{"most_read", PlanningNewsVisibleOnMostRead},
		{"trending", PlanningNewsVisibleOnMostRead},

		// visible_on = 50
		{"quick_news", PlanningNewsVisibleOnQuickNews},
		{"breaking", PlanningNewsVisibleOnQuickNews},
	}

	for _, test := range tests {
		test := test

		t.Run(test.feed, func(t *testing.T) {
			t.Parallel()

			value, ok := ResolvePlanningNewsFeed(test.feed)
			if !ok {
				t.Fatalf("expected feed %q to be supported", test.feed)
			}

			if value != test.expected {
				t.Fatalf(
					"unexpected visible_on for feed %q: got=%d expected=%d",
					test.feed,
					value,
					test.expected,
				)
			}
		})
	}
}

func TestResolvePlanningNewsFeedNormalizesInput(t *testing.T) {
	t.Parallel()

	value, ok := ResolvePlanningNewsFeed("  LATEST  ")
	if !ok {
		t.Fatal("expected normalized latest feed")
	}

	if value != PlanningNewsVisibleOnLatest {
		t.Fatalf(
			"unexpected latest value: got=%d expected=%d",
			value,
			PlanningNewsVisibleOnLatest,
		)
	}
}

func TestResolvePlanningNewsFeedRejectsUnknown(t *testing.T) {
	t.Parallel()

	invalid := []string{
		"",
		"unknown",
		"latest-news",
		"30",
	}

	for _, feed := range invalid {
		feed := feed

		t.Run(feed, func(t *testing.T) {
			t.Parallel()

			if _, ok := ResolvePlanningNewsFeed(feed); ok {
				t.Fatalf("feed %q must not be accepted", feed)
			}
		})
	}
}

func TestCanonicalPlanningNewsFeed(t *testing.T) {
	t.Parallel()

	tests := []struct {
		feed     string
		expected string
	}{
		{"featured", "top"},
		{"pinned", "top"},
		{"attention", "notable"},
		{"latest", "latest"},
		{"trending", "most_read"},
		{"breaking", "quick_news"},
	}

	for _, test := range tests {
		test := test

		t.Run(test.feed, func(t *testing.T) {
			t.Parallel()

			value, ok := CanonicalPlanningNewsFeed(test.feed)
			if !ok {
				t.Fatalf("expected %q to resolve", test.feed)
			}

			if value != test.expected {
				t.Fatalf(
					"unexpected canonical feed: got=%q expected=%q",
					value,
					test.expected,
				)
			}
		})
	}
}
