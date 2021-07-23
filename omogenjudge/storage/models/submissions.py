import dataclasses
import enum

from django.db import models

from .accounts import Account
from .langauges import Language
from .problems import Problem, ProblemTestcase, ProblemTestgroup, ProblemVersion
from ...util import django_fields, serialization
from ...util.django_fields import PrefetchIDMixin
from ...util.enums import EnumChoices


@dataclasses.dataclass(frozen=True)
class SubmissionFiles:
    files: dict[str, str]


class SubmissionFilesDecoder(serialization.DataclassJsonDecoder[SubmissionFiles]):
    def __init__(self):
        super().__init__(SubmissionFiles)


class Submission(PrefetchIDMixin, models.Model):
    submission_id = models.AutoField(primary_key=True)
    account = models.ForeignKey(Account, models.CASCADE)
    problem = models.ForeignKey(Problem, models.CASCADE)
    language = models.TextField(choices=Language.as_choices())
    current_run = models.ForeignKey('SubmissionRun', models.RESTRICT, db_column='current_run', related_name='+')
    date_created = models.DateTimeField(auto_now_add=True)
    submission_files = models.JSONField(
        encoder=serialization.DataclassJsonEncoder,
        decoder=SubmissionFilesDecoder,
    )

    def __str__(self):
        return f'Submission {self.submission_id} by {self.account_id} for {self.problem_id}'

    class Meta:
        db_table = 'submission'


class Verdict(EnumChoices['Verdict'], enum.Enum):
    UNJUDGED = 'unjudged'
    AC = 'accepted'
    WA = 'wrong answer'
    TLE = 'time limit exceeded'
    RTE = 'run-time error'


class SubmissionCaseRun(models.Model):
    submission_run = models.ForeignKey('SubmissionRun', models.CASCADE)
    problem_testcase = models.ForeignKey(ProblemTestcase, models.RESTRICT, related_name='+')
    date_created = models.DateTimeField()
    time_usage_ms = models.IntegerField()
    score = models.FloatField()
    verdict = models.TextField(choices=Verdict.as_choices())  # This field type is a guess.

    class Meta:
        db_table = 'submission_case_run'



class SubmissionGroupRun(models.Model):
    submission_run = models.ForeignKey('SubmissionRun', models.CASCADE)
    problem_testgroup = models.ForeignKey(ProblemTestgroup, models.RESTRICT, related_name='+')
    date_created = models.DateTimeField()
    time_usage_ms = models.IntegerField()
    score = models.FloatField()
    verdict = models.TextField(choices=Verdict.as_choices())

    class Meta:
        db_table = 'submission_group_run'


class Status(EnumChoices['Status'], enum.Enum):
    QUEUED = 'queued'
    COMPILING = 'compiling'
    RUNNING = 'running'
    COMPILE_ERROR = 'compile error'
    JUDGE_ERROR = 'judging error'
    DONE = 'done'


class SubmissionRun(models.Model):
    submission_run_id = models.AutoField(primary_key=True)
    submission = models.ForeignKey(Submission, models.CASCADE)
    problem_version = models.ForeignKey(ProblemVersion, models.CASCADE)
    date_created = models.DateTimeField(auto_now_add=True)
    status = models.TextField(choices=Status.as_choices())
    verdict = models.TextField(choices=Verdict.as_choices())
    time_usage_ms = models.IntegerField(null=True)
    score = models.FloatField(null=True)
    compile_error = django_fields.TextField(null=True)

    class Meta:
        db_table = 'submission_run'