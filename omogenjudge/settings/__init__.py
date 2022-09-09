import os
import sys

from omogenjudge.settings.base import *

if 'test' in sys.argv:
    from .tests import *
elif os.environ.get("PRODUCTION") == '1':
    from .production import *
else:
    from .local_development import *
