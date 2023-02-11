import dataclasses
from typing import Optional

from django.http import HttpResponse

from omogenjudge.contests.scoreboard import ScoreboardMaker, ScoreboardTeam, load_scoreboard
from omogenjudge.frontend.decorators import only_started_contests
from omogenjudge.storage.models import Contest
from omogenjudge.util.django_types import OmogenRequest
from omogenjudge.util.templates import render_template


@only_started_contests
def list_problems(request: OmogenRequest) -> HttpResponse:
    if request.contest:
        return list_contest_problems(request, request.contest)
    raise NotImplementedError


@dataclasses.dataclass
class ContestProblemArgs:
    scoreboard: ScoreboardMaker
    team_results: Optional[ScoreboardTeam] = None


@only_started_contests
def list_contest_problems(request: OmogenRequest, contest: Contest) -> HttpResponse:
    user = request.user
    scoreboard = load_scoreboard(contest)
    args = ContestProblemArgs(scoreboard=scoreboard)
    if user.is_authenticated:
        user_id = user.account_id
        if user_id in scoreboard.best_user_result:
            args.team_results = scoreboard.best_user_result[user_id]
    return render_template(request, 'problems/contest_problems.html', args)
