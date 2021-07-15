from django.contrib.auth.views import LogoutView
from django.urls import path, reverse_lazy

from omogenjudge.frontend.accounts.auth import login, register

urlpatterns = [
    path('login/', login, name='login'),
    path('logout/', LogoutView.as_view(next_page=reverse_lazy('home')), name='logout'),
    path('register/', register, name='register'),
    path('<str:username>/', login, name='profile'),
]
