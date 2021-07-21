import React from 'react';

import {ViewProblemRequest, ViewProblemResponse} from "webapi/proto/problem_service_pb";
import {ProblemService} from "webapi/proto/problem_service_pb_service";
import {ReactRequest, useRequest} from "../../api/react_request";
import {useParams} from 'react-router-dom';
import {ProblemStatement} from "webapi/proto/problems_pb";

function ViewProblem() {
  const {shortname} = useParams();
  const request = new ViewProblemRequest();
  request.setShortName(shortname);
  const data: ReactRequest<ViewProblemResponse> = useRequest(ProblemService.ViewProblem, request);
  return (
      <div>
        {data.loading ? "Loading problem..." :
            <Problem statement={data.message.getStatement()}/>}
      </div>
  );
}

type ProblemProps = {
  statement: ProblemStatement,
}

function Problem({statement}: ProblemProps) {
  return <div>
    <h1>{statement.getTitle()}</h1>
    <div dangerouslySetInnerHTML={{__html: statement.getHtml()}}/>
  </div>
}

export default ViewProblem;
