import dataclasses
import datetime
import math
from typing import List, Dict, TypeVar, Type, Optional

from omogenjudge.problems.lookup import contest_problems_with_grading
from omogenjudge.problems.testgroups import get_subtask_scores, get_submission_subtask_scores
from omogenjudge.storage.models import Contest, ContestProblem, Status, Submission, Team, Verdict, ScoringType, \
    SubmissionGroupRun
from omogenjudge.submissions.lookup import list_contest_submissions
from omogenjudge.teams.lookup import contest_teams


@dataclasses.dataclass
class ProblemResult:
    pending: int = 0
    tries: int = 0
    accepted: bool = False
    problem_score: float = 0
    subtask_scores: list[float] = dataclasses.field(default_factory=list)
    tiebreak: float = 0


@dataclasses.dataclass
class TeamResult:
    total_score: float = 0
    tiebreak: float = 0


@dataclasses.dataclass
class ScoreboardTeam:
    team: Team
    rank: int = 0
    total_score: float = 0
    tiebreak: float = 0
    results: list[ProblemResult] = dataclasses.field(default_factory=list)


@dataclasses.dataclass
class ScoreboardProblem:
    label: str
    problem: ContestProblem
    is_scoring: bool
    max_score: float = 0
    subtask_scores: list[float] = dataclasses.field(default_factory=list)


