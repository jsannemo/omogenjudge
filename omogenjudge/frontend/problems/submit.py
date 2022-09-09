import django.forms as forms
from crispy_forms.bootstrap import FieldWithButtons
from crispy_forms.helper import FormHelper
from crispy_forms.layout import Layout, Submit
from django.core.files.uploadhandler import FileUploadHandler, StopUpload
from django.http import Http404, HttpResponse, JsonResponse
from django.views.decorators.csrf import csrf_exempt

from omogenjudge.frontend.decorators import requires_user
from omogenjudge.problems.lookup import problem_by_name
from omogenjudge.problems.permissions import can_view_problem
from omogenjudge.storage.models import Problem, Account
from omogenjudge.storage.models.langauges import Language
from omogenjudge.submissions.create import create_submission
from omogenjudge.util.contest_urls import reverse_contest
from omogenjudge.util.django_types import OmogenRequest

SOURCE_CODE_LIMIT = 200000


class SubmitForm(forms.Form):
    upload_files = forms.FileField(
        label="",
        widget=forms.ClearableFileInput(attrs={'multiple': True, 'class': 'form-control'}))
    language = forms.ChoiceField(
        label="",
        choices=Language.as_choices(),
        widget=forms.Select(attrs={'class': 'form-select'}))

    def __init__(self, problem_short_name: str, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.helper = FormHelper()
        self.helper.attrs['id'] = 'submit'
        self.helper.layout = Layout(
            'upload_files',
            FieldWithButtons(
                'language',
                Submit('submit', 'Submit'),
            )
        )
        self.helper.form_action = reverse_contest('submit', short_name=problem_short_name)


class SourceLimitCappingHandler(FileUploadHandler):
    def receive_data_chunk(self, raw_data, start):
        if start + len(raw_data) > self.remaining:
            self.request.META['upload_was_capped'] = True
            raise StopUpload(connection_reset=True)
        return raw_data

    def file_complete(self, file_size):
        self.remaining -= file_size
        return None

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.remaining = SOURCE_CODE_LIMIT


@csrf_exempt
@requires_user
def submit(request: OmogenRequest, short_name: str, user: Account) -> HttpResponse:
    try:
        problem = problem_by_name(short_name)
    except Problem.DoesNotExist:
        raise Http404
    if not can_view_problem(problem):
        raise Http404
    request.upload_handlers.insert(0, SourceLimitCappingHandler(request))  # type: ignore
    exceeded_file_size = request.META.get('upload_was_capped', False)
    form = SubmitForm(problem.short_name, data=request.POST, files=request.FILES)
    if exceeded_file_size:
        return JsonResponse({'errors': {'upload_files': [f'The source code limit is {SOURCE_CODE_LIMIT // 1000} KB.']}})
    # Note: don't validate the rest of the form if we killed uploads
    if not form.is_valid():
        return JsonResponse({'errors': form.errors})
    language = Language(form.cleaned_data['language'])
    submission = create_submission(
        owner=user,
        problem=problem,
        language=language,
        files={f.name: f.read() for f in request.FILES.getlist('upload_files') if f.name is not None}
    )
    return JsonResponse({'submission_link': reverse_contest('submission', sub_id=submission.submission_id)})
