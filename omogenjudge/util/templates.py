from typing import Optional

from django.http import HttpRequest, HttpResponse
from django.shortcuts import render

from omogenjudge.util.serialization import IsDictable


def render_template(request: HttpRequest, template: str, args: Optional[IsDictable]) -> HttpResponse:
    return render(request, template, context=args.__dict__ if args else {})
