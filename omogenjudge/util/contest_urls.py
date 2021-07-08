from django.http import HttpResponse
from django.shortcuts import redirect
from django.urls import reverse

from omogenjudge.util.request_global import current_request


def redirect_contest(view_name: str, *args, **kwargs) -> HttpResponse:
    request = current_request()
    assert request.contest
    if not request.contest_site:
        kwargs["contest_short_name"] = request.contest.short_name
        view_name = "contest-" + view_name
    return redirect(view_name, *args, **kwargs)


def reverse_contest(view_name: str, *args, **kwargs) -> str:
    request = current_request()
    assert request.contest
    if args:
        if not request.contest_site:
            reverse_args = [request.contest.short_name] + list(args)
            view_name = "contest-" + view_name
        else:
            reverse_args = list(args)
        return reverse(view_name, args=reverse_args)
    if not request.contest_site:
        kwargs["contest_short_name"] = request.contest.short_name
        view_name = "contest-" + view_name
    return reverse(view_name, kwargs=kwargs)
