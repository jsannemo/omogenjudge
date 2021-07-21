from django.urls import include, path

import omogenjudge.frontend.problems.urls
import omogenjudge.frontend.submissions.urls
from omogenjudge.frontend.react.react import react_app

urlpatterns = [
    path('problems/', include(omogenjudge.frontend.problems.urls)),
    path('', react_app),
    path('<path:path>', react_app),
]
