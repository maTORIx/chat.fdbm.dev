import { Chat } from "../gen/chat/v1/chat_pb"
import { userStylesDict } from "../global"
import { generateRandomColorFromText } from "../utils"
import "./chats.css"

type ChatProps = {
  chat: Chat
}

export const ChatElement = (props: ChatProps) => {
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

type ChatsProps = {
  chats: Chat[],
  onScroll: (e: any) => void
}

export const ChatListElement = (props: ChatsProps) => {
  let chatElements = [];
  for (let chat of props.chats) {
    chatElements.push(<ChatElement key={chat.id} chat={chat} />);
  }
  return (
    <div className="chats-container" id="chats-container" onScroll={props.onScroll}>
      {chatElements}
    </div>
  )
}