from django.contrib import admin
from django.db import models

from .accounts import Account
from .contests import Contest


class Team(models.Model):
    team_id = models.AutoField(primary_key=True)
    contest = models.ForeignKey(Contest, models.CASCADE)
    team_name = models.TextField(null=True)
    team_data = models.JSONField(blank=True, default=dict)

    def __str__(self):
        return self.team_name + " (in " + str(self.contest) + ")"

    class Meta:
        db_table = 'team'


class TeamMember(models.Model):
    team = models.ForeignKey(Team, models.CASCADE)
    account = models.ForeignKey(Account, models.CASCADE)

    class Meta:
        db_table = 'team_member'
        unique_together = (('team', 'account'),)
