-- Intentionally data-preserving. The canonical tables may have received live
-- references after this migration; deleting rows during rollback would be
-- less correct than retaining the copied public catalog. Reapplying the up
-- migration is idempotent.
SELECT 1;
