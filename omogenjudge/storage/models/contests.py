from django.db import models

from omogenjudge.util import django_fields


class Contest(models.Model):
    contest_id = models.AutoField(primary_key=True)
    short_name = django_fields.TextField(unique=True)
    host_name = django_fields.TextField(null=True)
    start_time = models.DateTimeField()
    selection_window_end_time = models.DateTimeField(null=True)
    duration = models.DurationField()
    title = django_fields.TextField()
    allow_teams = models.BooleanField()

    class Meta:
        db_table = 'contest'


class ContestProblem(models.Model):
    contest = models.ForeignKey('Contest', models.CASCADE)
    problem = models.ForeignKey('Problem', models.CASCADE)
    label = django_fields.TextField()

    class Meta:
        db_table = 'contest_problem'
        unique_together = (('contest', 'problem'),)
