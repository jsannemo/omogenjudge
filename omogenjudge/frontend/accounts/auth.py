import dataclasses

from crispy_forms.helper import FormHelper
from crispy_forms.layout import ButtonHolder, Layout, Submit
from django import forms
import django.contrib.auth
from django.contrib.auth.forms import AuthenticationForm, UsernameField
from django.contrib.auth.validators import ASCIIUsernameValidator
from django.core.exceptions import ValidationError
from django.db import transaction
from django.forms import PasswordInput
from django.http import HttpRequest, HttpResponse
from django.shortcuts import redirect, render

from omogenjudge.accounts.lookup import email_exists, username_exists
from omogenjudge.accounts.register import register_account
from omogenjudge.storage.models import Account
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


def login(request: HttpRequest) -> HttpResponse:
    if request.method == 'POST':
        form = LoginForm(request=request, data=request.POST)
        if form.is_valid():
            django.contrib.auth.login(request, form.get_user())
            return redirect('/')
        print("errors", form.errors, form.is_valid())
    else:
        form = LoginForm()
    args = LoginArgs(
        login_form=form,
    )
    return render_template(request, 'accounts/login.html', args)


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


def register(request: HttpRequest) -> HttpResponse:
    with transaction.atomic():
        if request.method == 'POST':
            form = RegisterForm(request.POST)
            if form.is_valid():
                form_data = form.cleaned_data
                account = register_account(
                    username=form_data['username'],
                    full_name=form_data['full_name'],
                    email=form_data['email'],
                    password=form_data['password'])
                django.contrib.auth.login(request, account)
                return redirect('/')
        else:
            form = RegisterForm()
        args = RegisterArgs(
            register_form=form,
        )
    return render_template(request, 'accounts/register.html', args)
