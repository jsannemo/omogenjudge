import dataclasses
from typing import Iterable, Optional

import django.contrib.auth
import django.forms as forms
import requests
from crispy_forms.helper import FormHelper
from crispy_forms.layout import ButtonHolder, Layout, Submit
from django.contrib import messages
from django.contrib.auth.forms import AuthenticationForm, UsernameField
from django.contrib.auth.validators import ASCIIUsernameValidator
from django.core.exceptions import BadRequest, ValidationError
from django.db import transaction
from django.http import HttpResponse
from django.shortcuts import redirect
from django.urls import reverse
from oauthlib.oauth2 import WebApplicationClient

from omogenjudge.accounts.lookup import find_user_by_email, username_exists
from omogenjudge.accounts.register import register_account
from omogenjudge.settings import OAUTH_DETAILS
from omogenjudge.storage.models import Account
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
    social_logins: Iterable[str]


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
        social_logins=OAUTH_DETAILS.keys(),
    )
    return render_template(request, 'accounts/login.html', args)


class SocialCreateForm(forms.Form):
    username = UsernameField(
        widget=forms.TextInput(attrs={'autofocus': True}),
        validators=[ASCIIUsernameValidator()])
    full_name = forms.CharField(label='Full name',
                                min_length=3,
                                max_length=Account._meta.get_field('full_name').max_length)

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.helper = FormHelper()
        self.helper.layout = Layout(
            'username',
            'full_name',
            ButtonHolder(
                Submit('submit', 'Create account')
            )
        )
        self.helper.form_action = reverse('social-create')
        self.username_field = Account._meta.get_field(Account.USERNAME_FIELD)
        username_max_length = self.username_field.max_length
        self.fields['username'].max_length = username_max_length
        self.fields['username'].widget.attrs['maxlength'] = username_max_length

    def clean_username(self):
        username = self.data['username']
        if username_exists(username):
            raise ValidationError('This username is already in use')
        return username


@dataclasses.dataclass
class SocialCreateArgs:
    social_create_form: SocialCreateForm


def social_create(request: OmogenRequest) -> HttpResponse:
    if "social_email" not in request.session:
        raise BadRequest
    email = request.session["social_email"]
    username = request.session.get("social_username", None)
    full_name = request.session.get("social_full_name", None)
    if request.method == "POST":
        form = SocialCreateForm(request.POST)
        with transaction.atomic():
            if form.is_valid():
                form_data = form.cleaned_data
                # RegisterForm checks that the username and email are unique.
                request.session.pop("social_email")
                account = register_account(
                    username=form_data['username'],
                    full_name=form_data['full_name'],
                    email=email,
                    password=None)
                django.contrib.auth.login(request, account)
                return redirect("home")
    else:
        form = SocialCreateForm(
            initial={
                "username": username,
                "full_name": full_name,
            }
        )
    return render_template(request, 'accounts/social_create.html', SocialCreateArgs(form))


def oauth_login(request: OmogenRequest, username_suggestion: Optional[str], email: str,
                full_name: Optional[str]) -> HttpResponse:
    try:
        existing_account = find_user_by_email(email)
        django.contrib.auth.login(request, existing_account)
        return redirect("home")
    except Account.DoesNotExist:
        pass
    request.session["social_email"] = email
    request.session["social_username"] = username_suggestion
    request.session["social_full_name"] = full_name
    return redirect("social-create")


def github_auth(request: OmogenRequest) -> HttpResponse:
    oauth_details = OAUTH_DETAILS["github"]
    client = WebApplicationClient(oauth_details["client_id"])
    if "code" in request.GET:
        code = request.GET["code"]
        request_params = client.prepare_request_body(code=code, client_secret=oauth_details["client_secret"])
        response = requests.post("https://github.com/login/oauth/access_token", data=request_params)
        response.raise_for_status()
        client.parse_request_body_response(response.text)

        uri, headers, body = client.add_token("https://api.github.com/user", "get", )
        user_details = requests.get(uri, data=body, headers=headers).json()
        username = user_details["login"]
        full_name = user_details.get("name")

        uri, headers, body = client.add_token("https://api.github.com/user/emails", "get", )
        emails = requests.get(uri, data=body, headers=headers).json()
        for email in emails:
            if email["primary"] and email["verified"]:
                return oauth_login(request, username, email["email"], full_name)
        messages.error(request, "You need to verify your primary email on your Github account to log in.")
        return redirect("login")

    else:
        url, _, _ = client.prepare_authorization_request(
            "https://github.com/login/oauth/authorize",
            scope=["read:user", "user:email"],
        )
        return redirect(url)


def discord_auth(request: OmogenRequest) -> HttpResponse:
    oauth_details = OAUTH_DETAILS["discord"]
    client = WebApplicationClient(oauth_details["client_id"])
    if "code" in request.GET:
        code = request.GET["code"]
        request_params = client.prepare_request_body(
            code=code,
            redirect_uri=request.build_absolute_uri(reverse("discord-login")),
            client_secret=oauth_details["client_secret"],
        )
        response = requests.post("https://discord.com/api/oauth2/token", data=request_params,
                                 headers={'Content-Type': 'application/x-www-form-urlencoded'})
        response.raise_for_status()
        client.parse_request_body_response(response.text)
        uri, headers, body = client.add_token("https://discord.com/api/users/@me", "get")

        user_details = requests.get(uri, data=body, headers=headers).json()
        if not user_details["verified"]:
            messages.error(request, "You need to verify your email on your Discord account to log in.")
            return redirect("login")
        username = user_details["username"]
        email = user_details["email"]
        return oauth_login(request, username, email, None)
    else:
        url, _, _ = client.prepare_authorization_request(
            "https://discord.com/oauth2/authorize",
            redirect_url=request.build_absolute_uri(reverse("discord-login")),
            scope=["identify", "email"],
        )
        return redirect(url)
