from django.db import models

from omogenjudge.storage.models import Account, Contest
from omogenjudge.util import django_fields


class Team(models.Model):
    team_id = models.AutoField(primary_key=True)
    contest = models.ForeignKey(Contest, models.CASCADE)
    team_name = django_fields.TextField(null=True, blank=True)
    team_data = models.JSONField(blank=True, default=dict)
    # When this team starts the contest
    contest_start_time = models.DateTimeField(null=True, blank=True)
    practice = models.BooleanField(null=False, default=False)

    def display_name(self) -> str:
        if self.team_name:
            return self.team_name
        team_members = self.teammember_set.select_related("account").all()
        if team_members:
            return ', '.join(team_member.account.full_name for team_member in team_members)
        return 'Empty team'

    def __str__(self):
        return self.display_name() + " (in " + str(self.contest) + ")"

    class Meta:
        db_table = 'team'


class TeamMember(models.Model):
    team_member_id = models.AutoField(primary_key=True)
    team = models.ForeignKey(Team, models.CASCADE)
    account = models.ForeignKey(Account, models.CASCADE)

    class Meta:
        db_table = 'team_member'
        unique_together = (('team', 'account'),)
