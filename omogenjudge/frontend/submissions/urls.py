from django.urls import path

from omogenjudge.frontend.submissions.queue import submission_queue, my_submissions
from omogenjudge.frontend.submissions.view_submission import view_submission

urlpatterns = [
    path('<int:sub_id>/', view_submission, name='submission'),
    path('my', my_submissions, name='my-submissions'),
    path('', submission_queue, name='queue'),
]
