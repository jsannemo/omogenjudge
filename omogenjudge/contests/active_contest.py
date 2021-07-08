import typing

from django.http import HttpRequest, HttpResponse, Http404
from django.utils.functional import SimpleLazyObject

from omogenjudge.contests.lookup import contest_for_request, contest_from_shortname
from omogenjudge.problems.lookup import contest_problems
from omogenjudge.storage.models import Contest
from omogenjudge.teams.lookup import contest_team_for_user
from omogenjudge.util.django_types import OmogenRequest


class ActiveContestMiddleware:
    def __init__(self, get_response: typing.Callable[[HttpRequest], HttpResponse]):
        self.get_response = get_response

    def __call__(self, request: OmogenRequest) -> HttpResponse:
        return self.get_response(request)

    def process_view(self, request: OmogenRequest, view_func, view_args, view_kwargs):
        request.contest_site = False
        if "contest_short_name" in view_kwargs:
            short_name = view_kwargs.pop("contest_short_name")
            try:
                request.contest = contest_from_shortname(short_name)
            except Contest.DoesNotExist:
                raise Http404
        else:
            try:
                request.contest = contest_for_request(request)
                if request.contest:
                    request.contest_site = True
            except Contest.DoesNotExist:
                request.contest = None


def contest_context(request: OmogenRequest) -> typing.Dict[str, typing.Any]:
    if request.contest:
        contest = request.contest
        return {
            'contest': contest,
            'all_contest_problems': SimpleLazyObject(lambda: contest_problems(contest)),
            'contest_team': SimpleLazyObject(lambda: contest_team_for_user(contest, request.user)),
            'contest_problems': SimpleLazyObject(lambda: contest_problems(contest) if contest.has_started else []),
        }
    return {
        'contest': None
    }
