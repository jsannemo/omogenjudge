#include <unistd.h>

int main() {
  int ret = fork();
  if (ret == 0) {
    main();
  } else {
    while (true) sleep(1);
  }
}
