from django.urls import path

from omogenjudge.frontend.archive.archive import view_archive

urlpatterns = [
    path('', view_archive, name='archive'),
    path('<path:group_path>', view_archive, name='archive_group'),
]
