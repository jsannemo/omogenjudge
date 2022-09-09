from datetime import timedelta
from typing import Optional

from django.core.signing import TimestampSigner, SignatureExpired
from django.http import Http404
from django.urls import reverse

from omogenjudge.accounts.normalization import normalize_username
from omogenjudge.storage.models import Account
from omogenjudge.util.email import send_email, VERIFY_TEMPLATE
from omogenjudge.util.request_global import current_request


def register_account(*, username: str, full_name: str, email: str, password: Optional[str]) -> Account:
    """Registers a new account.

    No validation of the username or email is performed.
    Username and email is unique by-index.
    """
    account = Account(
        username=normalize_username(username),
        full_name=full_name,
        email=email)
    if password:
        account.set_password(password)
    else:
        account.set_unusable_password()
    account.save()
    return account


def send_verification_email(account: Account) -> None:
    response = send_email(to=account, subject="Verifiera ditt konto på judge.kodsport.dev", template=VERIFY_TEMPLATE,
                          variables={
                              "verify_url": create_verification_link(account)
                          })
    response.raise_for_status()


def create_verification_link(account: Account) -> str:
    return current_request().build_absolute_uri(reverse("verify-account", args=[TimestampSigner().sign_object({
        "account_id": account.account_id,
        "email": account.email,
    })]))


def verify_account_from_token(verification_token: str) -> tuple[Account, bool]:
    account_info = TimestampSigner().unsign_object(verification_token, max_age=timedelta(days=7))
    account = Account.objects.get(account_id=account_info["account_id"])
    try:
        TimestampSigner().unsign_object(verification_token, max_age=timedelta(days=7))
    except SignatureExpired:
        return account, False

    if account.email != account_info["email"]:
        raise Http404
    account.email_validated = True
    account.save(update_fields=["email_validated"])
    return account, True
