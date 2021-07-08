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


def list_contest_submissions(user_ids: list[int], problem_ids: list[int], contest: Contest) -> QuerySet[Submission]:
    qs = Submission.objects.filter(
        account_id__in=user_ids,
        problem_id__in=problem_ids,
    ).select_related('current_run') \
        .prefetch_related('current_run__group_runs') \
        .order_by('-submission_id')

    if contest.start_time:
        qs = qs.filter(
            date_created__gte=contest.start_time,
            date_created__lt=contest.start_time + contest.duration,
        )

    return qs


def all_submissions_for_queue() -> QuerySet[Submission]:
    return Submission.objects.select_related('current_run').select_related('problem').prefetch_related(
        Prefetch(
            'problem__statements',
            ProblemStatement.objects.all().only('problem_id', 'language', 'title'))).order_by('-submission_id').all()
