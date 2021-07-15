import dataclasses

from django.http import HttpRequest, HttpResponse

from omogenjudge.util.templates import render_template


@dataclasses.dataclass
class HomeArgs:
    pass


def home(request: HttpRequest) -> HttpResponse:
    return render_template(request, 'home/home.html', HomeArgs())
