import dataclasses

from django.http import HttpResponse

from omogenjudge.contests.scoreboard import ScoreboardMaker, load_scoreboard
from omogenjudge.frontend.decorators import requires_contest
from omogenjudge.storage.models import Contest
from omogenjudge.util.django_types import OmogenRequest
from omogenjudge.util.templates import render_template


@dataclasses.dataclass
class ScoreboardArgs:
    scoreboard: ScoreboardMaker


@requires_contest
def view_scoreboard(request: OmogenRequest, contest: Contest) -> HttpResponse:
    scoreboard = load_scoreboard(contest)
    return render_template(request, 'scoreboard/view_scoreboard.html', ScoreboardArgs(scoreboard))
