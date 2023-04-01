from typing import Optional, Tuple

from django.core.cache import cache
from django.db.models import Prefetch, QuerySet

from omogenjudge.storage.models import Contest, ContestProblem, Problem, ProblemStatement, ProblemStatementFile
from omogenjudge.util.i18n import preferred_languages


class NoSuchLanguage(Exception):
    pass


def get_problem_for_view(short_name: str, *, language: Optional[str] = None) -> Tuple[Problem, ProblemStatement, list[str]]:
    problem = (
        Problem.objects
        .prefetch_related('statements')
        .prefetch_related(Prefetch('statement_files', queryset=ProblemStatementFile.objects.filter(attachment=1)))
        .select_related('current_version')
        .only('short_name', 'author', 'source', 'license', 'current_version__time_limit_ms',
              'current_version__memory_limit_kb')
        .get(short_name=short_name)
    )
    statements: dict[str, ProblemStatement] = {}
    for statement in problem.statements.all():
        statements[statement.language] = statement
    available_languages = sorted(list(statements.keys()))
    if language and language not in statements:
        raise NoSuchLanguage
    for lang in [language] + preferred_languages():
        if lang in statements:
            selected_language = lang
            break
        else:
            selected_language = available_languages[0]
    return problem, statements[selected_language], available_languages


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
