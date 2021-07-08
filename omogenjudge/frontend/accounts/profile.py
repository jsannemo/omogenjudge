import dataclasses

import django.contrib.auth
from crispy_forms.helper import FormHelper
from crispy_forms.layout import ButtonHolder, Layout, Submit
from django.contrib.auth.forms import AuthenticationForm
from django.http import HttpRequest, HttpResponse, Http404
from django.shortcuts import redirect

from omogenjudge.accounts.lookup import find_user_by_username
from omogenjudge.storage.models import Account
from omogenjudge.util.django_types import OmogenRequest
from omogenjudge.util.templates import render_template


@dataclasses.dataclass
class ProfileArgs:
    username: str
    full_name: str


def profile(request: OmogenRequest, username: str) -> HttpResponse:
    try:
        user = find_user_by_username(username)
    except Account.DoesNotExist:
        raise Http404
    args = ProfileArgs(
        username=user.username,
        full_name=user.full_name,
    )
    return render_template(request, 'accounts/profile.html', args)
