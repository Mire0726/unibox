import React, { useState, useEffect } from "react";
import { useRouter } from "next/router";
import { onAuthStateChanged, getIdToken } from "firebase/auth";
import { auth } from "../firebase/auth";
import {
  Box,
  VStack,
  HStack,
  Text,
  Input,
  Button,
  Container,
  Heading,
  List,
  ListItem,
  useColorModeValue,
} from "@chakra-ui/react";

const Chat = () => {
  const [messages, setMessages] = useState<string[]>([]);
  const [inputText, setInputText] = useState("");
  const [ws, setWs] = useState<WebSocket | null>(null);
  const router = useRouter();
  const backendUrl = "http://localhost:8080";
  const workspaceId = router.query.workspaceID;
  console.log("workspaceId:", workspaceId);
  const channelId = "testchannelID";
  const borderColor = useColorModeValue("lavender.200", "lavender.700");

  useEffect(() => {
    const fetchMessages = async () => {
      const token = localStorage.getItem("idToken");
      if (!token) {
        console.error("No token found, please login again");
        return;
      }
  
      try {
        const response = await fetch(
          `${backendUrl}/workspaces/${workspaceId}/channels/${channelId}/messages`,
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          }
        );
        if (response.ok) {
          const data = await response.json();
          if (Array.isArray(data.messages)) {
            setMessages(data.messages.map(msg => msg.content));
          } else {
            console.error("Unexpected response structure:", data);
          }
        } else {
          console.error("Failed to fetch messages:", await response.text());
        }
      } catch (error) {
        console.error("Error fetching messages:", error);
      }
    };
  
    const setupWebSocket = () => {
      const newWs = new WebSocket(
        `${backendUrl}/ws?workspaceID=${workspaceId}&channelID=${channelId}`
      );
  
      newWs.onopen = () => {
        console.log("WebSocket connection established");
        fetchMessages();
      };
  
      newWs.onmessage = (event) => {
        try {
          const messageData = JSON.parse(event.data);
          console.log("Parsed message data:", messageData);
          if (Array.isArray(messageData.messages)) {
            setMessages(prevMessages => [...prevMessages, ...messageData.messages.map(msg => msg.content)]);
            console.log("Received messages:", messageData.messages);
          } else if (messageData.content) {
            setMessages(prevMessages => [...prevMessages, messageData.content]);
            console.log("Received message:", messageData.content);
          } else {
            console.error("Received data does not contain expected structure:", messageData);
          }
        } catch (error) {
          console.error("Error parsing message data:", error);
        }
      };
  
      newWs.onerror = (error) => {
        console.error("WebSocket error:", error);
      };
  
      setWs(newWs);
  
      return newWs;
    };
  
    const webSocket = setupWebSocket();
  
    const unsubscribe = onAuthStateChanged(auth, async (user) => {
      if (!user) {
        router.push("/index");
      } else {
        const token = await getIdToken(user);
        localStorage.setItem("idToken", token);
      }
    });
  
    return () => {
      webSocket.close();
      unsubscribe();
    };
  }, [router, workspaceId, channelId, backendUrl, auth]);
  
  // useEffect(() => {
  //   const newWs = new WebSocket(
  //     `${backendUrl}/ws?workspaceID=${workspaceId}&channelID=${channelId}`
  //   );
  //   newWs.onopen = () => {
  //     console.log("WebSocket connection established");
  //     console.log(
  //       "workspaceId!:",
  //       `${backendUrl}/ws?workspaceID=${workspaceId}&channelID=${channelId}`
  //     );
  //   };

  //   newWs.onmessage = (event) => {
  //     try {
  //       const messageData = JSON.parse(event.data);
  //       console.log("Parsed message data:", messageData);
  //       if (messageData.Content) {
  //         setMessages((prev) => [...prev, messageData.Content]);
  //         console.log("Received message:", messageData.Content);
  //       } else {
  //         console.error(
  //           "Received data does not contain 'Content' key:",
  //           messageData
  //         );
  //       }
  //     } catch (error) {
  //       console.error("Error parsing message data:", error);
  //     }
  //   };

  //   setWs(newWs);

  //   const unsubscribe = onAuthStateChanged(auth, async (user) => {
  //     if (!user) {
  //       router.push("/index");
  //     } else {
  //       const token = await getIdToken(user);
  //       localStorage.setItem("idToken", token);
  //     }
  //   });

  //   return () => {
  //     newWs.close();
  //     unsubscribe();
  //   };
  // }, [router, workspaceId, channelId]);

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
    <Container
      maxW="container.md"
      py={8}
      display="flex"
      flexDirection="column"
      alignItems="center"
    >
      <VStack spacing={6} align="stretch">
        <Heading as="h1" size="30px" textAlign="center" color="#2D3748">
          Chat Messages
        </Heading>
        <Box
          bg="#e6e6fa"
          borderRadius="md"
          p={4}
          borderWidth={1}
          borderColor={borderColor}
          height="60vh"
          overflowY="auto"
          alignItems={"center"}
        >
          <List spacing={3}>
            {messages.length === 0 && (
              <ListItem fontStyle="italic">No messages yet</ListItem>
            )}
            {messages.map((msg, index) => (
              <ListItem key={index} p={3} borderRadius="md" boxShadow="sm">
                <Text>{msg}</Text>
              </ListItem>
            ))}
          </List>
        </Box>
        <form onSubmit={handleFormSubmit}>
          <HStack>
            <Input
              value={inputText}
              onChange={(e) => setInputText(e.target.value)}
              placeholder="Type your message..."
              color="#696969"
              height={20}
              width={200}
            />
            <Button
              type="submit"
              bg="#e6e6fa"
              _hover={{ bg: "#2D3748", color: "#e6e6fa" }}
              width="full"
              borderRadius={10}
              alignItems="center"
            >
              Send
            </Button>
          </HStack>
        </form>
      </VStack>
    </Container>
  );
};

export default Chat;
