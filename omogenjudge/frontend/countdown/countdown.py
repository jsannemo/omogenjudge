from django.http import HttpResponse

from omogenjudge.frontend.decorators import requires_contest
from omogenjudge.util.contest_urls import redirect_contest
from omogenjudge.util.django_types import OmogenRequest
from omogenjudge.util.templates import render_template


@requires_contest
def countdown(request: OmogenRequest, contest) -> HttpResponse:
    if contest.has_started:
        return redirect_contest('problems')
    return render_template(request, 'countdown/countdown.html', None)
