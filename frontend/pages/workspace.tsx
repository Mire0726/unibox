import React, { useState, useEffect } from "react";
import axios from "axios";
import { useRouter } from "next/router";
import { onAuthStateChanged, getIdToken } from "firebase/auth";
import { auth } from "../firebase/auth";
import {
  Input,
  Button,
  Heading,
  Flex,
  VStack,
  FormControl,
  FormLabel,
} from "@chakra-ui/react";

export default function Login() {
  const [loginID, setLoginID] = useState("");
  const [loginPassword, setLoginPassword] = useState("");
  const [registerName, setRegisterName] = useState("");
  const [registerPassword, setRegisterPassword] = useState("");
  const router = useRouter();
  const backendUrl = "http://localhost:8080";

  useEffect(() => {
    const unsubscribe = onAuthStateChanged(auth, async (user) => {
      if (!user) {
        router.push("/");
      } else {
        const token = await getIdToken(user);
        localStorage.setItem("idToken", token);
      }
    });

    return () => {
      unsubscribe();
    };
  }, [router]);

  const handleSubmitLogin = async (e) => {
    e.preventDefault();
    const token = localStorage.getItem("idToken");
    if (!token) {
      console.error("No token found, please login again");
      return;
    }
    try {
      const response = await fetch(`${backendUrl}/workspaces/login`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          ID: loginID,
          password: loginPassword,
        }),
      });
      if (response.ok) {
        router.push(`/home?workspaceID=${loginID}`);
      } else {
        const data = await response.json();
        console.error("Failed to login workspace:", data.message);
      }
    } catch (error) {
      console.error(error);
      alert("Login failed.");
    }
  };

  const handleSubmitRegistration = async (e) => {
    e.preventDefault();
    const token = localStorage.getItem("idToken");
    if (!token) {
      console.error("No token found, please login again");
      return;
    }
    try {
      const response = await fetch(`${backendUrl}/workspaces`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          name: registerName,
          password: registerPassword,
        }),
      });
      const data = await response.json();
      if (response.ok) {
        setLoginID(data.ID);
        alert("This workspace ID is " + data.ID);
        setLoginPassword(registerPassword);
        handleSubmitLogin(e);
      } else {
        alert("Failed to register. No ID returned: " + data.message);
      }
    } catch (error) {
      console.error(error);
      alert(
        "Registration failed. Error: " +
          (error.response?.data.message || error.message)
      );
    }
  };

  return (
    <Flex height="100vh" alignItems="center" justifyContent="center">
      <Flex
        direction="column"
        padding={100}
        rounded={50}
        minWidth="550px"
        maxWidth="90%"
      >
        <VStack
          spacing={2}
          align="center"
          bg="#e6e6fa"
          width={400}
          p={8}
          borderRadius="12px"
          boxShadow="lg"
        >
          <form onSubmit={handleSubmitLogin}>
            <Heading
              padding={10}
              color="#696969"
              textAlign="center"
              fontFamily="Arial"
            >
              ログイン
            </Heading>
            <FormControl id="loginID" mb={4}>
              <FormLabel color="#696969">Workspace ID</FormLabel>
              <Input
                placeholder="Workspace ID"
                type="text"
                value={loginID}
                onChange={(e) => setLoginID(e.target.value)}
                required
                height={30}
                width={300}
              />
            </FormControl>
            <FormControl id="loginPassword" mb={4}>
              <FormLabel color="#696969">Password</FormLabel>
              <Input
                placeholder="Password"
                type="password"
                value={loginPassword}
                onChange={(e) => setLoginPassword(e.target.value)}
                required
                color="#696969"
                height={30}
                width={300}
              />
            </FormControl>
            <Button
              type="submit"
              bg="#e6e6fa"
              _hover={{ bg: "#2D3748", color: "#e6e6fa" }}
              width="full"
              borderRadius={10}
              alignItems="center"
            >
              ログイン
            </Button>
          </form>

          <form onSubmit={handleSubmitRegistration}>
            <Heading
              padding={1}
              color="#696969"
              textAlign="center"
              fontFamily="Arial"
              fontSize={15}
            >
              新規登録
            </Heading>
            <FormControl id="registerName" mb={4}>
              <FormLabel color="#696969">Workspace Name</FormLabel>
              <Input
                placeholder="Workspace Name"
                type="text"
                value={registerName}
                onChange={(e) => setRegisterName(e.target.value)}
                required
                height={20}
                width={200}
              />
            </FormControl>
            <FormControl id="registerPassword" mb={4}>
              <FormLabel color="#696969">Password</FormLabel>
              <Input
                placeholder="Password"
                type="password"
                value={registerPassword}
                onChange={(e) => setRegisterPassword(e.target.value)}
                required
                color="#696969"
                height={20}
                width={200}
              />
            </FormControl>
            <Button
              type="submit"
              bg="#e6e6fa"
              _hover={{ bg: "#2D3748", color: "#e6e6fa" }}
              width="full"
              borderRadius={10}
              alignItems="center"
            >
              登録
            </Button>
          </form>
        </VStack>
      </Flex>
    </Flex>
  );
}
