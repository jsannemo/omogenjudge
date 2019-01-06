#pragma once

#include "glog/logging.h"
#include "grpc++/grpc++.h"

using grpc::Status;

// A result of a request that can fail. An actual result value is only given
// when the request succeeded, and the status will then be set to OK. When
// the request failed, a non-OK status will be set and no value can be
// retrieved.
template <typename T>
class StatusOr {
  Status _status;
  T _result;

 public:
  // Construct a status with a non-OK code.
  StatusOr(const Status& status) : _status(status) {
    CHECK(!status.ok()) << "Constructing OK status without result";
  }

  // Construct a status with a given result and an OK status.
  StatusOr(const T& result) : _status(Status()), _result(result) {}

  T value() {
    CHECK(_status.ok()) << "Attempting to get a value of a failed result";
    return _result;
  }

  Status status() { return _status; }

  bool ok() { return _status.ok(); }
};
