from django.db import models
from django.utils import timezone
from django.utils.functional import cached_property

from omogenjudge.util import django_fields


class Contest(models.Model):
    contest_id = models.AutoField(primary_key=True)
    short_name = django_fields.TextField(unique=True)
    host_name = django_fields.TextField(blank=True, null=True)
    start_time = models.DateTimeField()
    duration = models.DurationField()
    title = django_fields.TextField()
    problems = models.ManyToManyField('Problem', through='ContestProblem')

    @cached_property
    def has_started(self):
        return timezone.now() >= self.start_time

    @cached_property
    def has_ended(self):
        return timezone.now() >= self.start_time + self.duration

    def __str__(self):
        repr = self.title
        if self.host_name:
            repr += " @ " + self.host_name
        else:
            repr += " (" + self.short_name + ")"
        return repr

    class Meta:
        db_table = 'contest'


class ContestProblem(models.Model):
    contest = models.ForeignKey('Contest', models.CASCADE)
    problem = models.ForeignKey('Problem', models.CASCADE)
    label = django_fields.TextField(blank=True, null=True)

    class Meta:
        db_table = 'contest_problem'
        unique_together = (('contest', 'problem'),)
