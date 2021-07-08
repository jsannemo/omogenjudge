from django.contrib.auth.base_user import AbstractBaseUser
from django.db import models


class Account(AbstractBaseUser):
    USERNAME_FIELD = 'username'
    REQUIRED_FIELDS = ['full_name', 'email']

    account_id = models.AutoField(primary_key=True)
    username = models.TextField(unique=True)
    password_hash = models.TextField()
    full_name = models.TextField()
    email = models.TextField(unique=True)
    date_created = models.DateTimeField(auto_now_add=True)
    last_login = models.DateTimeField()

    class Meta:
        db_table = 'account'
