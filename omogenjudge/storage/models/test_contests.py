from datetime import timedelta

from django.test import TestCase
from django.utils import timezone

from omogenjudge.storage.models import Contest, ScoringType


class ContestsTest(TestCase):

    def test_create_contest(self):
        contest = Contest(
            title='Contest',
            short_name='slug',
            start_time=timezone.now(),
            duration=timedelta(hours=5),
            scoring_type=ScoringType.SCORING,
        )
        contest.save()
        contest.refresh_from_db()
        self.assertEqual(ScoringType.SCORING, contest.scoring_type)
