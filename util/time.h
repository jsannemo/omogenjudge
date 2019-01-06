#include <chrono>

namespace omogenexec {

// A stopwatch that starts measuring time when created, and can return the
// elapsed time since.
class Stopwatch {
  std::chrono::time_point<std::chrono::steady_clock> start;

 public:
  Stopwatch();
  long long millis();
};

}  // namespace omogenexec
