import dataclasses
import datetime
import math
from typing import Dict, List, Optional, Type

from django.utils import timezone

from omogenjudge.problems.lookup import contest_problems_with_grading
from omogenjudge.problems.testgroups import get_submission_subtask_scores, get_subtask_scores
from omogenjudge.storage.models import Contest, ContestProblem, ScoringType, Status, Submission, SubmissionGroupRun, \
    Team, Verdict
from omogenjudge.submissions.lookup import list_queue_submissions
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
    virtual: bool = False
    virtual_time: Optional[datetime.timedelta] = None
    practice: bool = False


@dataclasses.dataclass
class ScoreboardProblem:
    label: str
    problem: ContestProblem
    is_scoring: bool
    max_score: float = 0
    subtask_scores: list[float] = dataclasses.field(default_factory=list)


class ScoreboardMaker:
    def __init__(self, contest, problems: List[ContestProblem], teams: List[Team], *,
                 now: datetime.datetime,
                 at_time: Optional[datetime.timedelta] = None,
                 tiebreak=sum):
        self.now = now
        self.contest: Contest = contest
        self.problems = problems
        self.teams = teams
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
        self.upsolve_teams: List[ScoreboardTeam] = []
        self.user_to_rank: Dict[int, int] = {}
        self.best_user_result: Dict[int, ScoreboardTeam] = {}
        self._at_time = at_time

    def process_submissions(self, submissions: List[Submission]):
        account_to_team = {tm.account_id: team for team in self.teams for tm in team.teammember_set.all()}
        team_problem_submissions: Dict[Team, Dict[int, List[Submission]]] = {
            team: {problem.problem_id: [] for problem in self.problems} for team in
            self.teams}
        for s in submissions:
            team_problem_submissions[account_to_team[s.account_id]][s.problem_id].append(s)

        self.scoreboard_teams = [t for team in self.teams for t in
                                 self._process_team(team, team_problem_submissions[team])]
        for team in self.scoreboard_teams:
            team.total_score = self._round(sum(p.problem_score for p in team.results))
        self._sort_teams()

    def _process_team(self, team: Team, submissions: Dict[int, List[Submission]]) -> list[
        ScoreboardTeam]:
        problem_results = []
        upsolved_results = []
        start_time = team.contest_start_time if team.contest_start_time else self.contest.start_time
        for scoreboard_problem in self.scoreboard_problems:
            all_submissions = submissions[scoreboard_problem.problem.problem_id]
            contest_submissions = self._contest_submissions(team, all_submissions)
            problem_results.append(self._process_problem(contest_submissions, scoreboard_problem, start_time))
            upsolved_results.append(self._process_problem(all_submissions, scoreboard_problem, start_time))

        teams = []
        if not team.practice or team.contest_start_time:
            sc_team = ScoreboardTeam(
                team=team,
                results=problem_results,
                tiebreak=self._tiebreak_aggregate(res.tiebreak for res in problem_results),
                virtual=team.practice and team.contest_start_time is not None,
                practice=False,
            )
            if team.contest_start_time:
                elapsed = self.now - team.contest_start_time
                if (not self._at_time or elapsed <= self._at_time) and elapsed <= self.contest.duration:
                    sc_team.virtual_time = elapsed
            teams.append(sc_team)
        if not self._at_time and (not teams or problem_results != upsolved_results):
            upsolve_team = ScoreboardTeam(
                team=team,
                results=upsolved_results,
                tiebreak=self._tiebreak_aggregate(res.tiebreak for res in upsolved_results),
                practice=True,
            )
            teams.append(upsolve_team)
        return teams

    def _process_problem(self, submissions: List[Submission], problem: ScoreboardProblem,
                         start_time: Optional[datetime.datetime]) -> ProblemResult:
        raise NotImplementedError

    def _sort_teams(self):
        self.scoreboard_teams.sort(key=lambda t: (t.practice, self._team_sort_key(t)))
        if self.scoreboard_teams:
            prev = self._team_sort_key(self.scoreboard_teams[0])
            at_rank = 1
            for idx, team in enumerate(self.scoreboard_teams):
                if team.practice:
                    continue
                now = self._team_sort_key(team)
                if prev != now:
                    at_rank = idx + 1
                team.rank = at_rank
                prev = now

        for i, team in enumerate(self.scoreboard_teams):
            if not team.practice:
                for user in team.team.teammember_set.all():
                    self.user_to_rank[user.account_id] = i
            for user in team.team.teammember_set.all():
                self.best_user_result[user.account_id] = team

    def _round(self, num):
        return round(num, ndigits=2)

    def _team_sort_key(self, team: ScoreboardTeam):
        return -team.total_score, team.tiebreak

    def format_tiebreak(self, minutes: float) -> str:
        minutes = int(minutes)
        neg = minutes < 0
        if neg: minutes = -minutes
        return ("-" if neg else "") + "{:02d}:{:02d}".format(minutes // 60, minutes % 60)

    def max_score(self) -> float:
        return self._round(sum(self._round(p.max_score) for p in self.scoreboard_problems))

    def _contest_submissions(self, team: Team, submissions: list[Submission]) -> list[Submission]:
        if team.practice and not team.contest_start_time:
            return []
        current_duration = self._at_time or self.contest.duration
        start_time = team.contest_start_time or self.contest.start_time
        end_time = team.contest_start_time + current_duration if team.contest_start_time else (
            self.contest.start_time + current_duration if self.contest.start_time else None)
        contest_submissions = []
        for submission in submissions:
            if start_time and end_time and start_time <= submission.date_created <= end_time:
                contest_submissions.append(submission)
        return contest_submissions

    @property
    def has_penalty(self):
        return NotImplementedError


def _subtask_runs(group_runs: List[SubmissionGroupRun], subtasks: int):
    # First run is sample; assume all other subtasks have depth 1
    return group_runs[1:subtasks + 1]


class Scoring(ScoreboardMaker):
    def __init__(self, contest: Contest, problems: List[ContestProblem], teams: List[Team], **kwargs):
        super(Scoring, self).__init__(contest, problems, teams, tiebreak=max, **kwargs)
        self.problem_subtasks = {}
        for problem, scoreboard_problem in zip(self.problems, self.scoreboard_problems):
            subtask_scores = get_subtask_scores(problem.problem.current_version)
            self.problem_subtasks[problem.problem_id] = subtask_scores
            scoreboard_problem.subtask_scores = subtask_scores
            scoreboard_problem.max_score = sum(subtask_scores)

    def _process_problem(self, submissions: List[Submission], scoreboard_problem: ScoreboardProblem,
                         start_time: Optional[datetime.datetime]) -> ProblemResult:
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

            if start_time and score_after > score_now:
                problem_result.tiebreak = (sub.date_created - start_time).total_seconds() // 60

            # Problem is maxed; additional problems no longer count as tries
            if problem_result.subtask_scores == scoreboard_problem.subtask_scores:
                break
        problem_result.problem_score = self._round(sum(problem_result.subtask_scores))
        return problem_result

    def _team_sort_key(self, team):
        return -team.total_score, team.tiebreak

    @property
    def has_penalty(self):
        return False


class BinaryWithPenalty(ScoreboardMaker):

    def _process_problem(self, submissions: List[Submission], scoreboard_problem: ScoreboardProblem,
                         start_time: Optional[datetime.datetime]) -> ProblemResult:
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
                if start_time:
                    seconds = math.floor((sub.date_created - start_time).total_seconds())
                    problem_result.tiebreak = failures * self.contest.try_penalty + seconds // 60
            else:
                failures += 1
        return problem_result

    @property
    def has_penalty(self):
        return True


_SCOREBOARDS: Dict[ScoringType, Type[ScoreboardMaker]] = {
    ScoringType.BINARY_WITH_PENALTY: BinaryWithPenalty,
    ScoringType.SCORING: Scoring,
}


def load_scoreboard(contest: Contest, *, now: Optional[datetime.datetime] = None,
                    at_time: Optional[datetime.timedelta] = None) -> ScoreboardMaker:
    if not now:
        now = timezone.now()
    problems = list(contest_problems_with_grading(contest))
    teams = list(contest_teams(contest))
    registered_user_ids = [tm.account_id for team in teams for tm in team.teammember_set.all()]
    submissions = list(
        list_queue_submissions(registered_user_ids, [problem.problem_id for problem in problems], ascending=True))

    scoreboard = _SCOREBOARDS[ScoringType(contest.scoring_type)](contest, problems, teams, now=now, at_time=at_time)
    scoreboard.process_submissions(submissions)
    return scoreboard
