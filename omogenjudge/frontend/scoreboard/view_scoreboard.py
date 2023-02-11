import dataclasses

from django.http import HttpResponse
from django.utils import timezone

from omogenjudge.contests.scoreboard import ScoreboardMaker, load_scoreboard
from omogenjudge.frontend.decorators import requires_contest
from omogenjudge.storage.models import Contest
from omogenjudge.teams.lookup import contest_team_for_user
from omogenjudge.util.django_types import OmogenRequest
from omogenjudge.util.templates import render_template


@dataclasses.dataclass
class ScoreboardArgs:
    scoreboard: ScoreboardMaker


@requires_contest
def view_scoreboard(request: OmogenRequest, contest: Contest) -> HttpResponse:
    my_team = contest_team_for_user(contest, request.user)
    scoreboard = None
    if my_team and my_team.contest_start_time:
        now = timezone.now()
        elapsed = now - my_team.contest_start_time
        if elapsed <= contest.duration:
            scoreboard = load_scoreboard(contest, now=now, at_time=elapsed)
    if not scoreboard:
        scoreboard = load_scoreboard(contest)
    return render_template(request, 'scoreboard/view_scoreboard.html', ScoreboardArgs(scoreboard))
