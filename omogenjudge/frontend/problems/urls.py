from django.urls import path

from omogenjudge.frontend.problems.view_problem import problem_statement_file, view_problem
from omogenjudge.frontend.problems.list_problems import list_problems
from omogenjudge.frontend.problems.submit import submit

urlpatterns = [
    path('', list_problems, name='problems'),
    path('<slug:short_name>/', view_problem, name='problem'),
    path('<slug:short_name>/submit', submit, name='submit'),
    path('<slug:short_name>/img/<path:file_path>', problem_statement_file),
    path('<slug:short_name>/<slug:language>/', view_problem),
]
