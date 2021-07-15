import django.forms as forms
from crispy_forms.bootstrap import FieldWithButtons
from crispy_forms.helper import FormHelper
from crispy_forms.layout import Layout, Submit
from django.contrib.auth.decorators import login_required
from django.core.files.uploadhandler import FileUploadHandler, SkipFile, StopUpload
from django.http import HttpRequest, HttpResponse, JsonResponse
from django.urls import reverse
from django.views.decorators.csrf import csrf_exempt, csrf_protect

from omogenjudge.problems.lookup import lookup_problem
from omogenjudge.storage.models import Problem
from omogenjudge.storage.models.langauges import Language
from omogenjudge.submissions.create import create_submission

SOURCE_CODE_LIMIT = 200000


class SubmitForm(forms.Form):
    files = forms.FileField(
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
            'files',
            FieldWithButtons(
                'language',
                Submit('submit', 'Submit'),
            )
        )
        self.helper.form_action = reverse('submit', kwargs={'short_name': problem_short_name})


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
def submit(request: HttpRequest, short_name: str) -> HttpResponse:
    problem = lookup_problem(short_name)
    request.upload_handlers.insert(0, SourceLimitCappingHandler(request))
    return _submit(request, problem)


@login_required
@csrf_protect
def _submit(request: HttpRequest, problem: Problem):
    exceeded_file_size = request.META.get('upload_was_capped', False)
    form = SubmitForm(problem.short_name, data=request.POST, files=request.FILES)
    if exceeded_file_size:
        return JsonResponse({'errors': {'files': [f'The source code limit is {SOURCE_CODE_LIMIT // 1000} KB.']}})
    # Note: don't validate the rest of the form if we killed uploads
    if not form.is_valid():
        return JsonResponse({'errors': form.errors})
    language = Language(form.cleaned_data['language'])
    submission = create_submission(
        owner=request.user,
        problem=problem,
        language=language,
        files={f.name: f.read() for f in request.FILES.getlist('files')}
    )
    return JsonResponse({'submission_id': submission.submission_id})
