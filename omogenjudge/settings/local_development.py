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

STATICFILES_DIRS = [
    BASE_DIR.parent / "output" / "frontend_assets",
]
