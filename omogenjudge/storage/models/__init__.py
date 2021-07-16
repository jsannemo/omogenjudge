from .stored_file import *
from .accounts import *
from .problems import *
from .contests import *
from .groups import *
from .submissions import *
from .teams import *
from .langauges import *


class TeamMemberInline(admin.TabularInline):
    model = TeamMember


class TeamAdmin(admin.ModelAdmin):
    inlines = (TeamMemberInline,)


class TeamInline(admin.TabularInline):
    model = Team


class ProblemInline(admin.TabularInline):
    model = ContestProblem


class ContestAdmin(admin.ModelAdmin):
    inlines = (ProblemInline, TeamInline)


admin.site.register(Account)
admin.site.register(Contest, ContestAdmin)
admin.site.register(ContestProblem)
admin.site.register(Team, TeamAdmin)
admin.site.register(TeamMember)
