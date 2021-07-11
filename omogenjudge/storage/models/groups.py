from django.db import models

from .accounts import Account
from ...util import django_fields


class GroupContest(models.Model):
    group = models.ForeignKey('Group', models.CASCADE)
    contest = models.ForeignKey('Contest', models.CASCADE)

    class Meta:
        db_table = 'group_contest'
        unique_together = (('group', 'contest'),)


class Group(models.Model):
    group_id = models.AutoField(primary_key=True)
    group_name = django_fields.TextField()

    class Meta:
        db_table = 'group'


class GroupMember(models.Model):
    group = models.ForeignKey(Group, models.CASCADE)
    account = models.ForeignKey(Account, models.CASCADE)
    admin = models.BooleanField()

    class Meta:
        db_table = 'account_group_member'
        unique_together = (('group', 'account'),)
