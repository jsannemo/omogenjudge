from django.urls import path

from omogenjudge.frontend.scoreboard.view_scoreboard import view_scoreboard

urlpatterns = [
    path('', view_scoreboard, name='scoreboard'),
]
