import collections
import dataclasses
import datetime
import math

from omogenjudge.problems.lookup import contest_problems
from omogenjudge.storage.models import Contest, ContestProblem, Status, Submission, Team, Verdict
from omogenjudge.submissions.lookup import list_contest_submissions
from omogenjudge.teams.lookup import contest_teams


@dataclasses.dataclass
class ScoreboardResult:
    tries: int = 0
    pending: bool = False
    accepted: bool = False
    time_str: str = ""


@dataclasses.dataclass
class ScoreboardTeam:
    name: str
    rank: int = 0
    score: int = 0
    time: int = 0
    results: list[ScoreboardResult] = dataclasses.field(default_factory=list)


@dataclasses.dataclass
class Scoreboard:
    problems: list[ContestProblem]
    teams: list[ScoreboardTeam]


def _submission_by_team(submissions: list[Submission], teams: list[Team]):
    account_to_team = {tm.account_id: team for team in teams for tm in team.teammember_set.all()}
    subs = collections.defaultdict(list)
    for s in submissions:
        subs[account_to_team[s.account_id]].append(s)
    return subs


def _aggregate_team_submissions(team: Team, subs: list[Submission], problem_ids: list[int],
                                start_time: datetime.datetime) -> ScoreboardTeam:
    team = ScoreboardTeam(name=team.team_name)
    results = {problem_id: ScoreboardResult() for problem_id in problem_ids}

    time = 0
    for sub in subs:
        result = results[sub.problem_id]
        if result.accepted:
            continue

        status = Status(sub.current_run.status)
        if status in (Status.QUEUED, Status.COMPILING, Status.COMPILE_ERROR, Status.JUDGE_ERROR):
            continue

        result.tries += 1
        if status == Status.RUNNING:
            result.pending = True
            continue

        verdict = Verdict(sub.current_run.verdict)
        if verdict == Verdict.AC:
            result.accepted = True
            team.score += 1
            seconds = math.floor((sub.date_created - start_time).total_seconds())
            time += seconds // 60
            result.time_str = "{:d}:{:02d}".format(seconds // 60, seconds % 60)

    # TODO: take contest penalty from contest
    team.time = time # + sum(20 * (result.tries - 1) if result.accepted else 0 for result in results.values())
    team.results = [results[problem_id] for problem_id in problem_ids]
    return team


def load_scoreboard(contest: Contest) -> Scoreboard:
    problems = contest_problems(contest)
    teams = contest_teams(contest)
    registered_user_ids = [tm.account_id for team in teams for tm in team.teammember_set.all()]
    problem_ids = [p.problem_id for p in problems]
    submissions = list_contest_submissions(registered_user_ids, problem_ids, contest)

    team_submissions = _submission_by_team(submissions, teams)
    scoreboard_teams = [_aggregate_team_submissions(team, team_submissions[team], problem_ids, contest.start_time) for
                        team in teams]

    scoreboard_teams.sort(key=lambda team: (-team.score, team.time))
    if scoreboard_teams:
        prev = (0, 0)
        atrank = 1
        for idx, team in enumerate(scoreboard_teams):
            now = (team.score, team.time)
            if prev != now:
                atrank = idx + 1
            team.rank = atrank
            prev = now

    return Scoreboard(
        problems=problems,
        teams=scoreboard_teams,
    )
