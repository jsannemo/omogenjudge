from django.urls import include, path

import omogenjudge.frontend.accounts.urls
import omogenjudge.frontend.home.urls
import omogenjudge.frontend.problems.urls
import omogenjudge.frontend.submissions.urls

urlpatterns = [
    path('', include(omogenjudge.frontend.home.urls)),
    path('problems/', include(omogenjudge.frontend.problems.urls)),
    path('accounts/', include(omogenjudge.frontend.accounts.urls)),
    path('submissions/', include(omogenjudge.frontend.submissions.urls))
]
