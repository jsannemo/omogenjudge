from typing import Optional, Tuple

from django.core.cache import cache
from django.db.models import Prefetch, QuerySet

from omogenjudge.storage.models import Contest, ContestProblem, Problem, ProblemStatement, ProblemStatementFile


class NoSuchLanguage(Exception):
    pass


def get_problem_for_view(short_name: str, *, language: Optional[str] = None) -> Tuple[Problem, ProblemStatement]:
    problem = (
        Problem.objects.prefetch_related('statements')
        .select_related('current_version')
        .only('short_name', 'author', 'source', 'license', 'current_version__time_limit_ms',
              'current_version__memory_limit_kb')
        .get(short_name=short_name)
    )
    statements = {}
    for statement in problem.statements.all():
        if statement.language == language:
            return problem, statement
        statements[statement.language] = statement
    if language:
        raise NoSuchLanguage
    # TODO: look at user lang
    for lang in ["sv", "en"]:
        if lang in statements:
            return problem, statements[lang]
    return problem, statements[0]


def problem_by_name(short_name: str) -> Problem:
    return Problem.objects.get(short_name=short_name)


def find_statement_file(problem_short_name: str, path: str) -> ProblemStatementFile:
    return ProblemStatementFile.objects.get(
        problem__short_name=problem_short_name,
        file_path=path,
    )


def _statements_with_title():
    return ProblemStatement.objects.all().only('problem_id', 'language', 'title')


def problems_with_name_query():
    return Problem.objects.prefetch_related(Prefetch('statements', _statements_with_title()))


def list_public_problems() -> list[Problem]:
    return problems_with_name_query().all()


def contest_problems_query(contest: Contest, *, problem_ids: Optional[list[int]] = None) -> QuerySet[ContestProblem]:
    qs = contest.contestproblem_set.select_related('problem') \
        .prefetch_related(Prefetch('problem__statements', _statements_with_title())) \
        .order_by('label', 'problem__short_name') \
        .all()
    if problem_ids:
        qs = qs.filter(problem_id__in=problem_ids)
    return qs


def contest_problems(contest: Contest, *, problem_ids: Optional[list[int]] = None) -> QuerySet[ContestProblem]:
    return contest_problems_query(contest, problem_ids=problem_ids).all()


def contest_problems_with_grading(contest: Contest, *, problem_ids: Optional[list[int]] = None) -> QuerySet[
    ContestProblem]:
    return contest_problems_query(contest, problem_ids=problem_ids) \
        .select_related('problem__current_version') \
        .prefetch_related('problem__current_version__testgroups').all()
