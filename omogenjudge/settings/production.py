import toml

with open("/etc/omogen/web.toml", "r") as f:
    config = toml.load(f)

SECRET_KEY = config['web']['secret_key']

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
