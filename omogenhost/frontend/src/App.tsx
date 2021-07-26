import React from "react";
import {BrowserRouter as Router, Route, Switch} from "react-router-dom";
import Menu from "./shared/menu/Menu";
import ViewProblem from "./problems/view/ViewProblem";
import Page from "./shared/container/Page";
import "./App.scss";
import RegisterAccount from "./accounts/register/RegisterAccount";
import {ToastComponent} from "./shared/toasts/Toasts";
import Login from "./accounts/login/Login";

export default function App(): JSX.Element {
  return (
    <Router>
      <div className={"page-wrapper bg-light"}>
        <Menu/>
        <Page>
          <Switch>
            <Route path={"/problems/:shortname"}>
              <ViewProblem/>
            </Route>
            <Route path={"/login"}>
              <Login/>
            </Route>
            <Route path={"/register"}>
              <RegisterAccount/>
            </Route>
          </Switch>
        </Page>
      </div>
      <ToastComponent/>
    </Router>
  );
}