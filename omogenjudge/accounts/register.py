from omogenjudge.accounts.normalization import normalize_username
from omogenjudge.storage.models import Account


def register_account(*, username: str, full_name: str, email: str, password: str) -> Account:
    """Registers a new account.

    No validation of the username or email is performed.
    Username and email is unique by-index.
    """
    account = Account(
        username=normalize_username(username),
        full_name=full_name,
        email=email)
    account.set_password(password)
    account.save()
    return account
