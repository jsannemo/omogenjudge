from typing import Optional, Union

from django.contrib.auth.models import AnonymousUser
from django.db.models import QuerySet

from omogenjudge.storage.models import Contest, Team, Account


def contest_teams(contest: Contest) -> QuerySet[Team]:
    return contest.team_set.prefetch_related('teammember_set').all()


def contest_team_for_user(contest: Contest, user: Union[AnonymousUser, Account]) -> Optional[Team]:
    if not user.is_authenticated:
        return None
    try:
        return contest.team_set.filter(teammember__account_id=user.account_id).prefetch_related('teammember_set').get()
    except Team.DoesNotExist:
        return None
