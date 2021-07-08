from django.contrib.auth.views import LogoutView
from django.urls import path, reverse_lazy

from omogenjudge.frontend.accounts.login import login
from omogenjudge.frontend.accounts.profile import profile
from omogenjudge.frontend.accounts.register import register

urlpatterns = [
    path('login/', login, name='login'),
    path('logout/', LogoutView.as_view(next_page=reverse_lazy('home')), name='logout'),
    path('register/', register, name='register'),
    path('<str:username>/', profile, name='profile'),
]
