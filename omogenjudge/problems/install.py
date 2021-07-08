import io
import os.path
import shlex
import tempfile
import zipfile
from typing import Optional

import problemtools.run
from django.db import transaction
from problemtools import problem2html
from problemtools.verifyproblem import Problem as ToolsProblem, TestCase as ToolsCase, TestCaseGroup as ToolsGroup

from omogenjudge.storage.models import IncludedFiles, Problem, ProblemOutputValidator, ProblemStatement, \
    ProblemStatementFile, ProblemTestcase, \
    ProblemTestgroup, ProblemVersion, StoredFile, ScoringMode, VerdictMode, ProblemGrader, License
from omogenjudge.storage.stored_files import insert_file


def _add_problem(problem: ToolsProblem) -> Problem:
    author = [author.strip() for author in problem.config.get('author').split(",")]
    return Problem(
        short_name=problem.shortname,
        author=author,
        license=License(problem.config.get('license')),
        source=problem.config.get('source'),
    )


def _add_case(db_group: ProblemTestgroup, case: ToolsCase) -> ProblemTestcase:
    name = os.path.basename(case._base)
    with open(case.infile, 'rb') as infile:
        input_file = insert_file(infile.read())
    with open(case.ansfile, 'rb') as outfile:
        output_file = insert_file(outfile.read())
    db_case = ProblemTestcase(
        problem_testgroup=db_group,
        testcase_name=name,
        input_file=input_file,
        output_file=output_file,
    )
    db_case.save()
    return db_case


def _add_group(parent: Optional[ProblemTestgroup], group: ToolsGroup, db_version: ProblemVersion) -> ProblemTestgroup:
    group_name = os.path.basename(group._datadir)
    if parent:
        group_name = f'{parent.testgroup_name}/{group_name}'
    output = shlex.split(group._problem.config.get('validator_flags') + ' ' + group.config['output_validator_flags'])

    scoring_mode = ScoringMode.SUM
    verdict_mode = VerdictMode.WORST_ERROR
    ignore_sample = False
    accept_if_any_accepted = False
    grader_flags = shlex.split(group.config['grader_flags'])
    custom_grading = group.config.get("grading", "default") == "custom"
    if not custom_grading:
        for flag in grader_flags:
            try:
                scoring_mode = ScoringMode(flag)
            except ValueError:
                pass
            try:
                verdict_mode = VerdictMode(flag)
            except ValueError:
                pass
            if flag == 'ignore_sample':
                ignore_sample = True
            if flag == 'accept_if_any_accepted':
                accept_if_any_accepted = True
    db_group = ProblemTestgroup(
        parent=parent,
        problem_version=db_version,
        testgroup_name=group_name,
        break_on_reject=group.config['on_reject'] == 'break',
        scoring_mode=scoring_mode,
        verdict_mode=verdict_mode,
        accept_if_any_accepted=accept_if_any_accepted,
        ignore_sample=ignore_sample,
        output_validator_flags=output,
        grader_flags=grader_flags,
        custom_grading=custom_grading
    )
    if db_version.scoring:
        db_group.min_score, db_group.max_score = group.get_score_range()
        db_group.accept_score = group.config['accept_score']
        db_group.reject_score = group.config['reject_score']
    db_group.save()
    for case in group.get_testcases():
        _add_case(db_group, case)
    for subgroup in group.get_subgroups():
        _add_group(db_group, subgroup, db_version)
    return db_group


def _add_testdata(problem: ToolsProblem, db_version: ProblemVersion) -> ProblemTestgroup:
    return _add_group(None, problem.testdata, db_version)


def _included_files(problem: ToolsProblem) -> IncludedFiles:
    include_dict: dict[str, dict[str, str]] = {}
    includes = os.path.join(problem.probdir, 'include')
    if os.path.isdir(includes):
        for langname in os.listdir(includes):
            include_dict[langname] = {}
            for lang_dir, _, files in os.walk(os.path.join(includes, langname)):
                for file_name in files:
                    abs_path = os.path.join(lang_dir, file_name)
                    rel_path = os.path.relpath(abs_path, lang_dir)
                    with open(abs_path, 'r') as file:
                        include_dict[langname][rel_path] = file.read()
    return IncludedFiles(files_by_language=include_dict)


def _zip_program(path) -> StoredFile:
    zip_buf = io.BytesIO()
    zip_handler = zipfile.ZipFile(zip_buf, 'w', zipfile.ZIP_DEFLATED)
    for root, dirs, files in os.walk(path):
        for file in files:
            zip_handler.write(os.path.join(root, file), os.path.join(os.path.relpath(root, path), file))
    zip_handler.close()
    return insert_file(zip_buf.getbuffer())


