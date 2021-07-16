import dataclasses
import typing

from django.http import HttpRequest, HttpResponse
from django.utils.functional import SimpleLazyObject

from omogenjudge.contests.lookup import contest_for_request
from omogenjudge.problems.lookup import contest_problems
from omogenjudge.storage.models import Contest


class ActiveContestMiddleware:
    def __init__(self, get_response: typing.Callable[[HttpRequest], HttpResponse]):
        self.get_response = get_response

    def __call__(self, request: HttpRequest) -> HttpResponse:
        request.contest = SimpleLazyObject(lambda: contest_for_request(request))
        response = self.get_response(request)
        return response


@dataclasses.dataclass
class ContestContext:
    title: str
    has_started: bool
    has_ended: bool


def _to_context(contest: Contest) -> ContestContext:
    return ContestContext(
        title=contest.title,
        has_started=contest.has_started,
        has_ended=contest.has_ended,
    )


def contest_context(request: HttpRequest):
    contest: Contest = request.contest
    return {
        'contest': _to_context(contest),
        'all_contest_problems': SimpleLazyObject(lambda: contest_problems(contest)),
        'contest_problems': SimpleLazyObject(lambda: contest_problems(contest) if contest.has_started else []),
    }
