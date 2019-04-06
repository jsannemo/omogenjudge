#include <fcntl.h>
#include <unistd.h>
#include <cstring>
#include <fstream>
#include <iostream>
#include <sstream>

using namespace std;

int main(int argc, char** argv) {
  string s(200000, 'a');
  string path = string(argv[argc - 1]) + "/output";
  cout << path << endl;
  ofstream ofs(path);
  ofs << s << endl;
  ofs.close();
}
