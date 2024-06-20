import React, { useState, useEffect } from "react";
import { useRouter } from "next/router";
import { onAuthStateChanged, getIdToken } from "firebase/auth"; // Firebaseの認証関連の関数
import { auth } from "../firebase/auth"; // 正確なパスに修正

const Chat = () => {
  const [messages, setMessages] = useState<string[]>([]);
  const [inputText, setInputText] = useState("");
  const [ws, setWs] = useState<WebSocket | null>(null);
  const router = useRouter();
  const backendUrl = "http://localhost:8080";

  useEffect(() => {
    const newWs = new WebSocket(`${backendUrl}/ws`);
    newWs.onopen = () => {
      console.log("WebSocket connection established");
    };
    newWs.onmessage = (event) => {
      try {
        const messageData = JSON.parse(event.data);
        setMessages((prev) => [...prev, messageData.message]); 
      } catch (error) {
        console.error("Error parsing message data:", error);
      }
    };
    setWs(newWs);

    const unsubscribe = onAuthStateChanged(auth, async (user) => {
      if (!user) {
        router.push("/index");
      } else {
        const token = await getIdToken(user);
        localStorage.setItem("idToken", token);
      }
    });

    return () => {
      newWs.close();
      unsubscribe();
    };
  }, [router]);

  const handleFormSubmit = async (e) => {
    e.preventDefault();
    const token = localStorage.getItem("idToken");
    if (!token) {
      console.error("No token found, please login again");
      return;
    }
    try {
      const response = await fetch(`${backendUrl}/messages`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`, // Bearer prefix added
        },
        body: JSON.stringify({
          channelId: "0d1b6af3-1ce0-11ef-bbe1-0242ac150003",
          content: inputText,
        }),
      });
      if (response.ok) {
        sendMessage(inputText);
        setInputText("");
      } else {
        const data = await response.json();
        console.error("Failed to send message:", data.message); 
      }
    } catch (error) {
      console.error("Failed to send message err:", error);
    }
  };

  const sendMessage = (message) => {
    if (message.trim() !== "") {
      ws?.send(message);
    }
  };
  return (
    <div>
      <h1>Chat Messages</h1>
      <ul>
        {messages.map((msg, index) => (
          <li key={index}>{msg}</li>
        ))}
      </ul>
      <form onSubmit={handleFormSubmit}>
        <input
          type="text"
          value={inputText}
          onChange={(e) => setInputText(e.target.value)}
          placeholder="Type your message..."
        />
        <button type="submit">Send Message</button>
      </form>
    </div>
  );
};

export default Chat;
