SECRET_KEY = 'not a very secret key'

ALLOWED_HOSTS = []

DEBUG = True

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