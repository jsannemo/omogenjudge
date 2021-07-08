from django.urls import path

from omogenjudge.frontend.home.home import home

urlpatterns = [
    path('', home, name='home'),
]
