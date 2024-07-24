import { Chat } from "../components/chat";
import { Heading } from "@chakra-ui/react";

export default function Home() {
  return (
    <>
      <Heading as="h1" size="30px" textAlign="center" color="#2D3748">
        Chat Messages
      </Heading>
      <Chat />;
    </>
  );
}
