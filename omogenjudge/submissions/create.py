import binascii

from django.db import transaction

from omogenjudge.storage.models import Account, Language, Problem, Status, Submission, SubmissionFiles, SubmissionRun, \
    Verdict


def create_submission(*, owner: Account, problem: Problem, language: Language, files: dict[str, bytes]) -> Submission:
    with transaction.atomic():
        submission = Submission(
            account=owner,
            problem=problem,
            language=language,
            submission_files=SubmissionFiles(
                {name: binascii.b2a_base64(value).decode("ascii") for name, value in files.items()}),
        )
        submission.prefetch_id()
        run = SubmissionRun(
            submission_id=submission.submission_id,
            status=Status.QUEUED,
            problem_version_id=problem.current_version_id,
            verdict=Verdict.UNJUDGED,
        )
        run.save()
        submission.current_run = run
        submission.save()
    return submission
