from omogenjudge.accounts.normalization import normalize_username
from omogenjudge.storage.models import Account


def username_exists(username: str) -> bool:
    return Account.objects.filter(
        username__iexact=normalize_username(username)).count() > 0


def email_exists(email: str) -> bool:
    return Account.objects.filter(email__iexact=email).count() > 0
