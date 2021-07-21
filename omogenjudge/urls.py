from django.urls import include, path
import django.contrib.admin

import omogenjudge.frontend.urls
import omogenjudge.api.urls

urlpatterns = [
    path('api/', include(omogenjudge.api.urls)),
    path('admin/', django.contrib.admin.site.urls),
    path('', include(omogenjudge.frontend.urls)),
]
