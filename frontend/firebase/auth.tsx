import { useEffect } from "react";
import { useRouter } from "next/router";
import { onAuthStateChanged, getIdToken } from "firebase/auth";
import { auth } from "./firebaseConfig";

const useAuth = () => {
  const router = useRouter();

  useEffect(() => {
    const unsubscribe = onAuthStateChanged(auth, async (user) => {
      if (!user) {
        router.push("/index");
      } else {
        const token = await getIdToken(user);
        localStorage.setItem("token", token);
      }
    });

    return () => unsubscribe();
  }, [router]);

  return null;
};

export default useAuth;
export { auth };
