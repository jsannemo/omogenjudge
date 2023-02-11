from django.core.exceptions import BadRequest
from django.core.exceptions import BadRequest
from django.http import HttpResponse
from django.views.decorators.http import require_http_methods

from omogenjudge.frontend.decorators import requires_contest, requires_user
from omogenjudge.storage.models import Account, Contest
from omogenjudge.teams.register import register_user_for_practice, register_user_for_virtual
from omogenjudge.util.contest_urls import redirect_contest
from omogenjudge.util.django_types import OmogenRequest


@requires_user()
@requires_contest
@require_http_methods(["POST"])
def register(request: OmogenRequest, user: Account, contest: Contest) -> HttpResponse:
    type = request.POST["type"]
    if not contest.open_for_practice:
        raise BadRequest()
    if type == "practice":
        register_user_for_practice(contest, user)
    elif type == "virtual":
        register_user_for_virtual(contest, user)
    else:
        raise BadRequest()
    return redirect_contest('problems')
