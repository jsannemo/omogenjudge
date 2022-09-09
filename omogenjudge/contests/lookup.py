from django.db.models import QuerySet

from omogenjudge.storage.models import Contest, ContestGroup
from omogenjudge.util.django_types import OmogenRequest


def _only_published_queryset() -> QuerySet[Contest]:
    return Contest.objects.filter(published=True)


def active_contest_queryset() -> QuerySet[Contest]:
    return _only_published_queryset().prefetch_related("problems")


def contest_for_request(request: OmogenRequest) -> Contest:
    return active_contest_queryset().get(
        host_name=request.get_host()
    )


def contest_from_shortname(short_name: str) -> Contest:
    return active_contest_queryset().get(
        short_name=short_name
    )


def contests_in_group(group: ContestGroup) -> list[Contest]:
    return list(group.contests.filter(published=True))
