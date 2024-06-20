import { useState } from "react";
import axios from "axios";
import { useRouter } from "next/router";
import styles from "./index.module.scss";

export default function Login() {
  const [loginEmail, setLoginEmail] = useState("");
  const [loginPassword, setLoginPassword] = useState("");
  const [name, setName] = useState("");
  const [registerEmail, setRegisterEmail] = useState("");
  const [registerPassword, setRegisterPassword] = useState("");
  const router = useRouter();
  const backendUrl = "http://localhost:8080";

  const handleSubmitLogin = async (e) => {
    e.preventDefault();
    try {
      const response = await axios.post(`${backendUrl}/signIn`, {
        email: loginEmail,
        password: loginPassword,
      });
      if (response.data.idToken) {  
        localStorage.setItem("idToken", response.data.idToken);
        router.push("/home");
      } else {
        console.error("トークンがレスポンスに含まれていません。");
        alert("ログインに問題がありました。もう一度お試しください。");
      }
    } catch (error) {
      console.error(error);
      alert("ログインに失敗しました。");
    }
  };
  

  const handleSubmitRegistration = async (e) => {
    e.preventDefault();
    try {
      const response = await axios.post(`${backendUrl}/signUp`, {
        name: name,
        email: registerEmail,
        password: registerPassword,
      });

      alert("登録が完了しました。ログインしてください。");
      setLoginEmail(registerEmail);
      setLoginPassword(registerPassword);
      handleSubmitLogin(e); // 登録後自動でログイン
    } catch (error) {
      console.error(error);
      alert("登録に失敗しました。エラー: " + error.response.data.message);
    }
  };

  return (
    <>
      <form onSubmit={handleSubmitLogin} className={styles.loginForm}>
        <h3>ログイン</h3>
        <input
          type="email"
          value={loginEmail}
          onChange={(e) => setLoginEmail(e.target.value)}
          placeholder="Email"
          required
        />
        <input
          type="password"
          value={loginPassword}
          onChange={(e) => setLoginPassword(e.target.value)}
          placeholder="Password"
          required
        />
        <button type="submit">ログイン</button>
      </form>
      <form
        onSubmit={handleSubmitRegistration}
        className={styles.registrationForm}
      >
        <h3>新規登録</h3>
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Name"
          required
        />
        <input
          type="email"
          value={registerEmail}
          onChange={(e) => setRegisterEmail(e.target.value)}
          placeholder="Email"
          required
        />
        <input
          type="password"
          value={registerPassword}
          onChange={(e) => setRegisterPassword(e.target.value)}
          placeholder="Password"
          required
        />
        <button type="submit">登録</button>
      </form>
      <img
        src="/cat1.png"
        alt="Description of image"
        className={styles.cardImage}
      />
    </>
  );
}
