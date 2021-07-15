from omogenjudge.accounts.normalization import normalize_username
from omogenjudge.storage.models import Account


def register_account(*, username: str, full_name: str, email: str, password: str) -> Account:
    account = Account(
        username=normalize_username(username),
        full_name=full_name,
        email=email)
    account.set_password(password)
    account.save()
    return account
