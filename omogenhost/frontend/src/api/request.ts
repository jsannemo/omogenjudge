import {grpc} from "@improbable-eng/grpc-web";


export interface Result<ResponseMessage> {
  status: grpc.Code;
  statusMessage: string;
  headers: grpc.Metadata;
  message: ResponseMessage | null;
  trailers: grpc.Metadata;
}

export type GrpcResult<ResponseMessage> = {
  error: false;
  message: ResponseMessage;
}

export type GrpcError = {
  error: true;
  status: grpc.Code;
  statusMessage: string;
}


const API_PORT = 56744;
const API_ADDRESS = "http://" + window.location.hostname + ":" + API_PORT;

console.log("API address: " + API_ADDRESS);

export async function grpcRequest<TRequest extends grpc.ProtobufMessage,
    ResponseMessage extends grpc.ProtobufMessage,
    M extends grpc.UnaryMethodDefinition<TRequest, ResponseMessage>>(
    methodDescriptor: M,
    request: TRequest,
): Promise<GrpcResult<ResponseMessage> | GrpcError> {
  return new Promise(
      (resolve, reject) => {
        internalRequest(methodDescriptor, request).then((result: Result<ResponseMessage>) => {
          if (result.status == grpc.Code.OK) {
            resolve({
              error: false,
              message: result.message
            });
          } else {
            resolve({
              error: true,
              status: result.status,
              statusMessage: result.statusMessage,
            })
          }
        }).catch(err => reject(err));
      })
}

async function internalRequest<TRequest extends grpc.ProtobufMessage,
    ResponseMessage extends grpc.ProtobufMessage,
    M extends grpc.UnaryMethodDefinition<TRequest, ResponseMessage>>(
    methodDescriptor: M,
    request: TRequest,
): Promise<Result<ResponseMessage>> {
  return new Promise(
      (resolve, reject) => {
        try {
          grpc.unary(methodDescriptor, {
            request: request,
            host: API_ADDRESS,
            onEnd: (response: grpc.UnaryOutput<ResponseMessage>) => {
              resolve({
                status: response.status,
                statusMessage: response.statusMessage,
                headers: response.headers,
                message: response.message,
                trailers: response.trailers,
              });
            }
          });
        } catch (err) {
          reject(err);
        }
      }
  );
}