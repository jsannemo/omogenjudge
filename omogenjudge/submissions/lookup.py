from typing import Optional

from django.db.models import Prefetch

from omogenjudge.storage.models import Account, Contest, Problem, ProblemStatement, Submission


def get_submission_for_view(sub_id: int):
    return Submission.objects.select_related('current_run').select_related('problem').prefetch_related(
        Prefetch(
            'problem__problemstatement_set',
            ProblemStatement.objects.all().only('problem_id', 'language', 'title'),
            to_attr='statements')).get(submission_id=sub_id)


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


def list_contest_submissions(user_ids: list[int], problem_ids: list[int], contest: Contest):
    return Submission.objects.filter(
        account_id__in=user_ids,
        problem_id__in=problem_ids,
        date_created__gte=contest.start_time,
        date_created__lt=contest.start_time + contest.duration
    ).select_related('current_run').order_by('date_created').all()


def all_submissions_for_queue() -> list[Submission]:
    return Submission.objects.select_related('current_run').select_related('problem').prefetch_related(
        Prefetch(
            'problem__problemstatement_set',
            ProblemStatement.objects.all().only('problem_id', 'language', 'title'),
            to_attr='statements')).order_by('-submission_id').all()
