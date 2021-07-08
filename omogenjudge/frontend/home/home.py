import dataclasses

from django.http import HttpRequest, HttpResponse
from django.shortcuts import redirect

from omogenjudge.util.contest_urls import redirect_contest
from omogenjudge.util.django_types import OmogenRequest


@dataclasses.dataclass
class HomeArgs:
    pass


def home(request: OmogenRequest) -> HttpResponse:
    if request.contest:
        return redirect_contest('problems')
    else:
        return redirect('archive')