def _add_validator(problem: ToolsProblem) -> ProblemOutputValidator:
    # We recompile the validator to ensure that we have a directory only with a single validator present.
    # Otherwise, it's annoying to handle the case of multiple single-file validators in the same directory.
    with tempfile.TemporaryDirectory() as tmp_validator:
        validator = problemtools.run.find_programs(
            os.path.join(problem.probdir, "output_validators"),
            language_config=problem.language_config,
            work_dir=tmp_validator)[0]
        validator.compile()
        db_validator = ProblemOutputValidator(
            run_command=validator.get_runcmd(tmp_validator),
            validator_zip=_zip_program(tmp_validator),
            scoring_validator=problem.config.get('grading')['custom_scoring']
        )
    db_validator.save()
    return db_validator


def _add_grader(problem: ToolsProblem) -> Optional[ProblemGrader]:
    # We recompile the grader to ensure that we have a directory only with a single grader present.
    # Otherwise, it's annoying to handle the case of multiple single-file validators in the same directory.
    with tempfile.TemporaryDirectory() as tmp_grader:
        graders = problemtools.run.find_programs(
            os.path.join(problem.probdir, "graders"),
            language_config=problem.language_config,
            work_dir=tmp_grader)
        if not graders:
            return None
        grader = graders[0]
        grader.compile()
        db_grader = ProblemGrader(
            run_command=grader.get_runcmd(tmp_grader),
            grader_zip=_zip_program(tmp_grader),
        )
    db_grader.save()
    return db_grader


def _add_version(problem: ToolsProblem, db_problem: Problem) -> ProblemVersion:
    limits = problem.config.get('limits')
    db_version = ProblemVersion(
        problem=db_problem,
        time_limit_ms=limits.get('time') * 1000,
        memory_limit_kb=limits.get('memory') * 1000,
        scoring=problem.is_scoring,
        interactive=problem.is_interactive,
        included_files=_included_files(problem),
    )
    db_version.prefetch_id()
    db_version.root_group = _add_testdata(problem, db_version)
    if db_version.scoring:
        grading_settings = problem.config.get('grading')
        db_version.score_maximization = grading_settings['objective'] == 'max'
    if problem.config.get('validation-type') == 'custom':
        db_version.output_validator = _add_validator(problem)
    db_version.custom_grader = _add_grader(problem)
    db_version.save()
    db_problem.current_version = db_version
    return db_version


def _add_statement(problem: ToolsProblem, language_code: str, db_problem: Problem):
    statement = ProblemStatement(
        language=language_code,
        problem=db_problem,
        title=problem.config.get('name')[language_code]
    )
    htmlopt = problem2html.ConvertOptions()
    with tempfile.TemporaryDirectory() as tmp_dest:
        htmlopt.destdir = tmp_dest
        htmlopt.quiet = True
        htmlopt.language = language_code
        htmlopt.bodyonly = True
        htmlopt.css = False
        htmlopt.headers = False
        htmlopt.imgbasedir = f"/problems/{problem.shortname}/img/{language_code}"
        problem2html.convert(problem.probdir, htmlopt)
        with open(os.path.join(tmp_dest, 'index.html'), 'r') as html:
            statement.html = html.read()
        statement.save()
        for root, _, files in os.walk(tmp_dest):
            for file in files:
                file_path = os.path.join(root, file)
                rel_path = os.path.relpath(file_path, tmp_dest)
                with open(file_path, 'rb') as statement_file:
                    ProblemStatementFile(
                        problem=db_problem,
                        file_path=f'{language_code}/{rel_path}',
                        statement_file=insert_file(statement_file.read()),
                        attachment=False,
                    ).save()
    statement.save()


def _add_statements(problem: ToolsProblem, db_problem: Problem):
    db_problem.statements.all().delete()
    db_problem.problemstatementfile_set.all().delete()
    statement = problem.statement
    for lang in statement.languages:
        _add_statement(problem, lang, db_problem)


def install_problem(problem: ToolsProblem, *, update_existing=False) -> Problem:
    with transaction.atomic():
        try:
            db_problem = Problem.objects.get(short_name=problem.shortname)
            if not update_existing:
                raise ValueError(
                    f"Problem {problem.shortname} already exists, but did not expect to update an existing problem")
        except Problem.DoesNotExist:
            db_problem = _add_problem(problem)
            db_problem.prefetch_id()

        _add_version(problem, db_problem)
        _add_statements(problem, db_problem)
        db_problem.save()
    return db_problem
