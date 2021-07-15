import dataclasses

from django.http import HttpRequest, HttpResponse

from omogenjudge.problems.lookup import list_public_problems
from omogenjudge.storage.models import Problem
from omogenjudge.util.templates import render_template


@dataclasses.dataclass
class ListArgs:
    problems: list[Problem]


def list_problems(request: HttpRequest) -> HttpResponse:
    problems = list_public_problems()
    return render_template(request, 'problems/list_problems.html', ListArgs(problems))
