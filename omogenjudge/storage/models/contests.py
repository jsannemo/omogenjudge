import enum

from django.db import models
from django.urls import reverse
from django.utils import timezone
from django.utils.functional import cached_property

from omogenjudge.storage.models import Account, Problem
from omogenjudge.util import django_fields
from omogenjudge.util.django_fields import EnumField, StrEnum


class ScoringType(StrEnum):
    BINARY_WITH_PENALTY = 'binary with penalty'
    SCORING = 'scoring'


class Contest(models.Model):
    contest_id = models.AutoField(primary_key=True)
    short_name = django_fields.TextField(unique=True)
    title = django_fields.TextField()
    host_name = django_fields.TextField(blank=True, null=True)

    # An only virtual contest never runs as a contest, but should still have an e.g. duration because it can be done
    # virtually
    only_virtual_contest = models.BooleanField(default=False)  # TODO: not implemented
    start_time = models.DateTimeField(null=True, blank=True)
    duration = models.DurationField()
    # If set, this contest may be started at an arbitrary time between start_time and flexible_start_window_end_time.
    flexible_start_window_end_time = models.DateTimeField(null=True, blank=True)  # TODO: not implemented

    problems = models.ManyToManyField(Problem, through='ContestProblem')
    scoring_type = EnumField(enum_type=ScoringType)
    # Whether contestants should be able to view anything about each other.
    public_scoreboard = models.BooleanField(default=False)  # TODO: not implemented

    allow_registration = models.BooleanField(default=False)  # TODO: not implemented

    published = models.BooleanField(default=False)

    # These properties are cached in order to provide a consistent view of has started/has ended throughout rendering
    @cached_property
    def has_started(self):
        if not self.start_time:
            return False
        return timezone.now() >= self.start_time

    @cached_property
    def has_ended(self):
        if not self.start_time:
            return False
        return timezone.now() >= self.start_time + self.duration

    @cached_property
    def open_for_practice(self):
        if self.only_virtual_contest:
            return True
        if self.flexible_start_window_end_time and self.flexible_start_window_end_time <= timezone.now():
            return True
        return self.has_ended()

    def is_scoring(self):
        return ScoringType(self.scoring_type) == ScoringType.SCORING

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
    contest_problem_id = models.AutoField(primary_key=True)
    contest = models.ForeignKey(Contest, models.CASCADE)
    problem = models.ForeignKey(Problem, models.CASCADE)
    label = django_fields.TextField(blank=True, null=True)
    binary_pass_score = models.IntegerField(null=True, blank=True, default=None)

    def __str__(self):
        return str(self.contest) + ": " + self.label + " - " + str(self.problem)

    class Meta:
        db_table = 'contest_problem'
        unique_together = (('contest', 'problem'),)


# TODO: not implemented
class ContestStaff(models.Model):
    contest_staff_id = models.AutoField(primary_key=True)
    contest = models.ForeignKey(Contest, models.CASCADE)
    account = models.ForeignKey(Account, models.CASCADE)

    can_administer_contest = models.BooleanField()
    can_answer_clarifications = models.BooleanField()
    can_see_submissions = models.BooleanField()
    can_register_teams = models.BooleanField()

    class Meta:
        db_table = 'contest_staff'
        unique_together = (('contest', 'account'),)


class ContestGroup(models.Model):
    contest_group_id = models.AutoField(primary_key=True)
    name = django_fields.TextField()
    short_name = django_fields.TextField()
    description = django_fields.TextField(null=True, blank=True)
    homepage = django_fields.TextField(null=True, blank=True)
    order = models.IntegerField(null=False, default=0)
    parent = models.ForeignKey('ContestGroup', models.CASCADE, null=True, blank=True, related_name='groups')
    contests = models.ManyToManyField(Contest, through='ContestGroupContest')

    @cached_property
    def subgroups(self):
        return list(self.groups.order_by('order').all())

    @cached_property
    def subcontests(self):
        return list(self.group_contests.order_by('label').filter(contest__published=True).all())

    # TODO: this is inefficient
    def url(self):
        s = [self]
        while s[-1].parent:
            s.append(s[-1].parent)
        return reverse('archive_group', kwargs={'group_path': '/'.join(reversed([x.short_name for x in s]))})

    def __str__(self):
        return ((str(self.parent) + " ") if self.parent else "") + self.name

    class Meta:
        db_table = 'contest_group'
        unique_together = (('contest_group_id', 'short_name'),)


class ContestGroupContest(models.Model):
    contest_group_contest_id = models.AutoField(primary_key=True)
    contest = models.ForeignKey(Contest, models.CASCADE, related_name='group_contests')
    contest_group = models.ForeignKey(ContestGroup, models.CASCADE, related_name='group_contests')
    label = django_fields.TextField()

    def __str__(self):
        return str(self.contest_group) + ": " + self.label

    class Meta:
        db_table = 'contest_group_contest'
        unique_together = (('contest', 'contest_group'),)
