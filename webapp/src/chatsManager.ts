import { Chat } from "./gen/chat/v1/chat_pb";
import { client } from "./client";

export class ChatsManager {
  discussionId: string;
  lowPassword: string;
  private chats: Chat[];
  private callbacks: ((chats: Chat[]) => void)[];

  constructor(discussionId: string, lowPassword: string) {
    this.discussionId = discussionId;
    this.lowPassword = lowPassword;
    this.chats = [];
    this.callbacks = [];

    this.getChats();
    this.watchChatsStreaming();
  }

  async getChats() {
    let resp = await client.getChats({
      discussionId: this.discussionId,
      lowPassword: this.lowPassword,
    });
    this.chats = resp.chats;
    this.onChange();
  }

  async watchChatsStreaming() {
    let iter = client.getChatsStream({
      discussionId: this.discussionId,
      lowPassword: this.lowPassword,
    });
    for await (let resp of iter) {
      this.chats = this.chats.concat(resp.chats);
      this.onChange();
    }
  }

  addEventListener(func: (chats: Chat[]) => void) {
    this.callbacks.push(func);
  }

  onChange() {
    for (let func of this.callbacks) {
      func(this.chats);
    }
  }

  send(name: string, message: string) {
    return client.sendChat({
      discussionId: this.discussionId,
      name: name,
      message: message,
      lowPassword: this.lowPassword,
    });
  }
}
