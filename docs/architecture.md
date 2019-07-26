# Architecture
The judge is split up into the following binaries:
- a frontend: Golang web server where users read problems and submit solutions
- a judge master: Golang server that keeps track of the submission queue and judges them
- a local file handling server: Golang server that transfers files to and from judge hosts
- a runner: Golang server that is capable of compiling and running programs through the sandbox server
- a sandbox: C++ server that executes command on the local system in a sandbox

Additionally, a complete judging system needs a database to store persistent application data.
We use PostgreSQL for this.

Within a judging system, several instances of the frontend can run, as well as several instances
of a judge slave.  In between them, a judge master sits to coordinate judging of all submissions.
Since judge coordination is a very light-weight task in comparison to serving web requests or
running submissions, it is unlikely you will ever need more than one judge master instance.

Furthermore, the **evaluation** part of the judging, meaning actually deciding what test cases
to run, what the score of a test case should be, and so on, lies within the judge master. This
means that one can freely use any set of the components independently. For example, one could
replace the compilation and execution of programs in any judging system with the judge slave by
using the same API even if evaluation is to be done slightly differently or support other problem
formats. Furthermore, a judge slave could even be reused between different instances of the judging
system.

## Reliability
All servers are mostly stateless. The exception is the judging master, which keeps
an in-memory queue of submissions that should be judged. However, this list is
also stored in the database, and is reconstructed on startup. Judging is crash-tolerant,
in that the judging master will automatically 
