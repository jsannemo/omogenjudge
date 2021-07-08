import dataclasses
from typing import Optional

from django.http import HttpResponse

from omogenjudge.contests.scoreboard import load_scoreboard, ScoreboardTeam, ScoreboardMaker
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
    scoreboard = load_scoreboard(contest)
    args = ContestProblemArgs(scoreboard=scoreboard)
    if request.user.is_authenticated:
        user_id = request.user.account_id
        if user_id in scoreboard.user_to_rank:
            args.team_results = scoreboard.scoreboard_teams[scoreboard.user_to_rank[user_id]]
    return render_template(request, 'problems/contest_problems.html', args)
