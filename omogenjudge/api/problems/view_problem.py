import dataclasses
from typing import Optional

from django.http import Http404, HttpRequest, HttpResponse

from omogenjudge.api.problems.model import ApiProblemLimits, ApiStatement
from omogenjudge.problems.lookup import NoSuchLanguage, get_problem_for_view
from omogenjudge.problems.permissions import can_view_problem
from omogenjudge.storage.models import Problem
from omogenjudge.util.templates import render_json


@dataclasses.dataclass
class ViewResponse:
    statement: ApiStatement
    limits: ApiProblemLimits


def view_problem(request: HttpRequest, short_name: str, language: Optional[str] = None) -> HttpResponse:
    try:
        problem, statement = get_problem_for_view(short_name, language=language)
    except (Problem.DoesNotExist, NoSuchLanguage):
        raise Http404
    if not can_view_problem(request, problem):
        raise Http404
    return render_json(ViewResponse(
        statement=ApiStatement.from_db_statement(statement),
        limits=ApiProblemLimits.from_db_version(problem.current_version),
    ))
