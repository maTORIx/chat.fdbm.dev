import CryptoJS from "crypto-js"
import { Chat } from "./gen/chat/v1/chat_pb";
import { client } from "./client";

const SALT = "rnP3dHk$0zS.?'!jkr)ytjjj"

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
      discussionInfo: {
        id: this.discussionId,
        lowPassword: this.lowPassword,
      }
    });
    resp.chats.forEach((val, idx) => val.body = this.decryptMessage(val.body))
    this.chats = resp.chats;
    this.onChange();
  }

  encryptMessage(message: string): string {
    return CryptoJS.AES.encrypt(message, this.lowPassword + SALT).toString()
  }

  decryptMessage(message: string): string {
    return CryptoJS.AES.decrypt(message, this.lowPassword + SALT).toString(CryptoJS.enc.Utf8)
  }

  async watchChatsStreaming() {
    let iter = client.getChatsStream({
      discussionInfo: {
        id: this.discussionId,
        lowPassword: this.lowPassword,
      }
    });
    for await (let resp of iter) {
      if (!resp.chat) {
        continue
      }
      resp.chat.body = this.decryptMessage(resp.chat.body)
      this.chats.push(resp.chat);
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
      discussionInfo: {
        id: this.discussionId,
        lowPassword: this.lowPassword
      },
      name: name,
      body: this.encryptMessage(message),
    });
  }
}
