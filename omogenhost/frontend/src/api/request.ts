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

export type GrpcResponse<ResponseMessage> = GrpcError | GrpcResult<ResponseMessage>;

const API_PORT = 3000;
const API_ADDRESS = "http://" + window.location.hostname + ":" + API_PORT;

console.log("API address: " + API_ADDRESS);

export type GrpcOpts = {
  cache?: boolean;
}

type Cache = Map<grpc.UnaryMethodDefinition<grpc.ProtobufMessage, grpc.ProtobufMessage>, Record<string, unknown>>;
const cache: Cache = new Map<grpc.UnaryMethodDefinition<grpc.ProtobufMessage, grpc.ProtobufMessage>, Record<string, unknown>>();

export function checkCache<TRequest extends grpc.ProtobufMessage,
  ResponseMessage extends grpc.ProtobufMessage,
  M extends grpc.UnaryMethodDefinition<TRequest, ResponseMessage>>(
  methodDescriptor: M,
  request: TRequest,
): ResponseMessage | null {
  const cacheKey = JSON.stringify(request.toObject());
  const rCache = cache.get(methodDescriptor);
  if (!rCache) {
    return null;
  }
  const cached = rCache[cacheKey];
  if (!cached) {
    return null;
  }
  return cached as ResponseMessage;
}

export async function grpcRequest<TRequest extends grpc.ProtobufMessage,
  ResponseMessage extends grpc.ProtobufMessage,
  M extends grpc.UnaryMethodDefinition<TRequest, ResponseMessage>>(
  methodDescriptor: M,
  request: TRequest,
  options?: GrpcOpts,
): Promise<GrpcResponse<ResponseMessage>> {
  let cacheKey: string | undefined;
  let rCache: Record<string, unknown> | undefined;
  const opts = options || {};
  if (opts.cache === true) {
    cacheKey = JSON.stringify(request.toObject());
    rCache = cache.get(methodDescriptor);
    if (!rCache) {
      cache.set(methodDescriptor, rCache = {});
    }
    if (rCache && cacheKey in rCache) {
      const msg = rCache[cacheKey] as ResponseMessage;
      return {
        error: false,
        message: msg,
      };
    }
  }

  const result = await internalRequest(methodDescriptor, request);
  const authHeaders = result.headers.get("authorization");
  if (authHeaders.length) {
    localStorage.setItem("authorization", authHeaders[0]);
  } else {
    localStorage.removeItem("authorization");
  }
  if (result.status == grpc.Code.OK) {
    if (rCache) {
      rCache[cacheKey!] = result.message;
    }
    return {
      error: false,
      message: result.message as ResponseMessage,
    };
  } else {
    return {
      error: true,
      status: result.status,
      statusMessage: result.statusMessage,
    };
  }
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
          metadata: {
            authorization: localStorage.getItem("authorization") || [],
          },
          onEnd: (response: grpc.UnaryOutput<ResponseMessage>) => {
            const result = {
              status: response.status,
              statusMessage: response.statusMessage,
              headers: response.headers,
              message: response.message,
              trailers: response.trailers,
            };
            resolve(result);
          }
        });
      } catch (err) {
        reject(err);
      }
    }
  );
}