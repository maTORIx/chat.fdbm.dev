import { Value } from "@bufbuild/protobuf";
import { ChangeEvent, useState, useEffect, useRef, useCallback } from "react";
import { ChatsManager } from "./chatsManager";
import { Chat } from "./gen/chat/v1/chat_pb";
import "./styles/fonts.css";
import "./styles/reset.css";
import "./styles/style.css";

let chatsManager: ChatsManager;

function App() {
  const [name, setName] = useState("");
  const [message, setMessage] = useState("");
  const [chats, setChats] = useState<Chat[]>([]);
  const [isScrolled, setIsScrolled] = useState(false);

  const url = new URL(window.location.href);
  const discussionId = url.searchParams.get("d");
  if (discussionId === null) {
    window.location.href = url.href + "?d=" + generateRandomString(8);
    return <div />;
  }

  useEffect(() => {
    if (discussionId === null || !!chatsManager) return;
    chatsManager = new ChatsManager(discussionId, "");
    chatsManager.addEventListener((chats: Chat[]) => {
      setChats(chats);
      setIsScrolled(false)
    });
  }, []);


  const scrollToLatest = useCallback(() => {
    const chats = document.getElementById("chats-container")
    if (!chats) return;
    chats.scroll(0, chats.scrollHeight)
  }, [])

  const whenChatsUpdated = useEffect(() => {
    const container = document.getElementById("chats-container")
    const lastChatElem = document.getElementById(`chat-${chats.slice(-2, -1)[0].id}`)
    if (!container || isScrolled || !lastChatElem) {
    } else if (container.scrollHeight === container.scrollTop + container.clientHeight - lastChatElem.clientHeight) {
      scrollToLatest()
    }
    setIsScrolled(false)
  }, [chats])

  let chatElements = [];
  for (let chat of chats) {
    chatElements.push(
      <div key={chat.id} id={`chat-${chat.id}`} className="chat-container">
        <div className="account-icon">
          <p>person</p>
        </div>
        <div className="chat-body">
          <p className="name">{chat.name}</p>
          <p className="message">{chat.message}</p>
        </div>
        <div className="time">
          <time>{new Date(parseInt(chat.createdAt)).toLocaleTimeString().slice(0, -3)}</time>
        </div>
      </div>
    );
  }

  const onSubmit = useCallback(async (e: { preventDefault: () => void; }) => {
    e.preventDefault();
    if (message === "" || !chatsManager) return;
    await chatsManager.send(name ? name : "unknown", message);
    setMessage("");
  }, [message, name])

  return (
    <div className="App">
      <div className="top-bar">
        <img src="/chatbird.svg" alt="C" />
        <label className="profile">
          <input name="name" onChange={withEventValue(setName)} value={name} placeholder="unknown" maxLength={14} />
          <span className="account-icon material-symbols-outlined">person</span>
        </label>
      </div>
      <div className="chats-container" id="chats-container">
        {chatElements}
      </div>
      <form
        className="chat-form"
        onSubmit={onSubmit}
        placeholder="Aa"
      >
        <textarea
          name="message"
          onChange={withEventValue(setMessage)}
          value={message}
        ></textarea>
        <button className="material-symbols-outlined submit">send</button>
      </form>
      <small>© 2022 Takumi Jonen</small>
    </div>
  );
}

function withEventValue(f: Function): (e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => void {
  return (e) => f(e.target.value);
}

function generateRandomString(len: Number): String {
  const charList =
    "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
  return Array.from(Array(len))
    .map(() => charList[Math.floor(Math.random() * charList.length)])
    .join("");
}

export default App;
