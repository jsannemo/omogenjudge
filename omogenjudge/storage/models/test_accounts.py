from django.db import IntegrityError
from django.test import TestCase

from omogenjudge.storage.models import Account


class AccountsTests(TestCase):

    def setUp(self) -> None:
        Account.objects.all().delete()

    def testCreateAccount(self):
        Account(
            username='Username',
            full_name='Full Name',
            email='Email'
        ).save()

    def testUniqueUsername(self):
        Account(
            username='Username',
            full_name='Full Name',
            email='Email'
        ).save()

        self.assertRaises(IntegrityError,
                          lambda: Account(username='uSerName', full_name='Full Name', email='other email').save())

    def testUniqueEmail(self):
        Account(
            username='Username',
            full_name='Full Name',
            email='Email'
        ).save()

        self.assertRaises(IntegrityError,
                          lambda: Account(username='other username', full_name='Full Name', email='eMaIl').save())
