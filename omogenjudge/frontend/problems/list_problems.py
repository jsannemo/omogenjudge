import dataclasses

from django.http import HttpRequest, HttpResponse
from django.shortcuts import redirect

from omogenjudge.frontend.decorators import requires_started_contest
from omogenjudge.problems.lookup import contest_problems
from omogenjudge.storage.models import Problem
from omogenjudge.util.templates import render_template


@dataclasses.dataclass
class ListArgs:
    problems: list[Problem]


@requires_started_contest
def list_problems(request: HttpRequest) -> HttpResponse:
    problems = contest_problems(request.contest)
    if problems:
        return redirect('problem', short_name=problems[0].problem.short_name)
    return render_template(request, 'problems/list_problems.html', ListArgs(problems))
