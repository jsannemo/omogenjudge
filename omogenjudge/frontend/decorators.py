from typing import Callable

from django.contrib.auth.decorators import login_required
from django.core.exceptions import PermissionDenied
from django.http import HttpResponse

from omogenjudge.util.contest_urls import redirect_contest


def only_started_contests(f, *, allow_practice=True):
    def wrapped(*args, **kwargs):
        request = args[0]
        if request.contest and not request.contest.has_started and (
                not allow_practice or not request.contest.only_virtual_contest) and \
                not request.user.is_superuser:
            return redirect_contest('countdown')
        return f(*args, **kwargs)

    return wrapped


def requires_contest(f):
    def wrapped(*args, **kwargs):
        request = args[0]
        if not request.contest:
            return HttpResponse('No active contest')
        return f(*args, contest=request.contest, **kwargs)

    return wrapped


def requires_user(only_superuser: bool = False) -> Callable[
    [Callable[..., HttpResponse]], Callable[..., HttpResponse]]:
    def requires_user_decorator(f: Callable[..., HttpResponse]) -> Callable[..., HttpResponse]:
        def wrapped(*args, **kwargs):
            request = args[0]
            user = request.user
            if only_superuser and not user.is_superuser:
                raise PermissionDenied("Forbidden")
            return f(*args, user=user, **kwargs)

        return login_required(wrapped)
    return requires_user_decorator
