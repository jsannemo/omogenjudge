#include <cstdlib>
#include <sstream>

#include "format.h"
#include "util/error.h"

using std::string;
using std::stringstream;
using std::vector;

namespace omogenexec {

long long StringToLL(const std::string& str) {
    errno = 0;
    long long ret = strtoll(str.c_str(), nullptr, 10);
    if (errno == ERANGE) {
        OE_FATAL("strtoll");
    }
    return ret;
}

vector<string> Split(const string &s, char delim) {
    stringstream ss(s);
    string field;
    vector<string> res;
    while (std::getline(ss, field, delim)) {
        res.push_back(field);
    }
    return res;
}

}
