package db

import (
	"strings"
	"testing"
)

func TestSplitSQLStatements_DollarQuoteAndComment(t *testing.T) {
	script := `
-- comment; ignore
CREATE TABLE a (id int);
DO $$
BEGIN
  IF true THEN
    RAISE NOTICE 'x;y';
  END IF;
END $$;
SELECT 1;
`
	got := splitSQLStatements(script)
	if len(got) != 3 {
		t.Fatalf("want 3 stmts, got %d: %#v", len(got), got)
	}
	if !strings.Contains(got[1], "DO $$") || !strings.Contains(got[1], "END $$") {
		t.Fatalf("dollar-quote block split incorrectly: %q", got[1])
	}
}
