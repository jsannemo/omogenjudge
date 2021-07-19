from django.urls import include, path

import omogenjudge.frontend.accounts.urls
import omogenjudge.frontend.problems.urls
import omogenjudge.frontend.scoreboard.urls
import omogenjudge.frontend.submissions.urls
from omogenjudge.frontend.react.react import react_app

urlpatterns = [
    path('', react_app),
    path('problems/', include(omogenjudge.frontend.problems.urls)),
    path('accounts/', include(omogenjudge.frontend.accounts.urls)),
    path('submissions/', include(omogenjudge.frontend.submissions.urls)),
    path('scoreboard/', include(omogenjudge.frontend.scoreboard.urls))
]
