import dataclasses
from typing import Dict

from django.contrib.auth.decorators import login_required
from django.http import HttpResponse

from omogenjudge.frontend.decorators import only_started_contests, requires_contest, requires_user
from omogenjudge.frontend.submissions.view_submission import SubmissionWithSubtasks, ProblemWithScores
from omogenjudge.problems.lookup import contest_problems_with_grading
from omogenjudge.problems.testgroups import get_subtask_scores, get_submission_subtask_scores
from omogenjudge.storage.models import Contest, Account
from omogenjudge.submissions.lookup import list_contest_submissions
from omogenjudge.util.django_types import OmogenRequest
from omogenjudge.util.templates import render_template


@dataclasses.dataclass
class QueueArgs:
    submissions: list[SubmissionWithSubtasks]
    problems: Dict[int, ProblemWithScores]


@requires_user
def submission_queue(_request: OmogenRequest, user: Account) -> HttpResponse:
    raise NotImplementedError


@requires_user
@requires_contest
@only_started_contests
def my_submissions(request: OmogenRequest, user: Account, contest: Contest) -> HttpResponse:
    # TODO: add page for non-contests
    problems = [
        ProblemWithScores(problem=cp.problem, subtask_scores=get_subtask_scores(cp.problem.current_version))
        for cp in contest_problems_with_grading(contest)]

    problem_map = {p.problem.problem_id: p for p in problems}

    submissions = list_contest_submissions([user.account_id], list(problem_map.keys()), contest)
    submissions_with_subtasks = [
        SubmissionWithSubtasks(submission, get_submission_subtask_scores(list(submission.current_run.group_runs.all()),
                                                                         subtasks=len(problem_map[
                                                                                          submission.problem_id].subtask_scores)))
        for
        submission in submissions]

    return render_template(request, 'submissions/my.html', QueueArgs(submissions_with_subtasks, problem_map))
