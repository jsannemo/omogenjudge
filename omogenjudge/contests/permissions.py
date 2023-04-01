from typing import Optional

from omogenjudge.contests.contest_times import contest_has_started_for_team
from omogenjudge.storage.models import Contest, Team


def team_can_view_problems(contest: Contest, team: Optional[Team]) -> bool:
    if not contest.published:
        return False
    return contest_has_started_for_team(contest, team)