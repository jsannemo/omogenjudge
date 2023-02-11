import typing

from django.utils import timezone

from omogenjudge.storage.models import Contest, Team


def contest_start_for_team(contest: Contest, team: typing.Optional[Team]) -> typing.Optional[int]:
    if team and team.contest_start_time:
        return int(team.contest_start_time.timestamp())
    return int(contest.start_time.timestamp()) if contest.start_time else None


def contest_has_started_for_team(contest: Contest, team: typing.Optional[Team]) -> bool:
    if team and team.contest_start_time:
        return timezone.now() >= team.contest_start_time
    if contest.flexible_start_window_end_time:
        return contest.flexible_start_window_end_time >= timezone.now()
    if contest.only_virtual_contest:
        return True
    return contest.has_started


def contest_has_ended_for_team(contest: Contest, team: typing.Optional[Team]) -> bool:
    if team and team.contest_start_time:
        return timezone.now() > team.contest_start_time + contest.duration
    if contest.flexible_start_window_end_time:
        return contest.flexible_start_window_end_time >= timezone.now()
    if contest.only_virtual_contest:
        return True
    return contest.has_ended
