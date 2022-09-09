import dataclasses
from typing import Optional

from django.http import HttpResponse, Http404

from omogenjudge.contests.contest_groups import groups_by_shortnames, root_contest_groups, groups_by_parent
from omogenjudge.contests.lookup import contests_in_group
from omogenjudge.storage.models import ContestGroup, Contest
from omogenjudge.util.django_types import OmogenRequest
from omogenjudge.util.templates import render_template


@dataclasses.dataclass
class ArchiveArgs:
    current_groups: list[ContestGroup]
    groups: list[ContestGroup]
    contests: list[Contest]


def view_archive(request: OmogenRequest, group_path: Optional[str] = None) -> HttpResponse:
    if group_path:
        try:
            current_groups = groups_by_shortnames(group_path.split("/"))
        except ContestGroup.DoesNotExist:
            raise Http404
        groups = groups_by_parent(current_groups[-1])
        contests = contests_in_group(current_groups[-1])
    else:
        current_groups = []
        groups = root_contest_groups()
        contests = []

    return render_template(request, 'archive/view_archive.html', ArchiveArgs(current_groups, groups, contests))
