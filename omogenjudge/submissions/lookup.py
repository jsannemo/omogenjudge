from typing import Optional

from omogenjudge.storage.models import Account, Problem, Submission


def list_account_problem_submissions(*,
                                     account: Account,
                                     problem: Problem,
                                     limit: Optional[int] = None) -> list[Submission]:
    qs = (Submission.objects.prefetch_related('current_run')
          .filter(account=account, problem=problem)
          .order_by('-submission_id')
          .all()
          )
    if limit != None:
        qs = qs[:limit]
    return qs
