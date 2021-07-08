import logging
import os.path
import sys

import problemtools.verifyproblem
from django.core.management import BaseCommand
from problemtools.verifyproblem import Problem

import omogenjudge.storage.models
from omogenjudge.problems.install import install_problem
from omogenjudge.util.console import ask_yes_or_no
from omogenjudge.problems.lookup import problem_by_name

logger = logging.getLogger(__name__)


class Command(BaseCommand):
    help = 'Installs a new problem'

    def add_arguments(self, parser):
        parser.add_argument('path', type=str, nargs='+')
        parser.add_argument('--ignore-warnings', action='store_true')

    def handle(self, *args, **options):
        for path in options['path']:
            path = os.path.abspath(path)
            logger.info("Installing problem at path %s", path)
            with Problem(path) as problem:
                try:
                    problem_by_name(problem.shortname)
                    if not ask_yes_or_no("Problem already exists: update it (y/N)?", False):
                        sys.exit(1)
                    update_existing = True
                except omogenjudge.storage.models.Problem.DoesNotExist:
                    update_existing = False

                num_errors, num_warnings = problem.check(problemtools.verifyproblem.default_args())
                if num_errors:
                    logger.error("Problem has errors: exiting")
                    sys.exit(1)
                if num_warnings and not options['ignore_warnings']:
                    if not ask_yes_or_no("Problem has warnings: continue (y/N)? ", False):
                        sys.exit(1)
                install_problem(problem, update_existing=update_existing)
                logger.info("Problem %s installed", problem.shortname)
