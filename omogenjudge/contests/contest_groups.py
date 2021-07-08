from typing import Sequence

from omogenjudge.storage.models import ContestGroup


def groups_by_shortnames(short_names: Sequence[str]) -> list[ContestGroup]:
    groups: list[ContestGroup] = []
    for short_name in short_names:
        if not short_name:
            continue
        groups.append(ContestGroup.objects.get(parent=groups[-1] if groups else None, short_name=short_name))
    return groups


def groups_by_parent(parent: ContestGroup) -> list[ContestGroup]:
    return list(parent.groups.all())


def root_contest_groups() -> list[ContestGroup]:
    return list(ContestGroup.objects.filter(parent=None))
