from django.contrib.auth.views import LogoutView
from django.urls import path, reverse_lazy

from omogenjudge.frontend.accounts.login import login, github_auth, social_create, discord_auth
from omogenjudge.frontend.accounts.profile import profile
from omogenjudge.frontend.accounts.register import register, verify_account

urlpatterns = [
    path('login/', login, name='login'),
    path('logout/', LogoutView.as_view(next_page=reverse_lazy('home')), name='logout'),
    path('register/', register, name='register'),
    path('verify-account/<str:verify_token>', verify_account, name='verify-account'),
    path('o/github/', github_auth, name='github-login'),
    path('o/discord/', discord_auth, name='discord-login'),
    path('o/create/', social_create, name='social-create'),
    path('<str:username>/', profile, name='profile'),
]
