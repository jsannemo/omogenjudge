import dataclasses
from typing import Optional

from django.http import Http404, HttpResponse

from omogenjudge.contests.contest_groups import groups_by_shortnames, root_contest_groups
from omogenjudge.storage.models import Contest, ContestGroup, ContestGroupContest, Team
from omogenjudge.teams.lookup import contest_team_for_user
from omogenjudge.util.django_types import OmogenRequest
from omogenjudge.util.templates import render_template


@dataclasses.dataclass
class ArchiveContest:
    contest: ContestGroupContest
    my_team: Optional[Team] = None


@dataclasses.dataclass
class ArchiveArgs:
    current_groups: list[ContestGroup]
    groups: list[ContestGroup]
    contests: list[ArchiveContest]


def view_archive(request: OmogenRequest, group_path: Optional[str] = None) -> HttpResponse:
    if group_path:
        try:
            current_groups = groups_by_shortnames(group_path.split("/"))
        except ContestGroup.DoesNotExist:
            raise Http404
        groups = current_groups[-1].subgroups
        contests = [ArchiveContest(
            contest=cgc,
            my_team=contest_team_for_user(cgc.contest, request.user)
        )
            for cgc in current_groups[-1].subcontests
        ]
    else:
        current_groups = []
        groups = root_contest_groups()
        contests = []

    return render_template(request, 'archive/view_archive.html', ArchiveArgs(current_groups, groups, contests))
