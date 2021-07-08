from django.contrib import admin

from omogenjudge.storage.models.langauges import Language
from omogenjudge.storage.models.accounts import Account
from omogenjudge.storage.models.stored_file import StoredFile
from omogenjudge.storage.models.problems import ProblemTestgroup, ProblemTestcase, Problem, ProblemVersion, \
    ProblemOutputValidator, ProblemGrader, ProblemStatement, ProblemStatementFile, ScoringMode, VerdictMode, License, \
    IncludedFiles
from omogenjudge.storage.models.contests import ContestProblem, ContestGroupContest, Contest, ContestGroup, \
    ContestStaff, ScoringType
from omogenjudge.storage.models.submissions import SubmissionGroupRun, Submission, SubmissionRun, SubmissionCaseRun, \
    SubmissionFiles, Status, Verdict
from omogenjudge.storage.models.teams import TeamMember, Team


class TeamMemberInline(admin.TabularInline):
    model = TeamMember


class TeamAdmin(admin.ModelAdmin):
    inlines = (TeamMemberInline,)


class TeamInline(admin.TabularInline):
    model = Team


class ContestProblemInline(admin.TabularInline):
    model = ContestProblem


class ContestAdmin(admin.ModelAdmin):
    inlines = (ContestProblemInline, TeamInline)


class ProblemAdmin(admin.ModelAdmin):
    pass


class TestGroupInline(admin.TabularInline):
    model = ProblemTestgroup
    show_change_link = True
    readonly_fields = ['testgroup_name', 'max_score']
    fields = ['testgroup_name', 'max_score']


class TestCaseInline(admin.TabularInline):
    model = ProblemTestcase
    show_change_link = True


class TestGroupAdmin(admin.ModelAdmin):
    inlines = (TestGroupInline, TestCaseInline)


class SubmissionGroupInline(admin.TabularInline):
    model = SubmissionGroupRun
    show_change_link = True


class SubmissionRunAdmin(admin.ModelAdmin):
    inlines = (SubmissionGroupInline,)


class ContestGroupContestInline(admin.TabularInline):
    model = ContestGroupContest
    show_change_link = True


class ContestGroupAdmin(admin.ModelAdmin):
    inlines = (ContestGroupContestInline,)


admin.site.register(Account)
admin.site.register(Contest, ContestAdmin)
admin.site.register(ContestProblem)
admin.site.register(ContestGroup, ContestGroupAdmin)
admin.site.register(ContestGroupContest)
admin.site.register(Team, TeamAdmin)
admin.site.register(TeamMember)
admin.site.register(Submission)
admin.site.register(SubmissionRun, SubmissionRunAdmin)
admin.site.register(SubmissionGroupRun)
admin.site.register(SubmissionCaseRun)
admin.site.register(Problem, ProblemAdmin)
admin.site.register(ProblemTestgroup, TestGroupAdmin)
admin.site.register(ProblemVersion)
admin.site.register(ProblemTestcase)
admin.site.register(ProblemOutputValidator)
admin.site.register(ProblemGrader)
