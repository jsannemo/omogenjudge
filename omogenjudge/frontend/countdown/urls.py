from django.urls import path

from omogenjudge.frontend.countdown.countdown import countdown

urlpatterns = [
    path('', countdown, name='countdown'),
]
