import mimetypes

from django.http import Http404, HttpRequest, HttpResponse

from omogenjudge.frontend.decorators import requires_started_contest
from omogenjudge.problems.lookup import find_statement_file
from omogenjudge.storage.models import ProblemStatementFile


@requires_started_contest
def problem_attachment(request: HttpRequest,
                       short_name: str,
                       file_path: str) -> HttpResponse:
    try:
        file = find_statement_file(short_name, file_path)
    except ProblemStatementFile.DoesNotExist:
        raise Http404
    mime_type = mimetypes.guess_type(request.path)
    if mime_type[1]:
        mime_type = f'{mime_type[0]}; charset={mime_type[1]}'
    else:
        mime_type = mime_type[0]
    return HttpResponse(content=file.statement_file.file_contents,
                        content_type=mime_type)
