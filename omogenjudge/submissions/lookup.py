from typing import Optional

from django.db.models import Prefetch, QuerySet

from omogenjudge.storage.models import Account, Contest, Problem, ProblemStatement, Submission


def get_submission_for_view(sub_id: int) -> Submission:
    return Submission.objects.select_related('current_run').select_related('problem') \
        .prefetch_related('current_run__group_runs') \
        .prefetch_related(
        Prefetch(
            'problem__statements',
            ProblemStatement.objects.all().only('problem_id', 'language', 'title'))).get(submission_id=sub_id)


def list_account_problem_submissions(*,
                                     account: Account,
                                     problem: Problem,
                                     limit: Optional[int] = None) -> QuerySet[Submission]:
    qs = Submission.objects.select_related('current_run') \
        .prefetch_related('current_run__group_runs') \
        .filter(account=account, problem=problem) \
        .order_by('-submission_id')

    qs = qs.all()
    if limit is not None:
        qs = qs[:limit]
    return qs


def list_queue_submissions(user_ids: list[int] | None, problem_ids: list[int], *, ascending: bool=False) -> \
        QuerySet[Submission]:
    filters = {}
    if user_ids is not None:
        filters['account_id__in'] = user_ids
    if problem_ids is not None:
        filters['problem_id__in'] = problem_ids
    qs = Submission.objects.filter(**filters).select_related('current_run') \
        .prefetch_related('current_run__group_runs') \
        .order_by('submission_id' if ascending else '-submission_id')
    return qs
