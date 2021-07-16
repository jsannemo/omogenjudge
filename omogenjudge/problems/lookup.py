from typing import Optional, Tuple

from django.db.models import Prefetch

from omogenjudge.storage.models import Contest, ContestProblem, Problem, ProblemStatement, ProblemStatementFile


class NoSuchLanguage(Exception):
    pass


def get_problem_for_view(short_name: str, *, language: Optional[str] = None) -> Tuple[Problem, ProblemStatement]:
    problem = (
        Problem.objects.prefetch_related('problemstatement_set')
            .select_related('current_version')
            .only('short_name', 'author', 'source', 'license', 'current_version__time_limit_ms',
                  'current_version__memory_limit_kb')
            .get(short_name=short_name)
    )
    statements = {}
    for statement in problem.problemstatement_set.all():
        if statement.language == language:
            return problem, statement
        statements[statement.language] = statement
    if language:
        raise NoSuchLanguage
    # TODO: look at user lang
    for lang in ["sv", "en"]:
        if lang in statements:
            return problem, statements[lang]
    return problem, problem.problemstatement_set[0]


def problem_by_name(short_name: str) -> Problem:
    return Problem.objects.get(short_name=short_name)


def find_statement_file(problem_short_name: str, path: str) -> ProblemStatementFile:
    return ProblemStatementFile.objects.get(
        problem__short_name=problem_short_name,
        file_path=path,
    )


def problems_with_name_query():
    return Problem.objects.prefetch_related(
        Prefetch(
            'problemstatement_set',
            ProblemStatement.objects.all().only('problem_id', 'language', 'title'),
            to_attr='statements'))


def list_public_problems() -> list[Problem]:
    return problems_with_name_query().all()


def contest_problems(contest: Contest) -> list[ContestProblem]:
    return contest.contestproblem_set.select_related('problem').prefetch_related(
        Prefetch(
            'problem__problemstatement_set',
            ProblemStatement.objects.all().only('problem_id', 'language', 'title'),
            to_attr='statements')).order_by('label').all()
