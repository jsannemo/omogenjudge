from omogenjudge.contests.permissions import team_can_view_problems
from omogenjudge.problems.lookup import contest_problems
from omogenjudge.storage.models import Contest, Problem
from omogenjudge.teams.lookup import contest_team_for_user
from omogenjudge.util.request_global import current_request


def can_view_problem(problem: Problem) -> bool:
    request = current_request()
    user = request.user
    if user.is_superuser:
        return True
    contest = request.contest
    if not contest:
        return False
    cproblems = contest_problems(contest)
    if not any(problem.problem_id == contest_problem.problem_id for contest_problem in cproblems):
        return False
    return team_can_view_problems(contest, contest_team_for_user(contest, user))


def can_submit_in_contest(contest: Contest) -> bool:
    request = current_request()
    user = request.user
    if user.is_superuser:
        return True
    return contest_team_for_user(contest, user) is not None
