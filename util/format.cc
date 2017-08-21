#include <cstdlib>
#include <sstream>

#include "format.h"
#include "util/error.h"
#include "util/log.h"

using std::endl;
using std::string;
using std::stringstream;
using std::vector;

namespace omogenexec {

long long StringToLL(const std::string& str) {
    errno = 0;
    char *end  = const_cast<char*>(str.c_str() + str.size());
    long long ret = strtoll(str.c_str(), &end, 10);
    if (errno == ERANGE) {
        OE_FATAL("strtoll");
    }
    if (end == str.c_str()) {
        OE_LOG(FATAL) << "No numeric conversion could be performed on " << str << endl;
        OE_CRASH();
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
