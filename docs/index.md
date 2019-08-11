# OmogenJudge Documentation
This directory contains all the documentation for developing, running and administrating OmogenJudge.

- [Developer guide](dev/index.md): developing the judge.
- [Production guide](prod/index.md): installing and maintaining an OmogenJudge instance.
- [Administrator guide](admin/index.md): managing OmogenJudge as an administrator, installing problems, creating contests and so on.

To read more on how the judge works internally, you should first read the [architecture overview](architecture.md).
Further, Each component has its own documentation:

- [Sandbox docs](sandbox/index.md): the sandbox component, responsible for running untrusted code.
- [Evaluator docs](evaluator/index.md): the evaluator component, which handles all of the submission evaluation (including compiling and running programs).
- [Judging coordinator docs](master/index.md): the judging coordinator, which delegates submissions for evaluation to evaluators.
- [Frontend docs](frontend/index.md): the frontend component, used both for problem solvers, contest judges, etc.
- [Problem tools docs](problemtools/index.md): the problem tools component, used for verifying the validity of problems and courses.
