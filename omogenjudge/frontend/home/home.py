import dataclasses

from django.http import HttpRequest, HttpResponse
from django.shortcuts import redirect

from omogenjudge.util.templates import render_template


@dataclasses.dataclass
class HomeArgs:
    pass


def home(request: HttpRequest) -> HttpResponse:
    return redirect('scoreboard')
    # return render_template(request, 'home/home.html', HomeArgs())
