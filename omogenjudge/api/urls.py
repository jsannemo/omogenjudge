from django.urls import path

from omogenjudge.api.problems.view_problem import view_problem

urlpatterns = [
    path('problems/<str:short_name>/', view_problem),
    path('problems/<str:short_name>/<str:language>/', view_problem),
]
