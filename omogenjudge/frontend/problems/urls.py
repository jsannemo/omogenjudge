from django.urls import path

from omogenjudge.frontend.problems.list_problems import list_problems
from omogenjudge.frontend.problems.view_problem import problem_statement_file

urlpatterns = [
    path('', list_problems, name='problems'),
    path('<slug:short_name>/img/<path:file_path>', problem_statement_file),
    path('<slug:short_name>/attachment/<path:file_path>', problem_statement_file, name='problem_attachment'),
]

