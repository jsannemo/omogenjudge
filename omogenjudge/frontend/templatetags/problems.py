from django import template

from omogenjudge.storage.models import Status, SubmissionRun
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


@register.filter
def display_verdict(run: SubmissionRun):
    if run.status == Status.DONE.value:
        return run.verdict
    return run.status
