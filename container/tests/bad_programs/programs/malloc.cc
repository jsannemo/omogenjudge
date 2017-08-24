#include <cstdlib>

int main(int argc, char** argv) {
    int *ptr = static_cast<int*>(malloc(1000000000));
    for (int i = 0; i < 1000000000; i += 1024) {
        ptr[i] = i;
    }
}
