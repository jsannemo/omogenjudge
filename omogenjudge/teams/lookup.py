from omogenjudge.storage.models import Contest, Team


def contest_teams(contest: Contest) -> list[Team]:
    return contest.team_set.prefetch_related('teammember_set').all()
