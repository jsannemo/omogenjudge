import {grpc} from "@improbable-eng/grpc-web";
import {GrpcError, grpcRequest} from "./request";
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
export type ReactResponse<ResponseMessage> = ReactLoading | ReactError | ReactMessage<ResponseMessage>;

type SetResponse<ResponseMessage> = (r: ReactResponse<ResponseMessage>) => void;

export function useRequest<TRequest extends grpc.ProtobufMessage,
    ResponseMessage extends grpc.ProtobufMessage,
    M extends grpc.UnaryMethodDefinition<TRequest, ResponseMessage>>(
    methodDescriptor: M,
    request: TRequest,
): ReactResponse<ResponseMessage> {
    const initialState: ReactResponse<ResponseMessage> = {loading: true, message: null, error: null};
    const [state, setState] = useState<ReactResponse<ResponseMessage>>(initialState);
    useEffect(() => {
        (async () => {
            const response = await grpcRequest<TRequest, ResponseMessage, M>(methodDescriptor, request);
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