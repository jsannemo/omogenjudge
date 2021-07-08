import threading
import typing

from django.http import HttpRequest, HttpResponse

from omogenjudge.storage.models import Contest
from omogenjudge.util.django_types import OmogenRequest

_thread_locals = threading.local()


def current_request() -> OmogenRequest:
    return _thread_locals.request


def current_contest() -> typing.Optional[Contest]:
    return current_request().contest


class ThreadLocalMiddleware:
    def __init__(self, get_response: typing.Callable[[HttpRequest], HttpResponse]):
        self.get_response = get_response

    def __call__(self, request: OmogenRequest) -> HttpResponse:
        _thread_locals.request = request
        return self.get_response(request)

    def process_response(self, _request, response):
        if hasattr(_thread_locals, 'request'):
            del _thread_locals.request
        return response

    def process_exception(self, _request, _exception):
        if hasattr(_thread_locals, 'request'):
            del _thread_locals.request
