#include <iostream>

using namespace std;

int add(int a, int b) {
  if (a == 0) return b;
  int what[1000];
  what[0] = add(a - 1, b + 1);
  what[1] = add(0, -10);
  what[2] = add(0, 10);
  return what[0] + what[1] + what[2];
}

int main() {
  int a, b;
  cin >> a >> b;
  if (a + b <= 2000) {
    cout << add(a, b) << endl;
  } else {
    cout << a + b << endl;
  }
}
