import {grpc} from "@improbable-eng/grpc-web";
import {GrpcError, grpcRequest} from "./request";
import {useEffect, useState} from "react";

export interface ReactRequest<ResponseMessage> {
  loading: boolean;
  message: ResponseMessage | null;
  error: GrpcError | null;
}

export function useRequest<TRequest extends grpc.ProtobufMessage,
    ResponseMessage extends grpc.ProtobufMessage,
    M extends grpc.UnaryMethodDefinition<TRequest, ResponseMessage>>(
    methodDescriptor: M,
    request: TRequest,
): ReactRequest<ResponseMessage> {
  const [state, setState] = useState({loading: true, message: null, error: null});
  useEffect(() => {
    (async () => {
      const response = await grpcRequest(methodDescriptor, request);
      if (response.error === true) {
        setState({
          loading: false,
          error: response,
          message: null,
        });
      } else if (response.error === false) {
        setState({
          loading: false,
          message: response.message,
          error: null,
        });
      }
    })();
  }, []);
  return state;
}