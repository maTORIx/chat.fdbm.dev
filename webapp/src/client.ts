import {
  createConnectTransport,
  createPromiseClient,
} from "@bufbuild/connect-web";
import { ChatService } from "./gen/chat/v1/chat_connectweb";

const transport = createConnectTransport({
  baseUrl: "http://localhost:8080",
});

export const client = createPromiseClient(ChatService, transport);