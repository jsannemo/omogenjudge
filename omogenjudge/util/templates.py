from typing import Optional

from django.http import HttpResponse, JsonResponse
from django.shortcuts import render

from omogenjudge.util import serialization
from omogenjudge.util.django_types import OmogenRequest
from omogenjudge.util.serialization import IsDataclassClass, IsDictable


def render_template(request: OmogenRequest, template: str, args: Optional[IsDictable]) -> HttpResponse:
    return render(request, template, context=args.__dict__ if args else {})


def render_json(data: IsDataclassClass) -> HttpResponse:
    return JsonResponse(data, encoder=serialization.DataclassJsonEncoder, safe=False)
