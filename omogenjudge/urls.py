import django.contrib.admin
from django.urls import include, path

import omogenjudge.frontend.urls

urlpatterns = [
    path('admin/', django.contrib.admin.site.urls),
    path('', include(omogenjudge.frontend.urls)),
]
