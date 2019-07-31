// Submission storage models.
package submissions

import (
  runpb "github.com/jsannemo/omogenjudge/runner/api"
)

type Status string

const (
  StatusCompiling Status = "compiling"
  StatusRunning Status = "running"
  StatusSuccessful Status = "successful"
)

type Verdict string

// A stored file
type Submission struct {
  // A numeric ID for the submission.
  // This may exposed externally.
  SubmissionId int32

  // The account of the author of the submission.
  AccountId int32

  // The ID of the problem the submission was for.
  ProblemId int32

  // An identifier of the location of this resource.
  Files []*SubmissionFile

  Status Status

  Verdict Verdict
}

type SubmissionFile struct {
  // The file path of the submission file.
  Path string

  // The contents of the submitted file.
  Contents []byte
}

func (s *Submission) ToRunnerProgram() *runpb.Program {
  var files []*runpb.SourceFile
  for _, file := range s.Files {
    files = append(files, &runpb.SourceFile{
      Path: file.Path,
      Contents: file.Contents,
    })
  }
  return &runpb.Program{
    Sources: files,
    LanguageId: "gpp17",
  }
}
