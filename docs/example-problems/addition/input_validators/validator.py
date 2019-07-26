import re
import sys

INT_RE = "(0|[1-9][0-9]*)"

line = sys.stdin.readline()
assert re.match("^" + INT_RE + " " + INT_RE + "$", line)
a, b = map(int, line.split())

line = sys.stdin.readline()
assert len(line) == 0

sys.exit(42)
