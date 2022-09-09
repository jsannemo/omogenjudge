import dataclasses
import mimetypes
from typing import Optional, List

from django.http import Http404, HttpRequest, HttpResponse
from django.shortcuts import redirect

from omogenjudge.frontend.decorators import only_started_contests
from omogenjudge.frontend.problems.submit import SOURCE_CODE_LIMIT, SubmitForm
from omogenjudge.frontend.submissions.view_submission import SubmissionWithSubtasks
from omogenjudge.problems.lookup import NoSuchLanguage, find_statement_file, get_problem_for_view
from omogenjudge.problems.permissions import can_view_problem
from omogenjudge.problems.testgroups import get_subtask_scores, get_submission_subtask_scores
from omogenjudge.storage.models import Problem, ProblemStatementFile
from omogenjudge.submissions.lookup import list_account_problem_submissions
from omogenjudge.util.django_types import OmogenRequest
from omogenjudge.util.templates import render_template


@dataclasses.dataclass
class ViewArgs:
    statement_title: str
    statement_html: str
    statement_license: str
    statement_authors: str
    timelim_seconds: str
    timelim_ms: int
    memlim_mb: str
    is_scoring: bool
    subtask_scores: List[float]
    submit_form: SubmitForm
    source_code_limit: int
    submissions: list[SubmissionWithSubtasks]


@only_started_contests
def view_problem(request: OmogenRequest, short_name: str, language: Optional[str] = None) -> HttpResponse:
    try:
        problem, statement = get_problem_for_view(short_name, language=language)
    except Problem.DoesNotExist:
        raise Http404
    except NoSuchLanguage:
        return redirect('problem', short_name=short_name)

    if not can_view_problem(problem):
        raise Http404

    subtasks = get_subtask_scores(problem.current_version)
    submissions = list_account_problem_submissions(account=request.user, problem=problem,
                                                   limit=20) if request.user.is_authenticated else []
    submissions_with_subtasks = [
        SubmissionWithSubtasks(submission, get_submission_subtask_scores(list(submission.current_run.group_runs.all()),
                                                                         subtasks=len(subtasks))) for
        submission in submissions]

    args = ViewArgs(
        statement_title=statement.title,
        statement_html=statement.html,
        statement_license=problem.license,
        statement_authors=', '.join(problem.author),
        timelim_seconds=str(round(problem.current_version.time_limit_ms / 1000, ndigits=1)),
        timelim_ms=problem.current_version.time_limit_ms,
        memlim_mb='{:.0f}'.format(round(problem.current_version.memory_limit_kb / 1000)),
        is_scoring=problem.current_version.scoring,
        subtask_scores=subtasks,
        submit_form=SubmitForm(problem.short_name),
        source_code_limit=SOURCE_CODE_LIMIT,
        submissions=submissions_with_subtasks,
    )
    return render_template(request, 'problems/view_problem.html', args)


@only_started_contests
def problem_statement_file(request: OmogenRequest,
                           short_name: str,
                           file_path: str) -> HttpResponse:
    try:
        file = find_statement_file(short_name, file_path)
    except ProblemStatementFile.DoesNotExist:
        raise Http404
    mime_type, encoding = mimetypes.guess_type(request.path)
    if encoding:
        mime_type = f'{mime_type}; charset={encoding}'
    return HttpResponse(content=file.statement_file.file_contents,
                        content_type=mime_type)
