from datetime import timedelta

from django.test import TestCase
from django.utils import timezone

from omogenjudge.storage.models import Contest


class ContestsTest(TestCase):

    def testCreateContest(self):
        Contest(
            title='Contest',
            short_name='slug',
            start_time=timezone.now(),
            duration=timedelta(hours=5),
            allow_teams=False,
        ).save()
