#include <glog/logging.h>

#include "absl/strings/str_cat.h"
#include "absl/strings/numbers.h"
#include "util/config.h"
#include "util/files.h"

namespace omogenexec {

int ReadConfigInt(const std::string& filePath, const std::string& key) {
  std::vector<std::string> configTokens = TokenizeFile(filePath);
  for (size_t idx = 0; idx < configTokens.size() - 1; idx++) {
    if (configTokens[idx] == absl::StrCat(key, ":")) {
      std::string value = configTokens[idx + 1];
      int ret;
      CHECK(absl::SimpleAtoi(configTokens[idx + 1], &ret)) << "Wrong type of value in config for " << key;
      return ret;
    }
  }
  LOG(FATAL) << "Missing key " << key;
}

}
