import React, { useEffect } from "react";
import {useHistory} from "react-router-dom";
import {useApi} from "../../api/react_request";
import {LogoutRequest} from "webapi/proto/account_service_pb";
import {AccountService} from "webapi/proto/account_service_pb_service";

export default function Logout(): JSX.Element {
  const history = useHistory();
  const {send: apiSend} = useApi(AccountService.Logout);
  useEffect(() => {
    apiSend(new LogoutRequest()).then(() => {
      history.push("/");
    });
  }, []);
  return (
    <div>Logging out...</div>
  );
}