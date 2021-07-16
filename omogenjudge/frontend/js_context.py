import dataclasses
import json

from django.http import HttpRequest

from omogenjudge.storage.models import Contest


@dataclasses.dataclass
class JsContext:
    contest_start_timestamp: int
    contest_duration: int
    contest_started: bool
    contest_ended: bool


def js_context(request: HttpRequest):
    contest: Contest = request.contest
    return {
        'js_context':
            json.dumps(dataclasses.asdict(JsContext(
                contest_start_timestamp=int(contest.start_time.timestamp()),
                contest_duration=int(contest.duration.total_seconds()),
                contest_started=contest.has_started,
                contest_ended=contest.has_ended,
            ))),
    }
