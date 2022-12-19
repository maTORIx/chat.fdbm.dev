import CryptoJS from "crypto-js"
import { Chat } from "./gen/chat/v1/chat_pb";
import { client } from "./client";

const SALT = "rnP3dHk$0zS.?'!jkr)ytjjj"

export class ChatsManager {
  discussionId: string;
  private lowPassword: string;
  private hashedPassword: string;
  private chats: Chat[];
  private callbacks: ((chats: Chat[]) => void)[];

  constructor(discussionId: string, lowPassword: string) {
    this.discussionId = discussionId;
    this.lowPassword = lowPassword;
    this.hashedPassword = CryptoJS.SHA256(this.lowPassword).toString();
    this.chats = [];
    this.callbacks = [];
  }

  async getChats(limit = 20) {
    let pageingInfo = { limit: limit, cursor: "", earlierAt: 0 }
    if (this.chats.length > 0) {
      pageingInfo.cursor = this.chats[0].id
      pageingInfo.earlierAt = this.chats[0].createdAt
    }
    let resp = await client.getChats({
      discussionInfo: {
        id: this.discussionId,
        lowPassword: this.hashedPassword,
      },
      pageingInfo: pageingInfo
    });
    resp.chats.forEach((val, idx) => val.body = this.decryptMessage(val.body))
    resp.chats = resp.chats.reverse()
    this.chats = resp.chats.concat(this.chats);
    this.onChange();
    return this.chats
  }


  encryptMessage(message: string): string {
    return CryptoJS.AES.encrypt(message, this.hashedPassword + SALT).toString()
  }

  decryptMessage(message: string): string {
    return CryptoJS.AES.decrypt(message, this.hashedPassword + SALT).toString(CryptoJS.enc.Utf8)
  }

  async watchChatsStreaming() {
    try {
      let iter = client.getChatsStream({
        discussionInfo: {
          id: this.discussionId,
          lowPassword: this.hashedPassword,
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
    } catch(e) {
      console.error("Streaming disconnected.")
      window.setTimeout(this.watchChatsStreaming.bind(this), 5000)
      console.info("Retrying...")
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
        lowPassword: this.hashedPassword
      },
      name: name,
      body: this.encryptMessage(message),
    });
  }
}
