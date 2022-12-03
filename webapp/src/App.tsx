import { ChangeEvent, useState, useEffect, useRef } from "react";
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

  const url = new URL(window.location.href);
  const discussionId = url.searchParams.get("d");
  if (discussionId === null) {
    window.location.href = url.href + "?d=" + generateRandomString(8);
    return <div />;
  }

  useEffect(() => {
    if (discussionId == null || chatsManager != null) return;
    let ignore = false;
    chatsManager = new ChatsManager(discussionId, "");
    chatsManager.addEventListener((chats: Chat[]) => {
      setChats(chats);
    });
  }, []);

  let chatElements = [];
  for (let chat of chats) {
    chatElements.push(
      <div key={chat.id}>
        <p>
          <b>{chat.name}</b>
        </p>
        <p>{chat.message}</p>
      </div>
    );
  }

  return (
    <div className="App">
      {chatElements}
      <form
        className="chat-form"
        onSubmit={async (e) => {
          e.preventDefault();
          await chatsManager.send(name, message);
          setName("");
          setMessage("");
        }}
      >
        <input name="name" onChange={withEventValue(setName)} />
        <textarea
          name="message"
          onChange={withEventValue(setMessage)}
        ></textarea>
        <button>Submit</button>
      </form>
    </div>
  );
}

function withEventValue(
  func: Function
): (
  e: ChangeEvent<HTMLInputElement> | ChangeEvent<HTMLTextAreaElement>
) => void {
  return (
    e: ChangeEvent<HTMLInputElement> | ChangeEvent<HTMLTextAreaElement>
  ) => func(e.target.value);
}

function generateRandomString(len: Number): String {
  const charList =
    "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
  return Array.from(Array(len))
    .map(() => charList[Math.floor(Math.random() * charList.length)])
    .join("");
}

export default App;
