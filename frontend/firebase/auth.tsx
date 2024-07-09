import { useEffect, useCallback } from "react";
import { useRouter } from "next/router";
import {
  onAuthStateChanged,
  getIdToken,
  signInWithPopup,
  GoogleAuthProvider,
} from "firebase/auth";
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

  const signInWithGoogle = useCallback(async () => {
    try {
      const provider = new GoogleAuthProvider();
      const result = await signInWithPopup(auth, provider);
      const user = result.user;
      const token = await getIdToken(user);
      localStorage.setItem("token", token);
      router.push("/workspace");
    } catch (error) {
      console.error("Google sign in error:", error);
    }
  }, [router]);

  return { signInWithGoogle };
};

export default useAuth;
export { auth };
