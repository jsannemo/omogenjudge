from typing import Optional

from django.http import HttpRequest, HttpResponse

from omogenjudge.util.templates import render_template


def react_app(request: HttpRequest, path: Optional[str] = None) -> HttpResponse:
    return render_template(request, "react.html", {})