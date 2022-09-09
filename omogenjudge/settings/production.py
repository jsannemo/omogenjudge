from omogenjudge.settings.base import *  # noqa

import toml

with open("/etc/omogen/web.toml", "r") as f:
    config = toml.load(f)

SECRET_KEY = config['web']['secret_key']

MAILJET_API_KEY = config['email']['mailjet_api_key']
MAILJET_API_SECRET = config['email']['mailjet_api_secret']

SESSION_COOKIE_SECURE = True
CSRF_COOKIE_SECURE = True
SECURE_SSL_REDIRECT = True

ALLOWED_HOSTS = ['*']

DATABASES = {
    'default': {
        'ENGINE': 'django.db.backends.postgresql_psycopg2',
        'NAME': 'omogenjudge',
        'USER': 'omogenjudge',
        'PASSWORD': config['database']['password'],
        'HOST': 'localhost',
        'PORT': '5432',
    }
}

if "oauth" in config:
    OAUTH_DETAILS = config["oauth"]
