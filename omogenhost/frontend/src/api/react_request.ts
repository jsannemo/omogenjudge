import {grpc} from "@improbable-eng/grpc-web";
import {checkCache, GrpcError, GrpcOpts, grpcRequest} from "./request";
import {useEffect, useState} from "react";

export type ReactLoading = {
  loading: true;
  message: null;
  error: null;
}

export type ReactMessage<ResponseMessage> = {
  loading: false;
  message: ResponseMessage;
  error: null;
}

export type ReactError = {
  loading: false;
  message: null;
  error: GrpcError;
}

export type ReactApiState<ResponseMessage> = {
  loading: boolean;
  message: ResponseMessage | null;
  error: GrpcError | null
}

export type ReactApi<RequestMessage, ResponseMessage> = {
  state: ReactApiState<ResponseMessage>;
  send: (m: RequestMessage) => Promise<ResponseMessage>;
};

export function useApi<TRequest extends grpc.ProtobufMessage,
  ResponseMessage extends grpc.ProtobufMessage,
  M extends grpc.UnaryMethodDefinition<TRequest, ResponseMessage>>(
  methodDescriptor: M,
): ReactApi<TRequest, ResponseMessage> {
  const initialState: ReactApiState<ResponseMessage> = {loading: false, message: null, error: null};
  const [state, setState] = useState(initialState);
  const send = async (request: TRequest) => {
    console.debug("Sending request: ", request);
    setState({...state, loading: true});
    const response = await grpcRequest<TRequest, ResponseMessage, M>(methodDescriptor, request);
    console.debug("Response:", response);
    if (response.error === true) {
      setState({
        loading: false,
        error: response,
        message: null,
      });
      throw response;
    } else {
      setState({
        loading: false,
        message: response.message,
        error: null,
      });
      return response.message;
    }
  };
  return {state, send};
}

export type ReactResponse<ResponseMessage> = ReactLoading | ReactError | ReactMessage<ResponseMessage>;

export function useRequest<TRequest extends grpc.ProtobufMessage,
  ResponseMessage extends grpc.ProtobufMessage,
  M extends grpc.UnaryMethodDefinition<TRequest, ResponseMessage>>(
  methodDescriptor: M,
  request: TRequest,
  options?: GrpcOpts,
): ReactResponse<ResponseMessage> {
  let initialState: ReactResponse<ResponseMessage> = {loading: true, message: null, error: null};
  console.debug("Sending request: ", request);
  const shouldCache = options && options.cache;
  if (shouldCache) {
    const cached: ResponseMessage | null = checkCache(methodDescriptor, request);
    if (cached) {
      console.debug("Hit cache", request);
      initialState = {loading: false, message: cached, error: null};
    }
  }
  const [state, setState] = useState<ReactResponse<ResponseMessage>>(initialState);
  useEffect(() => {
    (async () => {
      if (state.loading) {
        const response = await grpcRequest<TRequest, ResponseMessage, M>(methodDescriptor, request, options);
        console.debug("Response:", response);
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
      }
    })();
  }, []);
  return state;
}