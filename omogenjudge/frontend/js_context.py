import dataclasses
import json
from typing import Dict, Optional

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
        return {
            'js_context':
                json.dumps(dataclasses.asdict(JsContext(
                    contest_start_timestamp=int(contest.start_time.timestamp()) if contest.start_time else None,
                    contest_duration=int(contest.duration.total_seconds()),
                    contest_started=contest.has_started,
                    contest_ended=contest.has_ended,
                    only_virtual=contest.only_virtual_contest
                ))),
        }
    return {}
