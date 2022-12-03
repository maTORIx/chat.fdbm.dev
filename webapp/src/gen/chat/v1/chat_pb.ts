// @generated by protoc-gen-es v0.3.0 with parameter "target=ts"
// @generated from file chat/v1/chat.proto (package chat.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";

/**
 * @generated from message chat.v1.GreetRequest
 */
export class GreetRequest extends Message<GreetRequest> {
  /**
   * @generated from field: string name = 1;
   */
  name = "";

  constructor(data?: PartialMessage<GreetRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime = proto3;
  static readonly typeName = "chat.v1.GreetRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GreetRequest {
    return new GreetRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GreetRequest {
    return new GreetRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GreetRequest {
    return new GreetRequest().fromJsonString(jsonString, options);
  }

  static equals(a: GreetRequest | PlainMessage<GreetRequest> | undefined, b: GreetRequest | PlainMessage<GreetRequest> | undefined): boolean {
    return proto3.util.equals(GreetRequest, a, b);
  }
}

/**
 * @generated from message chat.v1.GreetResponse
 */
export class GreetResponse extends Message<GreetResponse> {
  /**
   * @generated from field: string greeting = 1;
   */
  greeting = "";

  constructor(data?: PartialMessage<GreetResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime = proto3;
  static readonly typeName = "chat.v1.GreetResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "greeting", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GreetResponse {
    return new GreetResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GreetResponse {
    return new GreetResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GreetResponse {
    return new GreetResponse().fromJsonString(jsonString, options);
  }

  static equals(a: GreetResponse | PlainMessage<GreetResponse> | undefined, b: GreetResponse | PlainMessage<GreetResponse> | undefined): boolean {
    return proto3.util.equals(GreetResponse, a, b);
  }
}

/**
 * @generated from message chat.v1.SendChatRequest
 */
export class SendChatRequest extends Message<SendChatRequest> {
  /**
   * @generated from field: string discussion_id = 1;
   */
  discussionId = "";

  /**
   * @generated from field: string name = 2;
   */
  name = "";

  /**
   * @generated from field: string message = 3;
   */
  message = "";

  /**
   * @generated from field: string low_password = 4;
   */
  lowPassword = "";

  constructor(data?: PartialMessage<SendChatRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime = proto3;
  static readonly typeName = "chat.v1.SendChatRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "discussion_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "message", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 4, name: "low_password", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): SendChatRequest {
    return new SendChatRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): SendChatRequest {
    return new SendChatRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): SendChatRequest {
    return new SendChatRequest().fromJsonString(jsonString, options);
  }

  static equals(a: SendChatRequest | PlainMessage<SendChatRequest> | undefined, b: SendChatRequest | PlainMessage<SendChatRequest> | undefined): boolean {
    return proto3.util.equals(SendChatRequest, a, b);
  }
}

/**
 * @generated from message chat.v1.SendChatResponse
 */
export class SendChatResponse extends Message<SendChatResponse> {
  constructor(data?: PartialMessage<SendChatResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime = proto3;
  static readonly typeName = "chat.v1.SendChatResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): SendChatResponse {
    return new SendChatResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): SendChatResponse {
    return new SendChatResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): SendChatResponse {
    return new SendChatResponse().fromJsonString(jsonString, options);
  }

  static equals(a: SendChatResponse | PlainMessage<SendChatResponse> | undefined, b: SendChatResponse | PlainMessage<SendChatResponse> | undefined): boolean {
    return proto3.util.equals(SendChatResponse, a, b);
  }
}

/**
 * @generated from message chat.v1.GetChatsRequest
 */
export class GetChatsRequest extends Message<GetChatsRequest> {
  /**
   * @generated from field: string discussion_id = 1;
   */
  discussionId = "";

  /**
   * @generated from field: string low_password = 2;
   */
  lowPassword = "";

  constructor(data?: PartialMessage<GetChatsRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime = proto3;
  static readonly typeName = "chat.v1.GetChatsRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "discussion_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "low_password", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetChatsRequest {
    return new GetChatsRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetChatsRequest {
    return new GetChatsRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetChatsRequest {
    return new GetChatsRequest().fromJsonString(jsonString, options);
  }

  static equals(a: GetChatsRequest | PlainMessage<GetChatsRequest> | undefined, b: GetChatsRequest | PlainMessage<GetChatsRequest> | undefined): boolean {
    return proto3.util.equals(GetChatsRequest, a, b);
  }
}

/**
 * @generated from message chat.v1.Chat
 */
export class Chat extends Message<Chat> {
  /**
   * @generated from field: string id = 1;
   */
  id = "";

  /**
   * @generated from field: string name = 3;
   */
  name = "";

  /**
   * @generated from field: string message = 4;
   */
  message = "";

  /**
   * @generated from field: string created_at = 5;
   */
  createdAt = "";

  constructor(data?: PartialMessage<Chat>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime = proto3;
  static readonly typeName = "chat.v1.Chat";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 4, name: "message", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 5, name: "created_at", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Chat {
    return new Chat().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Chat {
    return new Chat().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Chat {
    return new Chat().fromJsonString(jsonString, options);
  }

  static equals(a: Chat | PlainMessage<Chat> | undefined, b: Chat | PlainMessage<Chat> | undefined): boolean {
    return proto3.util.equals(Chat, a, b);
  }
}

/**
 * @generated from message chat.v1.GetChatsResponse
 */
export class GetChatsResponse extends Message<GetChatsResponse> {
  /**
   * @generated from field: repeated chat.v1.Chat chats = 1;
   */
  chats: Chat[] = [];

  constructor(data?: PartialMessage<GetChatsResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime = proto3;
  static readonly typeName = "chat.v1.GetChatsResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "chats", kind: "message", T: Chat, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetChatsResponse {
    return new GetChatsResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetChatsResponse {
    return new GetChatsResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetChatsResponse {
    return new GetChatsResponse().fromJsonString(jsonString, options);
  }

  static equals(a: GetChatsResponse | PlainMessage<GetChatsResponse> | undefined, b: GetChatsResponse | PlainMessage<GetChatsResponse> | undefined): boolean {
    return proto3.util.equals(GetChatsResponse, a, b);
  }
}

/**
 * @generated from message chat.v1.GetChatsStreamRequest
 */
export class GetChatsStreamRequest extends Message<GetChatsStreamRequest> {
  /**
   * @generated from field: string discussion_id = 1;
   */
  discussionId = "";

  /**
   * @generated from field: string low_password = 2;
   */
  lowPassword = "";

  constructor(data?: PartialMessage<GetChatsStreamRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime = proto3;
  static readonly typeName = "chat.v1.GetChatsStreamRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "discussion_id", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "low_password", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetChatsStreamRequest {
    return new GetChatsStreamRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetChatsStreamRequest {
    return new GetChatsStreamRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetChatsStreamRequest {
    return new GetChatsStreamRequest().fromJsonString(jsonString, options);
  }

  static equals(a: GetChatsStreamRequest | PlainMessage<GetChatsStreamRequest> | undefined, b: GetChatsStreamRequest | PlainMessage<GetChatsStreamRequest> | undefined): boolean {
    return proto3.util.equals(GetChatsStreamRequest, a, b);
  }
}

/**
 * @generated from message chat.v1.GetChatsStreamResponse
 */
export class GetChatsStreamResponse extends Message<GetChatsStreamResponse> {
  /**
   * @generated from field: repeated chat.v1.Chat chats = 1;
   */
  chats: Chat[] = [];

  constructor(data?: PartialMessage<GetChatsStreamResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime = proto3;
  static readonly typeName = "chat.v1.GetChatsStreamResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "chats", kind: "message", T: Chat, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GetChatsStreamResponse {
    return new GetChatsStreamResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GetChatsStreamResponse {
    return new GetChatsStreamResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GetChatsStreamResponse {
    return new GetChatsStreamResponse().fromJsonString(jsonString, options);
  }

  static equals(a: GetChatsStreamResponse | PlainMessage<GetChatsStreamResponse> | undefined, b: GetChatsStreamResponse | PlainMessage<GetChatsStreamResponse> | undefined): boolean {
    return proto3.util.equals(GetChatsStreamResponse, a, b);
  }
}

