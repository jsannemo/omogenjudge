import React from "react";
import {SubmitHandler, useForm} from "react-hook-form";
import {Button, Card, Col, Form, FormControl, FormGroup, FormLabel, FormText, Row, Spinner} from "react-bootstrap";
import * as yup from "yup";
import {StringSchema} from "yup";
import {yupResolver} from "@hookform/resolvers/yup";
import {RegisterRequest, RegisterResponse} from "webapi/proto/account_service_pb";
import {AccountService} from "webapi/proto/account_service_pb_service";
import {useApi} from "../../api/react_request";
import { useHistory } from "react-router-dom";
import {addToast} from "../../shared/toasts/Toasts";

const schema = yup.object().shape({
  username: yup.string().label("Username")
    .trim()
    .min(3).max(25)
    .matches(/^[0-9a-zA-Z_\-\\.]*$/, "Only letters, digits, underscores, dashes and dots are allowed")
    .required(),
  full_name: yup.string().label("Full name")
    .trim()
    .min(3).max(255)
    .required(),
  password: yup.string().label("Password")
    .required(),
  password_confirmation: yup.string().label("Password confirmation")
    .when("password", (value: unknown, schema: StringSchema) => schema.equals([value], "Password confirmation doesn't match password")),
  email: yup.string().label("Email")
    .email()
    .required(),
});

export default function RegisterAccount(): JSX.Element {
  const {register, handleSubmit, trigger, formState: {errors}, getValues, setError} = useForm({
    resolver: yupResolver(schema),
  });
  const {state: apiState, send: apiSend} = useApi(AccountService.Register);
  const history = useHistory();
  const onSubmit: SubmitHandler<{
    username: string;
    email: string;
    password: string;
    full_name: string;
  }> = data => {
    const request = new RegisterRequest();
    request.setUsername(data["username"]);
    request.setFullName(data["full_name"]);
    request.setEmail(data["email"]);
    request.setPassword(data["password"]);
    apiSend(request).then((response: RegisterResponse) => {
      if (response.getErrorsList().length === 0) {
        addToast({
          title: "Account Created",
          message: "Your account was created successfully. You can now log in.",
          type: "success"
        });
        history.push("/login");
      }
      response.getErrorsList().forEach(err => {
        switch (err) {
        case RegisterResponse.RegisterError.EMAIL_INVALID:
          setError("email", {"message": "Invalid email"});
          break;
        case RegisterResponse.RegisterError.EMAIL_TAKEN:
          setError("email", {"message": "Email used for another account"});
          break;
        case RegisterResponse.RegisterError.USERNAME_INVALID:
          setError("username", {"message": "Invalid username"});
          break;
        case RegisterResponse.RegisterError.USERNAME_TAKEN:
          setError("username", {"message": "Username taken by another account"});
          break;
        }
      });
    });
  };
  const password = register("password");
  const confirmation = register("password_confirmation");
  return (
    <Row className={"justify-content-center"}>
      <Col lg={"5"}>
        <Card body={true}>
          <h1>Register Account</h1>
          <Form className={"mt-3"} onSubmit={handleSubmit(onSubmit)}>
            <FormGroup className={"mb-4"}>
              <FormLabel>Username *</FormLabel>
              <FormControl
                className={`form-control ${errors.username ? "is-invalid" : ""}`}
                {...register("username")}
              />
              <FormText>
                3-25 characters, letters, digits, underscores, dashes and dots.
              </FormText>
              <FormControl.Feedback type="invalid">
                {errors.username?.message}
              </FormControl.Feedback>
            </FormGroup>
            <FormGroup className={"mb-4"}>
              <FormLabel>Full name *</FormLabel>
              <FormControl
                className={`form-control ${errors.full_name ? "is-invalid" : ""}`}
                {...register("full_name")}
              />
              <FormControl.Feedback type="invalid">
                {errors.full_name?.message}
              </FormControl.Feedback>
            </FormGroup>
            <FormGroup className={"mb-4"}>
              <FormLabel>Password *</FormLabel>
              <FormControl
                type={"password"}
                className={`form-control ${errors.password ? "is-invalid" : ""}`}
                name={"password"}
                onBlur={password.onBlur}
                ref={password.ref}
                onChange={(e) => {
                  password.onChange(e);
                  if (getValues("password_confirmation")) {
                    trigger("password_confirmation");
                  }
                }}
              />
              <FormControl.Feedback type="invalid">
                {errors.password?.message}
              </FormControl.Feedback>
            </FormGroup>
            <FormGroup className={"mb-4"}>
              <FormLabel>Confirm password</FormLabel>
              <FormControl
                type={"password"}
                className={`form-control ${errors.password_confirmation ? "is-invalid" : ""}`}
                name={"password_confirmation"}
                onChange={(e) => {
                  confirmation.onChange(e).then(() => {
                    if (errors.password_confirmation) {
                      trigger("password_confirmation");
                    }
                  });
                }}
                ref={confirmation.ref}
                onBlur={(e: React.FocusEvent<HTMLInputElement>) => {
                  confirmation.onBlur(e);
                  if (getValues("password_confirmation")) {
                    trigger("password_confirmation");
                  }
                }}
              />
              <FormControl.Feedback type="invalid">
                {errors.password_confirmation?.message}
              </FormControl.Feedback>
            </FormGroup>
            <FormGroup className={"mb-4"}>
              <FormLabel>Email *</FormLabel>
              <FormControl
                type={"email"}
                className={`form-control ${errors.email ? "is-invalid" : ""}`}
                {...register("email")}
              />
              <FormControl.Feedback type="invalid">
                {errors.email?.message}
              </FormControl.Feedback>
            </FormGroup>
            {apiState.loading ?
              <Spinner animation={"border"}/> :
              <Button variant="primary" type="submit">
                Register
              </Button>
            }
          </Form>
        </Card>
      </Col>
    </Row>
  );
}