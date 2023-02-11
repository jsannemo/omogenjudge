from django.urls import path

from omogenjudge.frontend.contests.register import register
from omogenjudge.frontend.countdown.countdown import countdown
from omogenjudge.frontend.problems.list_problems import list_problems
from omogenjudge.frontend.problems.submit import submit
from omogenjudge.frontend.problems.view_problem import view_problem
from omogenjudge.frontend.scoreboard.view_scoreboard import view_scoreboard
from omogenjudge.frontend.submissions.queue import my_submissions, submission_queue
from omogenjudge.frontend.submissions.view_submission import view_submission

urlpatterns = [
    path('<slug:contest_short_name>/register', register, name='contest-register'),
    path('<slug:contest_short_name>/scoreboard', view_scoreboard, name='contest-scoreboard'),
    path('<slug:contest_short_name>/countdown', countdown, name='contest-countdown'),
    path('<slug:contest_short_name>/problems', list_problems, name='contest-problems'),
    path('<slug:contest_short_name>/problems/<slug:short_name>', view_problem, name='contest-problem'),
    path('<slug:contest_short_name>/problems/<slug:short_name>/submit', submit, name='contest-submit'),
    path('<slug:contest_short_name>/problems/<slug:short_name>/<slug:language>', view_problem,
         name='contest-problem-language'),
    path('<slug:contest_short_name>/submissions', submission_queue, name='contest-queue'),
    path('<slug:contest_short_name>/submissions/my', my_submissions, name='contest-my-submissions'),
    path('<slug:contest_short_name>/submissions/<int:sub_id>', view_submission, name='contest-submission'),
]
