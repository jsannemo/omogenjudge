from typing import Optional, Union

from django.contrib.auth.models import AnonymousUser
from django.http import HttpRequest

from omogenjudge.storage.models.accounts import Account

from omogenjudge.storage.models.contests import Contest, ContestProblem


class OmogenRequest(HttpRequest):
    user: Union[AnonymousUser, Account]
    contest: Optional[Contest]
    contest_site: bool
