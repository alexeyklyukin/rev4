#!/usr/bin/env bash
psql -f test/sql/test.sql >test/result/test.out
diff -u test/result/test.out test/expected/test.out >test/test.diff || echo "ERROR: results is different from expected, examine test/test.diff"
