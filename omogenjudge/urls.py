from django.urls import include, path
import django.contrib.admin

import omogenjudge.frontend.urls

urlpatterns = [
    path('', include(omogenjudge.frontend.urls)),
    path('admin/', django.contrib.admin.site.urls),
]
