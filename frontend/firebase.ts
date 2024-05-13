// firebase.ts
import { initializeApp } from "firebase/app";
import { getAnalytics } from "firebase/analytics";
import { getAuth, GoogleAuthProvider, signInWithPopup } from "firebase/auth";

const firebaseConfig = {
  apiKey: process.env.NEXT_PUBLIC_FIREBASE_API_KEY,
  authDomain: process.env.NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN,
  projectId: process.env.NEXT_PUBLIC_FIREBASE_PROJECT_ID,
  storageBucket: process.env.NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET,
  messagingSenderId: process.env.NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID,
  appId: process.env.NEXT_PUBLIC_FIREBASE_APP_ID,
  measurementId: process.env.NEXT_PUBLIC_FIREBASE_MEASUREMENT_ID,
};

// Firebase アプリの初期化
const app = initializeApp(firebaseConfig);

// Analytics を有効化
const analytics = typeof window !== "undefined" ? getAnalytics(app) : null;

// 認証サービスの取得
const auth = getAuth(app);

// Google プロバイダーのセットアップ
const provider = new GoogleAuthProvider();

// Googleログインの関数
const signInWithGoogle = async () => {
  try {
    const result = await signInWithPopup(auth, provider);
    // 成功した場合、ユーザー情報が result.user に含まれる
    console.log(result.user);
  } catch (error) {
    console.error(error);
  }
};

export { signInWithGoogle, auth };
