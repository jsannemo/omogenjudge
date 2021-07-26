import {getProfile, getUserId} from "./accounts/auth";

require("google-protobuf");
import React from "react";
import ReactDOM from "react-dom";
import App from "./App";

import "./index.scss";
import {grpcRequest} from "./api/request";
import { AccountService } from "webapi/proto/account_service_pb_service";
import {LogoutRequest} from "webapi/proto/account_service_pb";

async function initialize() {
  if (getUserId() != null && getProfile() == null) {
    await grpcRequest(AccountService.Logout, new LogoutRequest());
  }
  ReactDOM.render(
    <App/>,
    document.getElementById("root") as HTMLElement
  );
}

initialize().then(() => {});