from django.db import transaction
from django.utils import timezone

from omogenjudge.storage.models import Account, Contest, Team, TeamMember
from omogenjudge.teams.lookup import contest_team_for_user


class TeamExists(Exception): pass


def register_user_for_practice(contest: Contest, user: Account):
    with transaction.atomic():
        if contest_team_for_user(contest, user):
            raise TeamExists()
        team = Team(contest=contest, practice=True)
        team.save()
        TeamMember(team=team, account=user).save()


def register_user_for_virtual(contest: Contest, user: Account):
    with transaction.atomic():
        if contest_team_for_user(contest, user):
            raise TeamExists()
        team = Team(contest=contest, practice=True, contest_start_time=timezone.now())
        team.save()
        TeamMember(team=team, account=user).save()
