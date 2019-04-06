#include "time.h"

using std::chrono::duration_cast;
using std::chrono::steady_clock;
using std::chrono::time_point;

namespace omogen {
namespace util {

Stopwatch::Stopwatch() : start(steady_clock::now()) {}

long long Stopwatch::millis() {
  return duration_cast<std::chrono::milliseconds>(steady_clock::now() - start)
      .count();
}

}  // namespace util
}  // namespace omogen
