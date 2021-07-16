from django.urls import path

from omogenjudge.frontend.submissions.queue import submission_queue
from omogenjudge.frontend.submissions.view_submission import view_submission

urlpatterns = [
    path('<int:sub_id>/', view_submission, name='submission'),
    path('', submission_queue, name='queue'),
]
