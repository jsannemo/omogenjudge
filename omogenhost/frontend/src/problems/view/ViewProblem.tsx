import React from "react";

import {ViewProblemRequest, ViewProblemResponse} from "webapi/proto/problem_service_pb";
import ProblemPb from "webapi/proto/problems_pb";
import {ProblemService} from "webapi/proto/problem_service_pb_service";
import {ReactResponse, useRequest} from "../../api/react_request";
import {useParams} from "react-router-dom";
import {Card, Col, Row} from "react-bootstrap";

import "./ViewProblem.scss";
import {LatexContainer} from "../../shared/latex/LatexContainer";

export default function ViewProblem(): JSX.Element {
  const {shortname} = useParams();
  const request = new ViewProblemRequest();
  request.setShortName(shortname);
  const data: ReactResponse<ViewProblemResponse> = useRequest(ProblemService.ViewProblem, request);
  return (
    <Row>
      {data.loading ? "Loading problem..." :
        data.error ? "Could not load problem" :
          <ProblemStatement statement={data.message.getStatement()!}
            limits={data.message.getLimits()!}/>}
    </Row>
  );
}

type ProblemProps = {
    statement: ProblemPb.ProblemStatement,
    limits: ProblemPb.ProblemLimits,
}

function ProblemStatement({statement, limits}: ProblemProps): JSX.Element {
  return (
    <Col lg={8} className={"problemstatement"}>
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
            License: {statement.getLicense()} | Authors: {statement.getAuthorsList().join(", ")}
        </div>
      </Card>
    </Col>
  );
}