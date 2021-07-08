import dataclasses

import django.contrib.auth
from crispy_forms.helper import FormHelper
from crispy_forms.layout import ButtonHolder, Layout, Submit
from django.contrib.auth.forms import AuthenticationForm
from django.http import HttpRequest, HttpResponse
from django.shortcuts import redirect

from omogenjudge.util.django_types import OmogenRequest
from omogenjudge.util.templates import render_template


class LoginForm(AuthenticationForm):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.helper = FormHelper()
        self.helper.layout = Layout(
            'username',
            'password',
            ButtonHolder(
                Submit('submit', 'Log in')
            )
        )


@dataclasses.dataclass
class LoginArgs:
    login_form: LoginForm


def login(request: OmogenRequest) -> HttpResponse:
    if request.user.is_authenticated:
        return redirect('home')
    if request.method == 'POST':
        form = LoginForm(request=request, data=request.POST)
        if form.is_valid():
            django.contrib.auth.login(request, form.get_user())
            return redirect('/')
    else:
        form = LoginForm()
    args = LoginArgs(
        login_form=form,
    )
    return render_template(request, 'accounts/login.html', args)
