import dataclasses
import mimetypes
from typing import Optional

from django.http import Http404, HttpRequest, HttpResponse
from django.shortcuts import redirect
from django.urls import reverse

from omogenjudge.frontend.problems.submit import SOURCE_CODE_LIMIT, SubmitForm
from omogenjudge.problems.lookup import NoSuchLanguage, find_statement_file, lookup_for_viewing
from omogenjudge.storage.models import Problem, ProblemStatementFile, Submission
from omogenjudge.submissions.lookup import list_account_problem_submissions
from omogenjudge.util.templates import render_template


@dataclasses.dataclass
class ViewArgs:
    statement_title: str
    statement_html: str
    statement_license: str
    statement_authors: str
    timelim_seconds: str
    memlim_mb: str
    is_scoring: bool
    submit_form: SubmitForm
    source_code_limit: int
    submissions: list[Submission]


def view_problem(request: HttpRequest, short_name: str, language: Optional[str] = None) -> HttpResponse:
    try:
        problem, statement = lookup_for_viewing(short_name, language=language)
    except Problem.DoesNotExist:
        raise Http404
    except NoSuchLanguage:
        return redirect(reverse(view_problem, kwargs={'short_name': short_name}))
    args = ViewArgs(
        statement_title=statement.title,
        statement_html=statement.html,
        statement_license=problem.license,
        statement_authors=', '.join(problem.author),
        timelim_seconds='{:.2g}'.format(problem.current_version.time_limit_ms / 1000),
        memlim_mb='{:.0f}'.format(round(problem.current_version.memory_limit_kb / 1000)),
        is_scoring=problem.current_version.scoring,
        submit_form=SubmitForm(problem.short_name),
        source_code_limit=SOURCE_CODE_LIMIT,
        submissions=list_account_problem_submissions(account=request.user, problem=problem,
                                                     limit=20) if request.user.is_authenticated else [],
    )
    return render_template(request, 'problems/view_problem.html', args)


def problem_attachment(request: HttpRequest,
                       short_name: str,
                       file_path: str) -> HttpResponse:
    try:
        file = find_statement_file(short_name, file_path)
    except ProblemStatementFile.DoesNotExist:
        raise Http404
    mime_type = mimetypes.guess_type(request.path)
    if mime_type[1]:
        mime_type = f'{mime_type[0]}; charset={mime_type[1]}'
    else:
        mime_type = mime_type[0]
    return HttpResponse(content=file.statement_file.file_contents,
                        content_type=mime_type)
