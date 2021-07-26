import React from "react";
import {SubmitHandler, useForm} from "react-hook-form";
import {Button, Card, Col, Form, FormControl, FormGroup, FormLabel, Row, Spinner} from "react-bootstrap";
import * as yup from "yup";
import {yupResolver} from "@hookform/resolvers/yup";
import {AccountService} from "webapi/proto/account_service_pb_service";
import {LoginRequest, LoginResponse} from "webapi/proto/account_service_pb";
import {useApi} from "../../api/react_request";
import {useHistory} from "react-router-dom";

const schema = yup.object().shape({
  username: yup.string().label("Username").trim(),
  password: yup.string().label("Password")
});

export default function Login(): JSX.Element {
  const {register, handleSubmit, formState: {errors}, getValues, setError} = useForm({
    resolver: yupResolver(schema),
  });
  const {state: apiState, send: apiSend} = useApi(AccountService.Login);
  const history = useHistory();
  const onSubmit: SubmitHandler<{
    username: string;
    password: string;
  }> = data => {
    const request = new LoginRequest();
    request.setUsername(data["username"]);
    request.setPassword(data["password"]);
    apiSend(request).then((response: LoginResponse) => {
      if (response.getErrorsList().length === 0) {
        history.push("/");
      }
      response.getErrorsList().forEach(err => {
        switch (err) {
        case LoginResponse.LoginError.INVALID_CREDENTIALS:
          setError("password", {"message": "The username and password combination is not correct"});
          break;
        }
      });
    });
  };
  return (
    <Row className={"justify-content-center"}>
      <Col lg={"5"}>
        <Card body={true}>
          <h1>Login</h1>
          <Form className={"mt-3"} onSubmit={handleSubmit(onSubmit)}>
            <FormGroup className={"mb-4"}>
              <FormLabel>Username</FormLabel>
              <FormControl {...register("username")}/>
            </FormGroup>
            <FormGroup className={"mb-4"}>
              <FormLabel>Password</FormLabel>
              <FormControl
                type={"password"}
                className={`form-control ${errors.password ? "is-invalid" : ""}`}
                {...register("password")}
              />
              <FormControl.Feedback type="invalid">
                {errors.password?.message}
              </FormControl.Feedback>
            </FormGroup>
            {apiState.loading ?
              <Spinner animation={"border"}/> :
              <Button variant="primary" type="submit">
                Login
              </Button>
            }
          </Form>
        </Card>
      </Col>
    </Row>
  );
}