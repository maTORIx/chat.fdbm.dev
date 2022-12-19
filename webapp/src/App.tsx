import { ChangeEvent, useState, useEffect, useCallback} from "react";
import { Chat } from "./gen/chat/v1/chat_pb";
import { ChatsManager } from "./chatsManager";
import { client } from "./client";
import { PasswordModal } from "./components/passwordModal";
import { generateRandomString, withEventValue } from "./utils";
import "./styles/fonts.css";
import "./styles/reset.css";
import "./styles/style.css";
import { ChatListElement } from "./components/chats";

let chatsManager: ChatsManager;



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
    let lowPassword = window.localStorage.getItem(`${discussionId}-lowPassword`)
    if (lowPassword === null) lowPassword = ""
    chatsManager = new ChatsManager(discussionId, lowPassword);
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

  // scroll suggest or scroll ajustment function
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
    return () => { }
  }, [chats, isScrolled])


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
      {/* Topbar */}
      <div className="top-bar">
        <img src="/chatbird.svg" alt="C" />
        <div className="rightside">
          <label className="profile">
            <input name="name" onChange={withEventValue(setName)} value={name} placeholder="unknown" maxLength={14} />
            <span className="account-icon material-symbols-outlined">person</span>
          </label>
          <PasswordModal onSubmit={(lowPassword) => {
            window.localStorage.setItem(`${discussionId}-lowPassword`, lowPassword)
            location.reload()
          }}/>
        </div>
      </div>
      <ChatListElement onScroll={onScrollChats} chats={chats}/>
      {/* Suggest Latest Button */}
      <div className={`suggest-latest"}`} hidden={isLatestButtonSuggested} onClick={() => {
        setIsLatestButtonSuggested(false)
        scrollToLatest()
      }}>
        <p>
          <span className="material-symbols-outlined">south</span>
          新しいチャットが届いています
        </p>
      </div>
      {/* form */}
      <form
        className="chat-form"
        onSubmit={onSubmit}
        placeholder="Aa"
        style={{height: `calc(${message.split("\n").length - 1} * 24px + 40px)`}}
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

export default App;
