// Error utilities for database connections
package db

import (
  "github.com/lib/pq"
)

type PgError int

const (
  NotPgError PgError = iota
  OtherError
  SerializationFailure
  DeadlockDetected
  UniquenessViolation
)

// PgErrCode converts an error returned from a transaction into a typed Postgres error
func PgErrCode(err error) PgError {
  if pgerr, ok := err.(*pq.Error); ok {
    code := pgerr.Code.Name()
    if code == "serialization_failure" {
      return SerializationFailure
    }
    if code == "deadlock_detected" {
      return DeadlockDetected
    }
    if code == "unique_violation" {
      return UniquenessViolation
    }
    return OtherError
  }
  return NotPgError
}

// RetriableError checks if the given Postgres error should be retried.
// Generally, this indicates the user attempted a transaction that conflicted with a concurrent transaction that should be retried.
func RetriableError(err error) bool {
  code := PgErrCode(err)
  return code == SerializationFailure || code == DeadlockDetected
}
