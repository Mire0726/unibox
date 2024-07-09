import { useState } from "react";
import axios from "axios";
import { useRouter } from "next/router";
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
        router.push("/workspace");
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
      handleSubmitLogin(e);
    } catch (error) {
      console.error(error);
      alert("登録に失敗しました。エラー: " + error.response.data.message);
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
            <FormControl id="email" mb={4}>
              <FormLabel color="#696969">Email</FormLabel>
              <Input
                placeholder="mail@sample.com"
                type="email"
                value={loginEmail}
                onChange={(e) => setLoginEmail(e.target.value)}
                required
                height={30}
                width={300}
              />
            </FormControl>
            <FormControl id="password" mb={4}>
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
            <FormControl id="Name" mb={4}>
              <FormLabel color="#696969">Name</FormLabel>
              <Input
                placeholder="Name"
                type="name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                required
                height={20}
                width={200}
              />
            </FormControl>
            <FormControl id="email" mb={4}>
              <FormLabel color="#696969">Email</FormLabel>
              <Input
                placeholder="mail@sample.com"
                type="email"
                value={loginEmail}
                onChange={(e) => setRegisterEmail(e.target.value)}
                required
                height={20}
                width={200}
              />
            </FormControl>
            <FormControl id="password" mb={4}>
              <FormLabel color="#696969">Password</FormLabel>
              <Input
                placeholder="Password"
                type="password"
                value={loginPassword}
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
