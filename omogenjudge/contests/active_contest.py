import typing

from django.http import HttpRequest, HttpResponse
from django.utils.functional import SimpleLazyObject

from omogenjudge.contests.lookup import contest_for_request
from omogenjudge.problems.lookup import contest_problems


class ActiveContestMiddleware:
    def __init__(self, get_response: typing.Callable[[HttpRequest], HttpResponse]):
        self.get_response = get_response

    def __call__(self, request: HttpRequest) -> HttpResponse:
        request.contest = SimpleLazyObject(lambda: contest_for_request(request))
        response = self.get_response(request)
        return response


def contest_context(request: HttpRequest):
    return {
        'contest': request.contest,
        'contest_problems': SimpleLazyObject(lambda: contest_problems(request.contest) if request.contest else []),
    }
