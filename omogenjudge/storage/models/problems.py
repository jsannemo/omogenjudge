import dataclasses
import enum
import typing

from django.db import models
from django.contrib.postgres.fields import ArrayField

from omogenjudge.util import serialization, django_fields
from omogenjudge.util.django_fields import PrefetchIDMixin
from omogenjudge.util.enums import EnumChoices


class License(EnumChoices[str], enum.Enum):
    PUBLIC_DOMAIN = 'public domain'
    CC0 = 'cc0'
    CC_BY = 'cc by'
    CC_BY_SA = 'cc by-sa'
    EDUCATIONAL = 'educational'
    PERMISSION = 'permission'


class Problem(PrefetchIDMixin, models.Model):
    problem_id = models.AutoField(primary_key=True)
    short_name = django_fields.TextField(unique=True)
    author = ArrayField(django_fields.TextField())
    source = django_fields.TextField()
    license = django_fields.TextField(choices=License.as_choices())
    current_version = models.ForeignKey(
        'ProblemVersion', models.RESTRICT, db_column='current_version',
        related_name='+'
    )

    class Meta:
        db_table = 'problem'


@dataclasses.dataclass(frozen=True)
class ValidatorRunConfig:
    run_command: list[str]


class ValidatorRunConfigDecoder(serialization.DataclassJsonDecoder[ValidatorRunConfig]):
    def __init__(self):
        super().__init__(ValidatorRunConfig)


class ProblemOutputValidator(models.Model):
    validator_run_config = models.JSONField(
        encoder=serialization.DataclassJsonEncoder,
        decoder=ValidatorRunConfigDecoder,
    )
    validator_source_zip = models.ForeignKey('StoredFile', models.RESTRICT, db_column='file_hash', related_name='+')
    scoring_validator = models.BooleanField()

    class Meta:
        db_table = 'problem_output_validator'


class ProblemStatement(models.Model):
    problem = models.ForeignKey(Problem, models.CASCADE, db_column='problem')
    language = django_fields.TextField()
    title = django_fields.TextField()
    html = django_fields.TextField()

    class Meta:
        db_table = 'problem_statement'
        unique_together = (('problem', 'language'),)


class ProblemStatementFile(models.Model):
    problem = models.ForeignKey(Problem, models.CASCADE, db_column='problem')
    file_path = django_fields.TextField()
    file_hash = models.ForeignKey('StoredFile', models.RESTRICT, db_column='file_hash', related_name='+')
    attachment = models.BooleanField()

    class Meta:
        db_table = 'problem_statement_file'
        unique_together = (('problem', 'file_path'),)


class ProblemTestcase(models.Model):
    problem_testcase_id = models.AutoField(primary_key=True)
    problem_testgroup = models.ForeignKey('ProblemTestgroup', models.CASCADE)
    testcase_name = django_fields.TextField()
    input_file_hash = models.ForeignKey('StoredFile', models.RESTRICT, db_column='input_file_hash', related_name='+')
    output_file_hash = models.ForeignKey('StoredFile', models.RESTRICT, db_column='output_file_hash', related_name='+')

    class Meta:
        db_table = 'problem_testcase'


class ProblemTestgroup(models.Model):
    problem_testgroup_id = models.AutoField(primary_key=True)
    parent = models.ForeignKey('self', models.CASCADE, null=True)
    problem_version = models.ForeignKey('ProblemVersion', models.CASCADE)
    testgroup_name = django_fields.TextField()
    min_score = models.FloatField(null=True)
    max_score = models.FloatField(null=True)
    accept_score = models.FloatField(null=True)
    reject_score = models.FloatField(null=True)
    break_on_reject = models.BooleanField()
    output_validator_flags = ArrayField(django_fields.TextField())

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
    included_files = models.JSONField(
        encoder=serialization.DataclassJsonEncoder,
        decoder=IncludedFilesDecoder,
    )
    scoring = models.BooleanField()
    interactive = models.BooleanField()
    score_maximization = models.BooleanField(null=True)

    class Meta:
        db_table = 'problem_version'
