from django.urls import path

from omogenjudge.frontend.problems.view_problem import problem_attachment

urlpatterns = [
    path('<slug:short_name>/img/<path:file_path>', problem_attachment),
]
