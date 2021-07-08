from omogenjudge.settings.base import LOGGING as DEFAULT_LOGGING

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

LOGGING = DEFAULT_LOGGING
LOGGING["loggers"].update({
    'django.db.backends': {
        'level': 'DEBUG',
        'handlers': ['console'],
    }
})
