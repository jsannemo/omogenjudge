import React from 'react';
import {BrowserRouter as Router, Route, Switch} from 'react-router-dom';
import Menu from "./shared/menu/Menu";
import ViewProblem from "./problems/view/ViewProblem";
import Page from "./shared/container/Page";

import './App.scss';

function App() {
  return (
      <Router>
        <div className={"page-wrapper bg-light"}>
          <Menu/>
          <Page>
            <Switch>
              <Route path={"/problems/:shortname"}>
                <ViewProblem/>
              </Route>
            </Switch>
          </Page>
        </div>
      </Router>
  );
}

export default App;