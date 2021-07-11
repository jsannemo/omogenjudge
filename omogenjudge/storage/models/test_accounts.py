from django.test import TestCase

from omogenjudge.storage.models import Account


class AccountsTests(TestCase):

    def testCreateAccount(self):
        Account(
            username='Username',
            full_name='Full Name',
            email='Email'
        ).save()
