from omogenjudge.settings.base import *  # noqa

SECRET_KEY = 'not a very secret key'

DATABASES = {
    'default': {
        'ENGINE': 'django.db.backends.sqlite3',
    }
}
