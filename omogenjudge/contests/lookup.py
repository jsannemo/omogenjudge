from typing import Optional

from django.http import HttpRequest

from omogenjudge.storage.models import Contest


def contest_for_request(request: HttpRequest) -> Optional[Contest]:
    try:
        return Contest.objects.get(
            host_name=request.get_host()
        )
    except Contest.DoesNotExist:
        return None
