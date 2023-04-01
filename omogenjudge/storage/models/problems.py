import dataclasses
import enum

from django.contrib.postgres.fields import ArrayField
from django.db import models
from django.utils.functional import cached_property

from omogenjudge.util import django_fields, serialization
from omogenjudge.util.django_fields import PrefetchIDMixin, EnumField, TextField, StrEnum


class License(StrEnum):
    PUBLIC_DOMAIN = 'public domain'
    CC0 = 'cc0'
    CC_BY = 'cc by'
    CC_BY_SA = 'cc by-sa'
    EDUCATIONAL = 'educational'
    PERMISSION = 'permission'

    def display(self):
        return _LICENSE_NAME_AND_URL[self][0]

    def url(self):
        return _LICENSE_NAME_AND_URL[self][1]


_LICENSE_NAME_AND_URL = {
    License.PUBLIC_DOMAIN: ("Public Domain", "http://creativecommons.org/about/pdm"),
    License.CC0: ("CC0", "https://creativecommons.org/publicdomain/zero/1.0/"),
    License.CC_BY: ("CC BY", "https://creativecommons.org/licenses/by/3.0/"),
    License.CC_BY_SA: ("CC BY-SA", "https://creativecommons.org/licenses/by-sa/3.0/"),
    License.EDUCATIONAL: ("Used with permission for educational purposes", None),
    License.PERMISSION: ("Used with permission", None),
}


class Problem(PrefetchIDMixin, models.Model):
    problem_id = models.AutoField(primary_key=True)
    short_name = django_fields.TextField(unique=True)
    author = ArrayField(django_fields.TextField())
    source = django_fields.TextField()
    license = EnumField(enum_type=License)
    current_version = models.ForeignKey('ProblemVersion', models.RESTRICT, related_name='+')

    @cached_property
    def titles_by_language(self):
        return {
            s.language: s.title
            for s in self.statements.all()
        }

    def __str__(self):
        return self.short_name

    class Meta:
        db_table = 'problem'


class ProblemOutputValidator(models.Model):
    problem_output_validator_id = models.AutoField(primary_key=True)
    run_command = ArrayField(TextField())
    validator_zip = models.ForeignKey('StoredFile', models.RESTRICT, related_name='+')
    scoring_validator = models.BooleanField()

    class Meta:
        db_table = 'problem_output_validator'


class ProblemGrader(models.Model):
    problem_grader_id = models.AutoField(primary_key=True)
    run_command = ArrayField(TextField())
    grader_zip = models.ForeignKey('StoredFile', models.RESTRICT, related_name='+')

    class Meta:
        db_table = 'problem_grader'


class ProblemStatement(models.Model):
    problem_statement_id = models.AutoField(primary_key=True)
    problem = models.ForeignKey(Problem, models.CASCADE, related_name='statements')
    language = django_fields.TextField()
    title = django_fields.TextField()
    html = django_fields.TextField()

    class Meta:
        db_table = 'problem_statement'
        unique_together = (('problem', 'language'),)


class ProblemStatementFile(models.Model):
    problem_statement_file_id = models.AutoField(primary_key=True)
    problem = models.ForeignKey(Problem, models.CASCADE, related_name='statement_files')
    file_path = django_fields.TextField()
    statement_file = models.ForeignKey('StoredFile', models.RESTRICT, db_column='statement_file_hash', related_name='+')
    attachment = models.BooleanField()

    class Meta:
        db_table = 'problem_statement_file'
        unique_together = (('problem', 'file_path'),)


class ProblemTestcase(models.Model):
    problem_testcase_id = models.AutoField(primary_key=True)
    problem_testgroup = models.ForeignKey('ProblemTestgroup', models.CASCADE)
    testcase_name = django_fields.TextField()
    input_file = models.ForeignKey('StoredFile', models.RESTRICT, db_column='input_file_hash', related_name='+')
    output_file = models.ForeignKey('StoredFile', models.RESTRICT, db_column='output_file_hash', related_name='+')

    class Meta:
        db_table = 'problem_testcase'


class ScoringMode(enum.Enum):
    SUM = 'sum'
    AVG = 'avg'
    MIN = 'min'
    MAX = 'max'


class VerdictMode(enum.Enum):
    WORST_ERROR = 'worst_error'
    FIRST_ERROR = 'first_error'
    ALWAYS_ACCEPTED = 'always_accept'


class ProblemTestgroup(models.Model):
    problem_testgroup_id = models.AutoField(primary_key=True)
    parent = models.ForeignKey('self', models.CASCADE, null=True, related_name='children')
    problem_version = models.ForeignKey('ProblemVersion', models.CASCADE, related_name='testgroups')
    testgroup_name = django_fields.TextField()

    # Test group config
    min_score = models.FloatField(null=True)
    max_score = models.FloatField(null=True)
    accept_score = models.FloatField(null=True)
    reject_score = models.FloatField(null=True)
    break_on_reject = models.BooleanField()
    output_validator_flags = ArrayField(django_fields.TextField())

    # Grader flags
    scoring_mode = EnumField(enum_type=ScoringMode)
    verdict_mode = EnumField(enum_type=VerdictMode)
    accept_if_any_accepted = models.BooleanField()
    ignore_sample = models.BooleanField()
    grader_flags = ArrayField(django_fields.TextField())
    custom_grading = models.BooleanField()

    def __str__(self):
        return f"Testgroup {self.testgroup_name} for {self.problem_testgroup_id}"

    class Meta:
        db_table = 'problem_testgroup'


@dataclasses.dataclass(frozen=True)
class IncludedFiles:
    files_by_language: dict[str, dict[str, str]]


class IncludedFilesDecoder(serialization.DataclassJsonDecoder[IncludedFiles]):
    def __init__(self):
        super().__init__(IncludedFiles)


class ProblemVersion(PrefetchIDMixin, models.Model):
    problem_version_id = models.AutoField(primary_key=True)
    problem = models.ForeignKey(Problem, models.CASCADE)
    root_group = models.ForeignKey('ProblemTestgroup', models.RESTRICT, related_name='+')
    time_limit_ms = models.IntegerField()
    memory_limit_kb = models.IntegerField()
    output_validator = models.ForeignKey(ProblemOutputValidator, models.RESTRICT, null=True)
    custom_grader = models.ForeignKey(ProblemGrader, models.RESTRICT, null=True)
    included_files = models.JSONField(
        encoder=serialization.DataclassJsonEncoder,
        decoder=IncludedFilesDecoder,
    )
    scoring = models.BooleanField()
    interactive = models.BooleanField()
    score_maximization = models.BooleanField(null=True)

    class Meta:
        db_table = 'problem_version'
