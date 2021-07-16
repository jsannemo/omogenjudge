import dataclasses

from django.contrib.auth.decorators import login_required
from django.http import Http404, HttpRequest, HttpResponse

from omogenjudge.storage.models import Submission
from omogenjudge.submissions.lookup import all_submissions_for_queue
from omogenjudge.util.templates import render_template


@dataclasses.dataclass
class QueueArgs:
    submissions: list[Submission]


@login_required
def submission_queue(request: HttpRequest) -> HttpResponse:
    if not request.user.is_superuser:
        raise Http404
    submissions = all_submissions_for_queue()
    return render_template(request, 'submissions/queue.html', QueueArgs(submissions))
