from django import template
from django.utils.safestring import mark_safe

from omogenjudge.storage.models import Status, SubmissionRun, Verdict
from omogenjudge.storage.models.problems import ProblemStatement

register = template.Library()


@register.filter
def problem_title(statements: list[ProblemStatement]):
    titles = {
        s.language: s.title
        for s in statements
    }
    # TODO: look at user locale here
    for want_lang in ['en', 'sv']:
        if want_lang in titles:
            return titles[want_lang]
    return next(iter(titles.values()))


@register.filter()
def display_verdict(run: SubmissionRun):
    status = Status(run.status)
    verdict = Verdict(run.verdict)
    if status == Status.DONE:
        if verdict == Verdict.AC:
            return mark_safe('<span class="badge bg-success">Accepted</span>')
        elif verdict == Verdict.TLE:
            return mark_safe('<span class="badge bg-danger">Time Limit Exceeded</span>')
        elif verdict == Verdict.WA:
            return mark_safe('<span class="badge bg-danger">Wrong Answer</span>')
        elif verdict == Verdict.RTE:
            return mark_safe('<span class="badge bg-danger">Run-Time Error</span>')
        raise Exception("Unexpected verdict: " + verdict)
    elif status == Status.QUEUED:
        return mark_safe('<span class="badge bg-dark">Queued</span>')
    elif status == Status.RUNNING:
        return mark_safe('<span class="badge bg-dark">Running</span>')
    elif status == Status.COMPILING:
        return mark_safe('<span class="badge bg-dark">Compiling</span>')
    elif status == Status.COMPILE_ERROR:
        return mark_safe('<span class="badge bg-dark">Compile Error</span>')
    elif status == Status.JUDGE_ERROR:
        return mark_safe('<span class="badge bg-secondary">Judge Error</span>')
