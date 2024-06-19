import React, { useState, useEffect } from "react";

const Chat = () => {
  const [messages, setMessages] = useState<string[]>([]);
  const [ws, setWs] = useState<WebSocket | null>(null);

  useEffect(() => {
    // WebSocketの接続を開く
    const newWs = new WebSocket("ws://localhost:8080/ws");

    newWs.onmessage = (event) => {
      // サーバーからのメッセージを受信したら状態を更新
      setMessages((prev) => [...prev, event.data]);
    };

    setWs(newWs);

    // コンポーネントのアンマウント時にWebSocketを閉じる
    return () => {
      newWs.close();
    };
  }, []);

  // メッセージ送信用の関数
  const sendMessage = (message: string) => {
    ws?.send(message);
  };

  return (
    <div>
      <h1>Chat Messages</h1>
      <ul>
        {messages.map((msg, index) => (
          <li key={index}>{msg}</li>
        ))}
      </ul>
      <button onClick={() => sendMessage("Hello from Client!")}>
        Send Message
      </button>
    </div>
  );
};

export default Chat;
