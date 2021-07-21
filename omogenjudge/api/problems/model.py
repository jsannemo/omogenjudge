import dataclasses

from omogenjudge.storage.models import ProblemStatement, ProblemVersion


@dataclasses.dataclass
class ApiStatement:
    title: str
    html: str
    license: str
    authors: list[str]

    @staticmethod
    def from_db_statement(statement: ProblemStatement):
        return ApiStatement(
            title=statement.title,
            html=statement.html,
            license=statement.problem.license,
            authors=statement.problem.author,
        )


@dataclasses.dataclass
class ApiProblemLimits:
    timelim_ms: int
    memlim_kb: int

    @staticmethod
    def from_db_version(version: ProblemVersion):
        return ApiProblemLimits(
            timelim_ms=version.time_limit_ms,
            memlim_kb=version.memory_limit_kb,
        )
