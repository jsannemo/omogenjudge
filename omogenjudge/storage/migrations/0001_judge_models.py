# Generated by Django 3.2.5 on 2021-07-12 08:50

from django.conf import settings
import django.contrib.postgres.fields
from django.db import migrations, models
import django.db.models.deletion
import omogenjudge.storage.models.problems
import omogenjudge.storage.models.submissions
import omogenjudge.util.django_fields
import omogenjudge.util.serialization


class Migration(migrations.Migration):

    initial = True

    dependencies = [
    ]

    operations = [
        migrations.CreateModel(
            name='Account',
            fields=[
                ('password', models.CharField(max_length=128, verbose_name='password')),
                ('account_id', models.AutoField(primary_key=True, serialize=False)),
                ('username', omogenjudge.util.django_fields.TextField(default=None, unique=True)),
                ('full_name', omogenjudge.util.django_fields.TextField(default=None)),
                ('email', omogenjudge.util.django_fields.TextField(default=None, unique=True)),
                ('date_created', models.DateTimeField(auto_now_add=True)),
                ('last_login', models.DateTimeField(null=True)),
            ],
            options={
                'db_table': 'account',
            },
        ),
        migrations.CreateModel(
            name='Contest',
            fields=[
                ('contest_id', models.AutoField(primary_key=True, serialize=False)),
                ('short_name', omogenjudge.util.django_fields.TextField(default=None, unique=True)),
                ('host_name', omogenjudge.util.django_fields.TextField(default=None, null=True)),
                ('start_time', models.DateTimeField()),
                ('selection_window_end_time', models.DateTimeField(null=True)),
                ('duration', models.DurationField()),
                ('title', omogenjudge.util.django_fields.TextField(default=None)),
                ('allow_teams', models.BooleanField()),
            ],
            options={
                'db_table': 'contest',
            },
        ),
        migrations.CreateModel(
            name='Group',
            fields=[
                ('group_id', models.AutoField(primary_key=True, serialize=False)),
                ('group_name', omogenjudge.util.django_fields.TextField(default=None)),
            ],
            options={
                'db_table': 'group',
            },
        ),
        migrations.CreateModel(
            name='Problem',
            fields=[
                ('problem_id', models.AutoField(primary_key=True, serialize=False)),
                ('short_name', omogenjudge.util.django_fields.TextField(default=None, unique=True)),
                ('author', django.contrib.postgres.fields.ArrayField(base_field=omogenjudge.util.django_fields.TextField(default=None), size=None)),
                ('source', omogenjudge.util.django_fields.TextField(default=None)),
                ('license', omogenjudge.util.django_fields.TextField(choices=[('public domain', 'PUBLIC_DOMAIN'), ('cc0', 'CC0'), ('cc by', 'CC_BY'), ('cc by-sa', 'CC_BY_SA'), ('educational', 'EDUCATIONAL'), ('permission', 'PERMISSION')], default=None)),
            ],
            options={
                'db_table': 'problem',
            },
            bases=(omogenjudge.util.django_fields.PrefetchIDMixin, models.Model),
        ),
        migrations.CreateModel(
            name='ProblemOutputValidator',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('validator_run_config', models.JSONField(decoder=omogenjudge.storage.models.problems.ValidatorRunConfigDecoder, encoder=omogenjudge.util.serialization.DataclassJsonEncoder)),
                ('scoring_validator', models.BooleanField()),
            ],
            options={
                'db_table': 'problem_output_validator',
            },
        ),
        migrations.CreateModel(
            name='ProblemTestcase',
            fields=[
                ('problem_testcase_id', models.AutoField(primary_key=True, serialize=False)),
                ('testcase_name', omogenjudge.util.django_fields.TextField(default=None)),
            ],
            options={
                'db_table': 'problem_testcase',
            },
        ),
        migrations.CreateModel(
            name='ProblemTestgroup',
            fields=[
                ('problem_testgroup_id', models.AutoField(primary_key=True, serialize=False)),
                ('testgroup_name', omogenjudge.util.django_fields.TextField(default=None)),
                ('min_score', models.FloatField(null=True)),
                ('max_score', models.FloatField(null=True)),
                ('accept_score', models.FloatField(null=True)),
                ('reject_score', models.FloatField(null=True)),
                ('break_on_reject', models.BooleanField()),
                ('output_validator_flags', django.contrib.postgres.fields.ArrayField(base_field=omogenjudge.util.django_fields.TextField(default=None), size=None)),
                ('parent', models.ForeignKey(null=True, on_delete=django.db.models.deletion.CASCADE, to='storage.problemtestgroup')),
            ],
            options={
                'db_table': 'problem_testgroup',
            },
        ),
        migrations.CreateModel(
            name='ProblemVersion',
            fields=[
                ('problem_version_id', models.AutoField(primary_key=True, serialize=False)),
                ('time_limit_ms', models.IntegerField()),
                ('memory_limit_kb', models.IntegerField()),
                ('included_files', models.JSONField(decoder=omogenjudge.storage.models.problems.IncludedFilesDecoder, encoder=omogenjudge.util.serialization.DataclassJsonEncoder)),
                ('scoring', models.BooleanField()),
                ('interactive', models.BooleanField()),
                ('score_maximization', models.BooleanField(null=True)),
                ('output_validator', models.ForeignKey(null=True, on_delete=django.db.models.deletion.RESTRICT, to='storage.problemoutputvalidator')),
                ('problem', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='storage.problem')),
                ('root_group', models.ForeignKey(on_delete=django.db.models.deletion.RESTRICT, related_name='+', to='storage.problemtestgroup')),
            ],
            options={
                'db_table': 'problem_version',
            },
            bases=(omogenjudge.util.django_fields.PrefetchIDMixin, models.Model),
        ),
        migrations.CreateModel(
            name='StoredFile',
            fields=[
                ('file_hash', models.CharField(max_length=256, primary_key=True, serialize=False)),
                ('file_contents', models.BinaryField()),
            ],
            options={
                'db_table': 'stored_file',
            },
        ),
        migrations.CreateModel(
            name='Submission',
            fields=[
                ('submission_id', models.AutoField(primary_key=True, serialize=False)),
                ('language', omogenjudge.util.django_fields.TextField(default=None)),
                ('date_created', models.DateTimeField()),
                ('submission_files', models.JSONField(decoder=omogenjudge.storage.models.submissions.SubmissionFilesDecoder, encoder=omogenjudge.util.serialization.DataclassJsonEncoder)),
                ('account', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to=settings.AUTH_USER_MODEL)),
            ],
            options={
                'db_table': 'submission',
            },
        ),
        migrations.CreateModel(
            name='Team',
            fields=[
                ('team_id', models.AutoField(primary_key=True, serialize=False)),
                ('team_name', models.TextField(null=True)),
                ('virtual', models.BooleanField()),
                ('unofficial', models.BooleanField()),
                ('start_time', models.DateTimeField(null=True)),
                ('team_data', models.JSONField()),
                ('contest', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='storage.contest')),
            ],
            options={
                'db_table': 'team',
            },
        ),
        migrations.CreateModel(
            name='SubmissionRun',
            fields=[
                ('submission_run_id', models.AutoField(primary_key=True, serialize=False)),
                ('date_created', models.DateTimeField()),
                ('status', models.TextField(choices=[('queued', 'QUEUED'), ('compiling', 'COMPILING'), ('running', 'RUNNING'), ('compile error', 'COMPILE_ERROR'), ('judging error', 'JUDGE_ERROR'), ('done', 'DONE')])),
                ('verdict', models.TextField(choices=[('unjudged', 'UNJUDGED'), ('accepted', 'AC'), ('wrong answer', 'WA'), ('time limit exceeded', 'TLE'), ('run-time error', 'RTE')])),
                ('time_usage_ms', models.IntegerField()),
                ('score', models.IntegerField()),
                ('compile_error', omogenjudge.util.django_fields.TextField(default=None, null=True)),
                ('problem_version', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='storage.problemversion')),
                ('submission', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='storage.submission')),
            ],
            options={
                'db_table': 'submission_run',
            },
        ),
        migrations.CreateModel(
            name='SubmissionGroupRun',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('date_created', models.DateTimeField()),
                ('time_usage_ms', models.IntegerField()),
                ('score', models.IntegerField()),
                ('verdict', models.TextField(choices=[('unjudged', 'UNJUDGED'), ('accepted', 'AC'), ('wrong answer', 'WA'), ('time limit exceeded', 'TLE'), ('run-time error', 'RTE')])),
                ('problem_testgroup', models.ForeignKey(on_delete=django.db.models.deletion.RESTRICT, related_name='+', to='storage.problemtestgroup')),
                ('submission_run', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='storage.submissionrun')),
            ],
            options={
                'db_table': 'submission_group_run',
            },
        ),
        migrations.CreateModel(
            name='SubmissionCaseRun',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('date_created', models.DateTimeField()),
                ('time_usage_ms', models.IntegerField()),
                ('score', models.IntegerField()),
                ('verdict', models.TextField(choices=[('unjudged', 'UNJUDGED'), ('accepted', 'AC'), ('wrong answer', 'WA'), ('time limit exceeded', 'TLE'), ('run-time error', 'RTE')])),
                ('problem_testcase', models.ForeignKey(on_delete=django.db.models.deletion.RESTRICT, related_name='+', to='storage.problemtestcase')),
                ('submission_run', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='storage.submissionrun')),
            ],
            options={
                'db_table': 'submission_case_run',
            },
        ),
        migrations.AddField(
            model_name='submission',
            name='current_run',
            field=models.ForeignKey(db_column='current_run', null=True, on_delete=django.db.models.deletion.SET_NULL, related_name='+', to='storage.submissionrun'),
        ),
        migrations.AddField(
            model_name='submission',
            name='problem',
            field=models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='storage.problem'),
        ),
        migrations.AddField(
            model_name='problemtestgroup',
            name='problem_version',
            field=models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='storage.problemversion'),
        ),
        migrations.AddField(
            model_name='problemtestcase',
            name='input_file_hash',
            field=models.ForeignKey(db_column='input_file_hash', on_delete=django.db.models.deletion.RESTRICT, related_name='+', to='storage.storedfile'),
        ),
        migrations.AddField(
            model_name='problemtestcase',
            name='output_file_hash',
            field=models.ForeignKey(db_column='output_file_hash', on_delete=django.db.models.deletion.RESTRICT, related_name='+', to='storage.storedfile'),
        ),
        migrations.AddField(
            model_name='problemtestcase',
            name='problem_testgroup',
            field=models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='storage.problemtestgroup'),
        ),
        migrations.AddField(
            model_name='problemoutputvalidator',
            name='validator_source_zip',
            field=models.ForeignKey(db_column='file_hash', on_delete=django.db.models.deletion.RESTRICT, related_name='+', to='storage.storedfile'),
        ),
        migrations.AddField(
            model_name='problem',
            name='current_version',
            field=models.ForeignKey(db_column='current_version', on_delete=django.db.models.deletion.RESTRICT, related_name='+', to='storage.problemversion'),
        ),
        migrations.CreateModel(
            name='TeamMember',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('account', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to=settings.AUTH_USER_MODEL)),
                ('team', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='storage.team')),
            ],
            options={
                'db_table': 'team_member',
                'unique_together': {('team', 'account')},
            },
        ),
        migrations.CreateModel(
            name='ProblemStatementFile',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('file_path', omogenjudge.util.django_fields.TextField(default=None)),
                ('attachment', models.BooleanField()),
                ('file_hash', models.ForeignKey(db_column='file_hash', on_delete=django.db.models.deletion.RESTRICT, related_name='+', to='storage.storedfile')),
                ('problem', models.ForeignKey(db_column='problem', on_delete=django.db.models.deletion.CASCADE, to='storage.problem')),
            ],
            options={
                'db_table': 'problem_statement_file',
                'unique_together': {('problem', 'file_path')},
            },
        ),
        migrations.CreateModel(
            name='ProblemStatement',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('language', omogenjudge.util.django_fields.TextField(default=None)),
                ('title', omogenjudge.util.django_fields.TextField(default=None)),
                ('html', omogenjudge.util.django_fields.TextField(default=None)),
                ('problem', models.ForeignKey(db_column='problem', on_delete=django.db.models.deletion.CASCADE, to='storage.problem')),
            ],
            options={
                'db_table': 'problem_statement',
                'unique_together': {('problem', 'language')},
            },
        ),
        migrations.CreateModel(
            name='GroupMember',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('admin', models.BooleanField()),
                ('account', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to=settings.AUTH_USER_MODEL)),
                ('group', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='storage.group')),
            ],
            options={
                'db_table': 'account_group_member',
                'unique_together': {('group', 'account')},
            },
        ),
        migrations.CreateModel(
            name='GroupContest',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('contest', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='storage.contest')),
                ('group', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='storage.group')),
            ],
            options={
                'db_table': 'group_contest',
                'unique_together': {('group', 'contest')},
            },
        ),
        migrations.CreateModel(
            name='ContestProblem',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('label', omogenjudge.util.django_fields.TextField(default=None)),
                ('contest', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='storage.contest')),
                ('problem', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='storage.problem')),
            ],
            options={
                'db_table': 'contest_problem',
                'unique_together': {('contest', 'problem')},
            },
        ),
    ]
