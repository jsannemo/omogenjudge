#include "validator.h"

void run() {
	int maxp = Arg("maxp");
	int maxq = Arg("maxq");
  int P = Int(1, maxp);
  Space();
  int Q = Int(1, maxq);
  Endl();
  vector<int> cells = SpacedInts(Q, 1, P);
  Eof();

  assert(is_sorted(cells.begin(), cells.end()));
  assert(adjacent_find(cells.begin(), cells.end(), greater_equal<int>()) == cells.end());
}
