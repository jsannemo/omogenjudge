import dataclasses
from typing import Iterable

import django.contrib.auth
from crispy_forms.helper import FormHelper
from crispy_forms.layout import ButtonHolder, Layout, Submit
from django import forms
from django.contrib import messages
from django.contrib.auth.forms import UsernameField
from django.contrib.auth.validators import ASCIIUsernameValidator
from django.core.exceptions import ValidationError
from django.db import transaction
from django.forms import PasswordInput
from django.http import HttpResponse
from django.shortcuts import redirect

from omogenjudge.accounts.lookup import email_exists, username_exists
from omogenjudge.accounts.register import register_account, send_verification_email, verify_account_from_token
from omogenjudge.settings import OAUTH_DETAILS, REQUIRE_EMAIL_AUTH
from omogenjudge.storage.models import Account
from omogenjudge.util.django_types import OmogenRequest
from omogenjudge.util.templates import render_template


class RegisterForm(forms.Form):
    username = UsernameField(
        widget=forms.TextInput(attrs={'autofocus': True}),
        validators=[ASCIIUsernameValidator()])

    full_name = forms.CharField(label='Full name',
                                min_length=3,
                                max_length=Account._meta.get_field('full_name').max_length)
    email = forms.EmailField(label='Email', max_length=256)
    password = forms.CharField(label='Password', strip=False, min_length=5, max_length=100, widget=PasswordInput())
    password_confirmation = forms.CharField(label='Password confirmation', strip=False, widget=PasswordInput())

    def clean_username(self):
        username = self.data['username']
        if username_exists(username):
            raise ValidationError('This username is already in use')
        return username

    def clean_email(self):
        email = self.data['email']
        if email_exists(email):
            raise ValidationError('This email is already in use')
        return email

    def clean(self):
        cleaned_data = super().clean()
        password = cleaned_data.get('password')
        confirmation = cleaned_data.get('password_confirmation')
        if password != confirmation:
            raise ValidationError({'password_confirmation': 'The passwords does not match'})

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.helper = FormHelper()
        self.helper.layout = Layout(
            'username',
            'full_name',
            'email',
            'password',
            'password_confirmation',
            ButtonHolder(
                Submit('submit', 'Register')
            )
        )
        self.username_field = Account._meta.get_field(Account.USERNAME_FIELD)
        username_max_length = self.username_field.max_length
        self.fields['username'].max_length = username_max_length
        self.fields['username'].widget.attrs['maxlength'] = username_max_length


@dataclasses.dataclass
class RegisterArgs:
    register_form: RegisterForm
    social_logins: Iterable[str]


def register(request: OmogenRequest) -> HttpResponse:
    if request.method == 'POST':
        form = RegisterForm(request.POST)
        account = None
        with transaction.atomic():
            if form.is_valid():
                form_data = form.cleaned_data
                # RegisterForm checks that the username and email are unique.
                account = register_account(
                    username=form_data['username'],
                    full_name=form_data['full_name'],
                    email=form_data['email'],
                    password=form_data['password'])
        if account:
            if REQUIRE_EMAIL_AUTH:
                send_verification_email(account)
            else:
                account.email_validated = True
                account.save()
            messages.add_message(request, messages.INFO,
                                 'Your account was successfully created. '
                                 'To log in you must verify your account. '
                                 'We have sent an email to you with instructions on how to verify it.')
            return redirect('login')
    else:
        form = RegisterForm()
    args = RegisterArgs(
        register_form=form,
        social_logins=OAUTH_DETAILS.keys(),
    )
    return render_template(request, 'accounts/register.html', args)


@dataclasses.dataclass
class OldAccountArgs:
    pass


def verify_account(request: OmogenRequest, verify_token: str) -> HttpResponse:
    account, ok = verify_account_from_token(verify_token)
    if not ok:
        messages.error(request, "The verification link you clicked had expired since it's over 7 days old. "
                                "We have sent a new verification link to the same email.")
        return redirect('login')

    django.contrib.auth.login(request, account)
    messages.success(request, "Your account has been successfully verified!")
    return redirect('/')
