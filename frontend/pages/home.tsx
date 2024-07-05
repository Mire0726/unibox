import React, { useState, useEffect } from "react";
import { useRouter } from "next/router";
import { onAuthStateChanged, getIdToken } from "firebase/auth";
import { auth } from "../firebase/auth";

const Chat = () => {
  const [messages, setMessages] = useState<string[]>([]);
  const [inputText, setInputText] = useState("");
  const [ws, setWs] = useState<WebSocket | null>(null);
  const router = useRouter();
  const backendUrl = "http://localhost:8080";
  const workspaceId = router.query.workspaceID;
  console.log("workspaceId:", workspaceId);
  const channelId = "testchannelID";

  useEffect(() => {
    const newWs = new WebSocket(`${backendUrl}/ws`);
    newWs.onopen = () => {
      console.log("WebSocket connection established");
    };
    newWs.onmessage = (event) => {
      try {
        const messageData = JSON.parse(event.data);
        if (messageData.message) {
          setMessages((prev) => [...prev, messageData.message]);
          console.log("Received message:", messageData.message);
        } else {
          console.error(
            "Received data does not contain 'message' key:",
            messageData
          );
        }
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
      const response = await fetch(
        `${backendUrl}/workspaces/${workspaceId}/channels/${channelId}/messages`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            content: inputText,
          }),
        }
      );
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
      <p>Messages:</p>
      <ul>
        {messages.length === 0 && <li>No messages yet</li>}
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
