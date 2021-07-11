from django.db import models

from .accounts import Account
from .contests import Contest


class Team(models.Model):
    team_id = models.AutoField(primary_key=True)
    contest = models.ForeignKey(Contest, models.CASCADE)
    team_name = models.TextField(null=True)
    virtual = models.BooleanField()
    unofficial = models.BooleanField()
    start_time = models.DateTimeField(null=True)
    team_data = models.JSONField()

    class Meta:
        db_table = 'team'


class TeamMember(models.Model):
    team = models.ForeignKey(Team, models.CASCADE)
    account = models.ForeignKey(Account, models.CASCADE)

    class Meta:
        db_table = 'team_member'
        unique_together = (('team', 'account'),)
