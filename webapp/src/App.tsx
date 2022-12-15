import { Value } from "@bufbuild/protobuf";
import { ChangeEvent, useState, useEffect, useRef, useCallback, Reducer } from "react";
import { ChatsManager } from "./chatsManager";
import { client } from "./client";
import { Chat } from "./gen/chat/v1/chat_pb";
import "./styles/fonts.css";
import "./styles/reset.css";
import "./styles/style.css";

let chatsManager: ChatsManager;

let userStylesDict: {
  [user_id: string]: {
    color: string
  }
} = {};

function App() {
  const [name, setName] = useState("");
  const [message, setMessage] = useState("");
  const [chats, setChats] = useState<Chat[]>([]);
  const [isScrolled, setIsScrolled] = useState(false);
  const [isLatestButtonSuggested, setIsLatestButtonSuggested] = useState(false)
  const [firsttime, setFirstrtime] = useState(true)
  const [lastScrollHeight, setLastScrollHeight] = useState(0)

  const url = new URL(window.location.href);
  const discussionId = url.searchParams.get("d");
  if (discussionId === null) {
    window.location.href = url.href + "?d=" + generateRandomString(8);
    return <div />;
  }

  useEffect(() => {
    let firsttime = true;
    if (discussionId === null || !!chatsManager) return;
    chatsManager = new ChatsManager(discussionId, "");
    chatsManager.addEventListener((chats: Chat[]) => {
      setIsScrolled(false)
      setChats(chats);
    });
    chatsManager.watchChatsStreaming()
    chatsManager.getChats()
  }, []);

  const scrollToLatest = useCallback(() => {
    const container = document.getElementById("chats-container")
    if (!container) return;
    container.scroll(0, container.scrollHeight - container.clientHeight)
  }, [])

  useEffect(() => {
    const conn = document.getElementById("chats-container")
    const latestChat = chats.length > 0 ? document.getElementById(`chat-${chats.slice(-1)[0].id}`) : null

    if (!isScrolled && conn && latestChat && conn.scrollTop !== conn.scrollHeight - conn.clientHeight) {
      const inspectedHeight = conn.scrollHeight - latestChat.clientHeight - conn.clientHeight
      const isInspectedBottom = inspectedHeight <= conn.scrollTop
      if (firsttime || isInspectedBottom) {
        scrollToLatest()
        if (firsttime) setFirstrtime(false)
      } else if (conn.scrollTop === 0) {
        conn.scroll(0, conn.scrollHeight - lastScrollHeight)
      } else {
        setIsLatestButtonSuggested(true);
      }
      setLastScrollHeight(conn.scrollHeight)
      setIsScrolled(true)
    }
    return () => {}
  }, [chats, isScrolled])

  // set chats
  let chatElements = [];
  for (let chat of chats) {
    chatElements.push(<ChatElement key={chat.id} chat={chat} />);
  }

  const onSubmit = useCallback(async (e: { preventDefault: () => void; }) => {
    e.preventDefault();
    if (message === "" || !chatsManager) return;
    await chatsManager.send(name ? name : "unknown", message);
    setMessage("");
  }, [message, name])

  const onScrollChats = useCallback(async (e: any) => {
    if (e.target.scrollTop == e.target.scrollHeight - e.target.clientHeight) {
      setIsLatestButtonSuggested(false)
    } else if (e.target.scrollTop === 0) {
      chatsManager.getChats()
    }
  }, [])

  return (
    <div className="App">
      <div className="top-bar">
        <img src="/chatbird.svg" alt="C" />
        <label className="profile">
          <input name="name" onChange={withEventValue(setName)} value={name} placeholder="unknown" maxLength={14} />
          <span className="account-icon material-symbols-outlined">person</span>
        </label>
      </div>
      <div className="chats-container" id="chats-container" onScroll={onScrollChats}>
        {chatElements}
      </div>
      <div className={`suggest-latest ${isLatestButtonSuggested ? "" : "hidden"}`} onClick={() => {
        setIsLatestButtonSuggested(false)
        scrollToLatest()
      }}>
        <p>
          <span className="material-symbols-outlined">south</span>
          新しいチャットが届いています
        </p>
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
          maxLength={5000}
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

function generateRandomColorFromText(str: string): string {
  if (str.length < 1) {
    return "#ffffff";
  }
  while (str.length < 6) {
    str += str
  }
  let result = [0, 0, 0, 0, 0, 0]
  str.split("").forEach((v, idx) => result[idx % 6] += v.charCodeAt(0))
  return "#" + result.map((v) => (v % 16).toString(16)).join("")
}

const ChatElement = (props: { chat: Chat }) => {
  const chat = props.chat
  if (!userStylesDict[chat.userId]) {
    userStylesDict[chat.userId] = {
      color: generateRandomColorFromText(chat.userId)
    }
  }

  return (
    <div id={`chat-${chat.id}`} className="chat-container">
      <div className="account-icon" style={{ "backgroundColor": userStylesDict[chat.userId].color }}>
        <p>person</p>
      </div>
      <div className="chat-body">
        <p className="name">{chat.name}</p>
        <p className="message">{chat.body}</p>
      </div>
      <div className="time">
        <time>{new Date(chat.createdAt).toLocaleTimeString().slice(0, -3)}</time>
      </div>
    </div>
  )
}

export default App;
