import React from 'react';

import {ViewProblemRequest, ViewProblemResponse} from "webapi/proto/problem_service_pb";
import ProblemPb from "webapi/proto/problems_pb";
import {ProblemService} from "webapi/proto/problem_service_pb_service";
import {ReactRequest, useRequest} from "../../api/react_request";
import {useParams} from 'react-router-dom';
import {Card} from "react-bootstrap";

import './ViewProblem.scss';
import {LatexContainer} from "../../shared/latex/LatexContainer";

function ViewProblem() {
  const {shortname} = useParams();
  const request = new ViewProblemRequest();
  request.setShortName(shortname);
  const data: ReactRequest<ViewProblemResponse> = useRequest(ProblemService.ViewProblem, request);
  return (
      <div className={"row"}>
        {data.loading ? "Loading problem..." :
            <ProblemStatement statement={data.message.getStatement()}
                              limits={data.message.getLimits()}/>}
      </div>
  );
}

type ProblemProps = {
  statement: ProblemPb.ProblemStatement,
  limits: ProblemPb.ProblemLimits,
}

function ProblemStatement({statement, limits}: ProblemProps) {
  return (
      <div className="col-lg-8 problemstatement">
        <Card body={true}>
          <h1 className="problemheader">{statement.getTitle()}</h1>
          <div className="problemlimits mb-2">
            <div>Time limit: {limits.getTimeLimitMs() / 1000} s</div>
            <div>Memory limit: {limits.getMemoryLimitKb() / 1000} MB</div>
          </div>
          <LatexContainer>
            <div dangerouslySetInnerHTML={{__html: statement.getHtml()}}/>
          </LatexContainer>
          <hr/>
          <div>
            License: {statement.getLicense()} |
            Authors: {statement.getAuthorsList().join(", ")}
          </div>
        </Card>
      </div>
  );
}

export default ViewProblem;
