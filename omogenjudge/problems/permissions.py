from omogenjudge.contests.permissions import can_view_contest_problems
from omogenjudge.problems.lookup import contest_problems
from omogenjudge.storage.models import Problem
from omogenjudge.util.request_global import current_request


def can_view_problem(problem: Problem) -> bool:


    # TODO: shot and submit mode
    #        cproblems = contest_problems(contest)
    #        for idx, prob in enumerate(cproblems):
    #            if prob.problem == problem:
    #                break
    #            if not has_accepted(cproblems[idx].problem, request.user):
    #                return redirect('problem', short_name=cproblems[idx].problem.short_name)

    contest = current_request().contest
    if not contest:
        return False
    if not any(problem.problem_id == contest_problem.problem_id for contest_problem in contest_problems(contest)):
        return False
    return can_view_contest_problems()
