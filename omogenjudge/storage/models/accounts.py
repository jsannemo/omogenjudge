from django.contrib.auth.base_user import AbstractBaseUser
from django.contrib.auth.models import PermissionsMixin, UserManager
from django.db import models
from django.db.models.functions import Lower

from omogenjudge.util import django_fields


class Account(AbstractBaseUser, PermissionsMixin):
    USERNAME_FIELD = 'username'
    REQUIRED_FIELDS = ['full_name', 'email']

    objects = UserManager()

    account_id = models.AutoField(primary_key=True)
    username = django_fields.TextField(unique=True)
    full_name = django_fields.TextField()
    email = django_fields.TextField(unique=True)
    email_validated = models.BooleanField(default=False)
    date_created = models.DateTimeField(auto_now_add=True)
    last_login = models.DateTimeField(null=True)
    is_staff = models.BooleanField(default=False)
    is_superuser = models.BooleanField(default=False)

    class Meta:
        db_table = 'account'

        constraints = [
            models.UniqueConstraint(
                Lower('username'),  # type: ignore
                name='unique_case_insensitive_username'
            ),
            models.UniqueConstraint(
                Lower('email'),  # type: ignore
                name='unique_case_insensitive_email'
            ),
        ]
