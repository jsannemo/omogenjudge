import base64
import dataclasses
from typing import List, Optional

from django.http import Http404, HttpResponse

from omogenjudge.frontend.decorators import requires_user
from omogenjudge.problems.lookup import contest_problems
from omogenjudge.problems.testgroups import get_subtask_scores, get_submission_subtask_scores, \
    get_submission_subtask_groups
from omogenjudge.storage.models import Language, Problem, Submission, SubmissionGroupRun, Account, Contest
from omogenjudge.submissions.lookup import get_submission_for_view
from omogenjudge.util.django_types import OmogenRequest
from omogenjudge.util.templates import render_template


@dataclasses.dataclass
class SubmissionWithSubtasks:
    submission: Submission
    subtask_scores: List[float]
    subtask_groups: Optional[List[SubmissionGroupRun]] = None


@dataclasses.dataclass
class ProblemWithScores:
    problem: Problem
    subtask_scores: List[float]


@dataclasses.dataclass
class ViewArgs:
    author: Account
    problem_with_scores: ProblemWithScores
    submission_with_subtasks: SubmissionWithSubtasks
    language: str
    files: dict[str, str]


@requires_user()
def view_submission(request: OmogenRequest, sub_id: int, user: Account) -> HttpResponse:
    if request.contest:
        return view_contest_submission(request, request.contest, user, sub_id)
    raise NotImplementedError


def view_contest_submission(request: OmogenRequest, contest: Contest, user: Account, sub_id: int) -> HttpResponse:
    submission = get_submission_for_view(sub_id)
    problems = contest_problems(contest, problem_ids=[submission.problem_id])
    if not problems:
        raise Http404
    problem = problems[0].problem
    # TODO refactor into a can_see_submission
    if submission.account != user and not user.is_superuser:
        raise Http404
    problem_subtasks = get_subtask_scores(problem.current_version)
    submission_with_subtasks = SubmissionWithSubtasks(submission,
                                                      get_submission_subtask_scores(
                                                          list(submission.current_run.group_runs.all()),
                                                          subtasks=len(problem_subtasks)),
                                                      get_submission_subtask_groups(
                                                          list(submission.current_run.group_runs.all()),
                                                          subtasks=len(problem_subtasks)))
    args = ViewArgs(
        problem_with_scores=ProblemWithScores(problem=problem, subtask_scores=problem_subtasks),
        submission_with_subtasks=submission_with_subtasks,
        author=submission.account,
        language=Language(submission.language).display(),
        files={file: base64.b64decode(content).decode('UTF-8', errors='ignore') for file, content in
               submission.submission_files['files'].items()}
    )
    return render_template(request, 'submissions/view_submission.html', args)
