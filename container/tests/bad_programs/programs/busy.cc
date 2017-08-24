#include <unistd.h>

int busy() {
    int ret = 0;
    for (unsigned int i = 0; i < 1000000000; ++i) {
        ret ^= i * i;
    }
    return ret;
}

int main() {
    int res = 0;
    for (int i = 0; i < 8; ++i) {
        if (fork() == 0) {
            res += busy();
            break;
        }
    }
    busy();
}
