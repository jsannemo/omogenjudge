import typing

from django.http import Http404, HttpRequest, HttpResponse
from django.utils.functional import SimpleLazyObject

from omogenjudge.contests.contest_times import contest_has_ended_for_team, contest_has_started_for_team
from omogenjudge.contests.lookup import contest_for_request, contest_from_shortname
from omogenjudge.contests.permissions import team_can_view_problems
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
        if request.contest:
            request.team = contest_team_for_user(request.contest, request.user)


def contest_context(request: OmogenRequest) -> dict[str, typing.Any]:
    if request.contest:
        contest = request.contest
        ctx: dict[str, typing.Any] = {
            'contest': contest,
            'all_contest_problems': SimpleLazyObject(lambda: contest_problems(contest)),
            'contest_team': request.team,
            'contest_started': contest_has_started_for_team(contest, request.team),
            'contest_ended': contest_has_ended_for_team(contest, request.team),
        }
        if not team_can_view_problems(contest, request.team):
            ctx['contest_problems'] = []
        else:
            ctx['contest_problems'] = SimpleLazyObject(lambda: contest_problems(contest))
        return ctx
    return {
        'contest': None
    }
