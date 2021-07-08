from django.http import HttpRequest

from omogenjudge.problems.lookup import contest_problems
from omogenjudge.storage.models import Problem
from omogenjudge.util.django_types import OmogenRequest
from omogenjudge.util.request_global import current_contest


def can_view_contest_problems() -> bool:
    contest = current_contest()
    if not contest.published:
        return False
    if contest.only_virtual_contest:
        return True
    else:
        return contest.has_started
