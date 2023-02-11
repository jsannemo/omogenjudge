import dataclasses
import json
from typing import Dict, Optional

from omogenjudge.contests.contest_times import contest_has_ended_for_team, contest_has_started_for_team, \
    contest_start_for_team
from omogenjudge.teams.lookup import contest_team_for_user
from omogenjudge.util.django_types import OmogenRequest


@dataclasses.dataclass
class JsContext:
    contest_start_timestamp: Optional[int]
    contest_duration: int
    contest_started: bool
    contest_ended: bool
    only_virtual: bool


def js_context(request: OmogenRequest) -> Dict[str, str]:
    contest = request.contest
    if contest:
        team = contest_team_for_user(contest, request.user)
        return {
            'js_context':
                json.dumps(dataclasses.asdict(JsContext(
                    contest_start_timestamp=contest_start_for_team(contest, team),
                    contest_duration=int(contest.duration.total_seconds()),
                    contest_started=contest_has_started_for_team(contest, team),
                    contest_ended=contest_has_ended_for_team(contest, team),
                    only_virtual=contest.only_virtual_contest
                ))),
        }
    return {}
