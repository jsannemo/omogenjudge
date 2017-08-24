#include <unistd.h>

int main(int argc, char** argv) {
    sbrk(1000000000);
}
