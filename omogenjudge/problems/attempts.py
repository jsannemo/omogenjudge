from omogenjudge.storage.models import Account, Problem, SubmissionRun, Verdict


def has_accepted(problem: Problem, account: Account) -> bool:
    return SubmissionRun.objects.filter(
        submission__problem=problem,
        submission__account=account,
        verdict=Verdict.AC,
    ).count() > 0
