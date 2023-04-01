import django.contrib.admin
from django.conf import settings
from django.urls import include, path

import omogenjudge.frontend.urls

urlpatterns = [
    path('admin/', django.contrib.admin.site.urls),
    path('', include(omogenjudge.frontend.urls)),
]
if settings.DEBUG:
    urlpatterns += [path('__debug__/', include('debug_toolbar.urls'))]
