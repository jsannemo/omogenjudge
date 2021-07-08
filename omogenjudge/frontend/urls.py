from django.urls import include, path

import omogenjudge.frontend.accounts.urls
import omogenjudge.frontend.archive.urls
import omogenjudge.frontend.contests.urls
import omogenjudge.frontend.countdown.urls
import omogenjudge.frontend.home.urls
import omogenjudge.frontend.problems.urls
import omogenjudge.frontend.scoreboard.urls
import omogenjudge.frontend.submissions.urls

urlpatterns = [
    path('accounts/', include(omogenjudge.frontend.accounts.urls)),
    path('archive/', include(omogenjudge.frontend.archive.urls)),
    path('contests/', include(omogenjudge.frontend.contests.urls)),
    path('countdown/', include(omogenjudge.frontend.countdown.urls)),
    path('problems/', include(omogenjudge.frontend.problems.urls)),
    path('submissions/', include(omogenjudge.frontend.submissions.urls)),
    path('scoreboard/', include(omogenjudge.frontend.scoreboard.urls)),
    path('', include(omogenjudge.frontend.home.urls)),
]
