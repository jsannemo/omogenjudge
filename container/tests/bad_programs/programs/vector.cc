#include <unistd.h>
#include <vector>

using namespace std;

int main(int argc, char** argv) {
    vector<int> val(1000000000);
    for (int i = 0; i < 1000000000; i += 1024) {
        val[i] = i;
    }
}
