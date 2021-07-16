import dataclasses

from django.http import HttpRequest, HttpResponse

from omogenjudge.contests.scoreboard import Scoreboard, load_scoreboard
from omogenjudge.util.templates import render_template


@dataclasses.dataclass
class ScoreboardArgs:
    scoreboard: Scoreboard


def view_scoreboard(request: HttpRequest) -> HttpResponse:
    scoreboard = load_scoreboard(request.contest)
    return render_template(request, 'scoreboard/view_scoreboard.html', ScoreboardArgs(scoreboard))
