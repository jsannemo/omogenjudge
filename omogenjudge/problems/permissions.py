from django.http import HttpRequest

from omogenjudge.storage.models import Problem


def can_view_problem(request: HttpRequest, problem: Problem) -> bool:
    if request.user.is_superuser:
        return True
    if not request.contest:
        return True
    return request.contest.has_started
