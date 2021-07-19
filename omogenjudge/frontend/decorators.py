from django.http import HttpResponse
from django.shortcuts import redirect


def requires_started_contest(f):
    def wrapped(*args, **kwargs):
        request = args[0]
        if not request.contest.has_started:
            return redirect('home')
        return f(*args, **kwargs)

    return wrapped


def requires_contest(f):
    def wrapped(*args, **kwargs):
        request = args[0]
        if not request.contest:
            return HttpResponse('No active contest')
        return f(*args, **kwargs)

    return wrapped
