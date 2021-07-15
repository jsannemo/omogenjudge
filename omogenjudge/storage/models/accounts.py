from django.contrib import admin
from django.contrib.auth.base_user import AbstractBaseUser
from django.contrib.auth.models import PermissionsMixin, UserManager
from django.db import models

from omogenjudge.util import django_fields


class Account(AbstractBaseUser, PermissionsMixin):
    USERNAME_FIELD = 'username'
    REQUIRED_FIELDS = ['full_name', 'email']

    objects = UserManager()

    account_id = models.AutoField(primary_key=True)
    username = django_fields.TextField(unique=True)
    full_name = django_fields.TextField()
    email = django_fields.TextField(unique=True)
    date_created = models.DateTimeField(auto_now_add=True)
    last_login = models.DateTimeField(null=True)
    is_staff = models.BooleanField(default=False)
    is_superuser = models.BooleanField(default=False)

    class Meta:
        db_table = 'account'
