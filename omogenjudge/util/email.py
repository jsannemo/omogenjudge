import requests
from mailjet_rest import Client

from omogenjudge.settings import MAILJET_API_KEY, MAILJET_API_SECRET
from omogenjudge.storage.models import Account

_mailjet_client = Client(auth=(MAILJET_API_KEY, MAILJET_API_SECRET), version='v3.1')

VERIFY_TEMPLATE = 4182968


def send_email(*, to: Account, subject: str, template: int, variables=None) -> requests.Response:
    if variables is None:
        variables = {}
    variables.update({
        "full_name": to.full_name,
        "username": to.username,
    })
    data = {
        'Messages': [
            {
                "From": {
                    "Email": "domare@kodsport.se",
                    "Name": "Kodsport Sverige"
                },
                "To": [
                    {
                        "Email": to.email,
                        "Name": to.full_name,
                    }
                ],
                "TemplateID": template,
                "TemplateLanguage": True,
                "Subject": subject,
                "Variables": variables
            }
        ]
    }
    return _mailjet_client.send.create(data=data)
