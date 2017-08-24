#include <fcntl.h>
#include <sstream>
#include <iostream>
#include <fstream>
#include <unistd.h>
#include <cstring>

using namespace std;

int main(int argc, char** argv) {
    string s(200000, 'a');
    string path = string(argv[argc-1]) + "/output";
    ofstream ofs(path);
    ofs << s << endl;
    ofs.close();
}
