import math
from typing import List, Optional

from omogenjudge.storage.models import ProblemVersion, ProblemTestgroup, SubmissionGroupRun


def _none_score_to_inf(score: Optional[float]) -> float:
    if score is None:
        return math.inf
    return score


def get_subtask_scores(problem_version: ProblemVersion) -> List[float]:
    testgroups = problem_version.testgroups.all()
    secret_group = None
    subtasks: List[ProblemTestgroup] = []
    for group in testgroups:
        name = group.testgroup_name
        if name == "data/secret":
            secret_group = group
        if name.startswith("data/secret/"):
            subtasks.append(group)
    # Fallback for scoring problems that have no subtasks
    if not subtasks:
        assert secret_group, f"No secret group among {testgroups}"
        return [_none_score_to_inf(secret_group.max_score)]
    subtasks.sort(key=lambda key: key.testgroup_name)
    return [_none_score_to_inf(group.max_score) for group in subtasks]


def get_submission_subtask_scores(group_runs: List[SubmissionGroupRun], subtasks: int) -> List[float]:
    return [group.score for group in get_submission_subtask_groups(group_runs, subtasks)]


def get_submission_subtask_groups(group_runs: List[SubmissionGroupRun], subtasks: int) -> List[SubmissionGroupRun]:
    # First run is sample; assume all other subtasks have depth 1
    # 'secret' is judged last, so this also works where there are no subtasks (i.e. subtasks = 1, the secret group)
    return [group for group in group_runs[1:min(subtasks + 1, len(group_runs))]]