class ScoreboardMaker:
    def __init__(self, problems: List[ContestProblem], teams: List[Team], start_time: Optional[datetime.datetime], *,
                 tiebreak=sum):
        self.problems = problems
        self.teams = teams
        self._start_time = start_time
        self._tiebreak_aggregate = tiebreak

        self.scoreboard_problems = [
            ScoreboardProblem(
                label=p.label,
                problem=p,
                is_scoring=p.problem.current_version.scoring,
            )
            for p in problems
        ]
        self.scoreboard_teams: List[ScoreboardTeam] = []
        self.user_to_rank: Dict[int, int] = {}

    def process_submissions(self, submissions: List[Submission]):
        account_to_team = {tm.account_id: team for team in self.teams for tm in team.teammember_set.all()}
        team_problem_submissions: Dict[Team, Dict[int, List[Submission]]] = {
            team: {problem.problem_id: [] for problem in self.problems} for team in
            self.teams}
        for s in submissions:
            team_problem_submissions[account_to_team[s.account_id]][s.problem_id].append(s)

        self.scoreboard_teams = [self._process_team(team, team_problem_submissions[team]) for team in self.teams]
        for team in self.scoreboard_teams:
            team.total_score = self._round(sum(p.problem_score for p in team.results))
        self._sort_teams()

    def _process_team(self, team: Team, submissions: Dict[int, List[Submission]]) -> ScoreboardTeam:
        problem_results = [
            self._process_problem(submissions[scoreboard_problem.problem.problem_id], scoreboard_problem) for
            scoreboard_problem in
            self.scoreboard_problems]
        return ScoreboardTeam(
            team=team,
            results=problem_results,
            tiebreak=self._tiebreak_aggregate(res.tiebreak for res in problem_results)
        )

    def _process_problem(self, submissions: List[Submission], problem: ScoreboardProblem) -> ProblemResult:
        raise NotImplementedError

    def _sort_teams(self):
        self.scoreboard_teams.sort(key=lambda t: self._team_sort_key(t))
        if self.scoreboard_teams:
            prev = self._team_sort_key(self.scoreboard_teams[0])
            at_rank = 1
            for idx, team in enumerate(self.scoreboard_teams):
                now = self._team_sort_key(team)
                if prev != now:
                    at_rank = idx + 1
                team.rank = at_rank
                prev = now

        for i, team in enumerate(self.scoreboard_teams):
            for user in team.team.teammember_set.all():
                self.user_to_rank[user.account_id] = i

    def _round(self, num):
        return round(num, ndigits=2)

    def _team_sort_key(self, team: ScoreboardTeam):
        return -team.total_score, team.tiebreak

    def format_tiebreak(self, minutes: float) -> str:
        minutes = int(minutes)
        return "{:02d}:{:02d}".format(minutes // 60, minutes % 60)

    def max_score(self) -> float:
        return self._round(sum(self._round(p.max_score) for p in self.scoreboard_problems))


def _subtask_runs(group_runs: List[SubmissionGroupRun], subtasks: int):
    # First run is sample; assume all other subtasks have depth 1
    return group_runs[1:subtasks + 1]


class Scoring(ScoreboardMaker):
    def __init__(self, problems: List[ContestProblem], teams: List[Team], start_time: datetime.datetime, **kwargs):
        super(Scoring, self).__init__(problems, teams, start_time, tiebreak=max, **kwargs)
        self.problem_subtasks = {}
        for problem, scoreboard_problem in zip(self.problems, self.scoreboard_problems):
            subtask_scores = get_subtask_scores(problem.problem.current_version)
            self.problem_subtasks[problem.problem_id] = subtask_scores
            scoreboard_problem.subtask_scores = subtask_scores
            scoreboard_problem.max_score = sum(subtask_scores)

    def _process_problem(self, submissions: List[Submission], scoreboard_problem: ScoreboardProblem) -> ProblemResult:
        problem_result = ProblemResult(
            subtask_scores=[0] * len(scoreboard_problem.subtask_scores)
        )
        for sub in submissions:
            run = sub.current_run
            status = Status(run.status)
            if status in [Status.RUNNING, Status.QUEUED, Status.COMPILING]:
                problem_result.pending += 1
                continue
            if status in [Status.JUDGE_ERROR, Status.COMPILE_ERROR]:
                continue
            assert status == Status.DONE

            problem_result.tries += 1

            score_now = sum(problem_result.subtask_scores)
            for i, group_run_score in enumerate(
                    get_submission_subtask_scores(list(run.group_runs.all()), len(scoreboard_problem.subtask_scores))):
                problem_result.subtask_scores[i] = max(problem_result.subtask_scores[i], group_run_score)
            score_after = sum(problem_result.subtask_scores)

            if self._start_time and score_after > score_now:
                problem_result.tiebreak = (sub.date_created - self._start_time).total_seconds() // 60

            # Problem is maxed; additional problems no longer count as tries
            if problem_result.subtask_scores == scoreboard_problem.subtask_scores:
                break
        problem_result.problem_score = self._round(sum(problem_result.subtask_scores))
        return problem_result

    def _team_sort_key(self, team):
        return -team.total_score, team.tiebreak


class BinaryWithPenalty(ScoreboardMaker):

    def _process_problem(self, submissions: List[Submission], problem: ScoreboardProblem) -> ProblemResult:
        problem_result = ProblemResult()
        failures = 0
        for sub in submissions:
            run = sub.current_run
            status = Status(run.status)
            if status in [Status.RUNNING, Status.QUEUED, Status.COMPILING]:
                problem_result.pending += 1
                continue
            if status in [Status.JUDGE_ERROR, Status.COMPILE_ERROR]:
                continue
            assert status == Status.DONE
            if problem_result.accepted:
                break

            problem_result.tries += 1

            verdict = Verdict(sub.current_run.verdict)
            if verdict == Verdict.AC:
                problem_result.accepted = True
                problem_result.problem_score += 1
                if self._start_time:
                    seconds = math.floor((sub.date_created - self._start_time).total_seconds())
                    problem_result.tiebreak = failures * 20 + seconds // 60
        return problem_result


_SCOREBOARDS: Dict[ScoringType, Type[ScoreboardMaker]] = {
    ScoringType.BINARY_WITH_PENALTY: BinaryWithPenalty,
    ScoringType.SCORING: Scoring,
}


def load_scoreboard(contest: Contest, *, teams=None) -> ScoreboardMaker:
    problems = list(contest_problems_with_grading(contest))
    if not teams:
        teams = list(contest_teams(contest))
    registered_user_ids = [tm.account_id for team in teams for tm in team.teammember_set.all()]
    submissions = list(list_contest_submissions(registered_user_ids, [problem.problem_id for problem in problems], contest))

    scoreboard = _SCOREBOARDS[ScoringType(contest.scoring_type)](problems, teams, contest.start_time)
    scoreboard.process_submissions(submissions)
    return scoreboard
