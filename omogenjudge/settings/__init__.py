from .base import *
import os

if os.environ.get("PRODUCTION") == '1':
    from .production import *
else:
    from .development import *
