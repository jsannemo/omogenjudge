from omogenjudge.settings.base import *  # noqa

from .base import BASE_DIR

SECRET_KEY = 'not a very secret key'

DEBUG = True
ALLOWED_HOSTS = ['*']

DATABASES = {
    'default': {
        'ENGINE': 'django.db.backends.postgresql_psycopg2',
        'NAME': 'omogenjudge',
        'USER': 'omogenjudge',
        'PASSWORD': 'omogenjudge',
        'HOST': 'localhost',
        'PORT': '5432',
    }
}

INTERNAL_IPS = [
    '127.0.0.1',
]

LOGGING["loggers"].update({
    'django.db.backends': {
        'level': 'DEBUG',
        'handlers': ['console'],
    }
})
STATICFILES_DIRS = [
    BASE_DIR.parent / "output" / "frontend_assets",
]

REQUIRE_EMAIL_AUTH = False
