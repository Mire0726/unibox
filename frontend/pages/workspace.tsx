import React, { useState, useEffect } from "react";
import axios from "axios";
import { useRouter } from "next/router";
import styles from "./index.module.scss";
import { onAuthStateChanged, getIdToken } from "firebase/auth";
import { auth } from "../firebase/auth";

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
        router.push("/index");
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
        router.push("/home");
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
          Authorization: `Bearer ${token}`, // Bearer prefix added
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
    <>
      <form onSubmit={handleSubmitLogin} className={styles.loginForm}>
        <h3>Login</h3>
        <input
          type="text"
          value={loginID}
          onChange={(e) => setLoginID(e.target.value)}
          placeholder="Workspace ID"
          required
        />
        <input
          type="password"
          value={loginPassword}
          onChange={(e) => setLoginPassword(e.target.value)}
          placeholder="Password"
          required
        />
        <button type="submit">Log in to Workspace</button>
      </form>
      <form
        onSubmit={handleSubmitRegistration}
        className={styles.registrationForm}
      >
        <h3>Register</h3>
        <input
          type="text"
          value={registerName}
          onChange={(e) => setRegisterName(e.target.value)}
          placeholder="Workspace Name"
          required
        />
        <input
          type="password"
          value={registerPassword}
          onChange={(e) => setRegisterPassword(e.target.value)}
          placeholder="Password"
          required
        />
        <button type="submit">Register Workspace</button>
      </form>
    </>
  );
}
