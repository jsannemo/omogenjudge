import base64
import dataclasses

from django.contrib.auth.decorators import login_required
from django.http import Http404, HttpRequest, HttpResponse

from omogenjudge.frontend.decorators import requires_started_contest
from omogenjudge.storage.models import Problem, Submission
from omogenjudge.submissions.lookup import get_submission_for_view
from omogenjudge.util.templates import render_template


@dataclasses.dataclass
class ViewArgs:
    submission: Submission
    files: dict[str, str]


@requires_started_contest
@login_required
def view_submission(request: HttpRequest, sub_id: int) -> HttpResponse:
    try:
        submission = get_submission_for_view(sub_id)
    except Problem.DoesNotExist:
        raise Http404
    if submission.account != request.user and not request.user.is_superuser:
        raise Http404
    args = ViewArgs(
        submission=submission,
        files={file: base64.b64decode(content).decode('UTF-8', errors='ignore') for file, content in submission.submission_files['files'].items()}
    )
    return render_template(request, 'submissions/view_submission.html', args)
